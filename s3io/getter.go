package s3io

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"sync"
	"syscall"
)

const qWaitMax = 2

type Getter struct {
	url   url.URL
	b     *Bucket
	bufsz int64
	err   error

	chunkID    int
	rChunk     *chunk
	contentLen int64 // Total length of object
	bytesRead  int64
	chunkTotal int

	readCh   chan *chunk
	getCh    chan *chunk
	errCh    chan error // errors from chunk workers
	quit     chan struct{}
	qWait    map[int]*chunk
	qWaitLen uint
	cond     sync.Cond

	sp BufferPool

	closed bool
	c      *Config

	md5  hash.Hash
	cIdx int64
}

type chunk struct {
	id     int
	header http.Header
	start  int64
	size   int64
	b      []byte
	r      io.ReadCloser
}

//func newGetter(getURL url.URL, c *Config, b *Bucket) (*Getter, http.Header, error) {
func newGetter(getURL url.URL, c *Config, b *Bucket) (*Getter, http.Header, error) {
	g := new(Getter)
	g.url = getURL
	g.c, g.b = new(Config), new(Bucket)
	*g.c, *g.b = *c, *b
	g.bufsz = c.PartSize
	g.c.NTry = max(c.NTry, 1)
	g.c.Concurrency = max(c.Concurrency, 1)

	g.getCh = make(chan *chunk)
	g.readCh = make(chan *chunk)
	g.errCh = make(chan error, g.c.Concurrency)
	g.quit = make(chan struct{})
	g.qWait = make(map[int]*chunk)
	g.b = b
	g.md5 = md5.New()
	g.cond = sync.Cond{L: &sync.Mutex{}}

	// use get instead of head for error messaging
	resp, err := g.retryRequest("GET", g.url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	// If we got redirected by the v2 files service then set the the file URL to the new location.
	if resp.Request.URL.Host != g.url.Host {
		g.url = *resp.Request.URL
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, nil, newRespError(resp)
	}

	// Golang changes content-length to -1 when chunked transfer encoding / EOF close response detected
	if resp.ContentLength == -1 {
		return nil, nil, fmt.Errorf("Retrieving objects with undefined content-length " +
			" responses (chunked transfer encoding / EOF close) is not supported")
	}

	g.contentLen = resp.ContentLength
	g.chunkTotal = int((g.contentLen + g.bufsz - 1) / g.bufsz) // round up, integer division
	//logger.Printf("object size: %3.2g MB", float64(g.contentLen)/float64((1*mb)))

	g.sp = bufferPool(g.bufsz)

	for i := 0; i < g.c.Concurrency; i++ {
		go g.worker()
	}
	go g.initChunks()
	return g, resp.Header, nil
}

func (g *Getter) retryRequest(method, urlStr string, body io.ReadSeeker) (resp *http.Response, err error) {
	for i := 0; i < g.c.NTry; i++ {
		var req *http.Request
		req, err = http.NewRequest(method, urlStr, body)
		if err != nil {
			return
		}

		// apply custom headers
		for k, v := range g.c.Headers {
			req.Header[k] = v
		}

		if body != nil {
			req.Header.Set(sha256Header, shaReader(body))
		}

		g.b.Sign(req)

		resp, err = g.c.Client.Do(req)
		err = checkRequest(i, resp, err)
		if err == nil {
			return
		} else {
			if body != nil {
				if _, err = body.Seek(0, 0); err != nil {
					return
				}
			}
		}
	}
	return
}

func (g *Getter) initChunks() {
	id := 0
	defer close(g.getCh)
	for i := int64(0); i < g.contentLen; {
		size := min64(g.bufsz, g.contentLen-i)
		c := &chunk{
			id: id,
			header: http.Header{
				"Range": {fmt.Sprintf("bytes=%d-%d",
					i, i+size-1)},
			},
			start: i,
			size:  size,
			b:     nil,
		}

		for k, v := range g.c.Headers {
			c.header[k] = v
		}

		i += size
		id++

		select {
		case <-g.quit:
			return
		case g.getCh <- c:
		}
	}
}

func (g *Getter) worker() {
	for c := range g.getCh {
		g.retryGetChunk(c)
	}
}

func (g *Getter) retryGetChunk(c *chunk) {
	var err error
	select {
	case c.b = <-g.sp.Get():
	case <-g.quit:
		return
	}

	for i := 0; i < g.c.NTry; i++ {
		err = g.getChunk(c)
		if err == nil {
			return
		}
		logger.Printf("error on attempt %d: retrying chunk: %v, error: %s", i, c.id, err)
		sleep(i)
	}
	// An error occurred, either exit on quit or send error
	select {
	case <-g.quit:
	case g.errCh <- err:
	}
}

func (g *Getter) getChunk(c *chunk) (err error) {
	// ensure buffer is empty
	r, err := http.NewRequest("GET", g.url.String(), nil)
	if err != nil {
		return err
	}
	r.Header = c.header
	g.b.Sign(r)
	resp, err := g.c.Client.Do(r)
	if err != nil {
		return err
	}
	defer checkClose(resp.Body, err)
	if resp.StatusCode != 206 && resp.StatusCode != 200 {
		return newRespError(resp)
	}
	n, err := io.ReadAtLeast(resp.Body, c.b, int(c.size))
	if err != nil {
		return err
	}
	if err := resp.Body.Close(); err != nil {
		return err
	}
	if int64(n) != c.size {
		return fmt.Errorf("chunk %d: Expected %d bytes, received %d",
			c.id, c.size, n)
	}
	select {
	case <-g.quit:
		return syscall.EINVAL
	case g.readCh <- c:
	}

	// wait for qWait to drain before starting next chunk
	g.cond.L.Lock()
	for g.qWaitLen >= qWaitMax {
		if g.closed {
			return nil
		}
		g.cond.Wait()
	}
	g.cond.L.Unlock()
	return nil
}

func (g *Getter) Read(p []byte) (int, error) {
	var err error
	if g.closed {
		return 0, syscall.EINVAL
	}
	if g.err != nil {
		return 0, g.err
	}
	nw := 0
	for nw < len(p) {
		if g.bytesRead == g.contentLen {
			return nw, io.EOF
		} else if g.bytesRead > g.contentLen {
			// Here for robustness / completeness
			// Should not occur as golang uses LimitedReader up to content-length
			return nw, fmt.Errorf("expected %d bytes, received %d (too many bytes)",
				g.contentLen, g.bytesRead)
		}

		// If for some reason no more chunks to be read and bytes are off, error, incomplete result
		if g.chunkID >= g.chunkTotal {
			return nw, fmt.Errorf("expected %d bytes, received %d and chunkID %d >= chunkTotal %d (no more chunks remaining)",
				g.contentLen, g.bytesRead, g.chunkID, g.chunkTotal)
		}

		if g.rChunk == nil {
			g.rChunk, err = g.nextChunk()
			if err != nil {
				return nw, err
			}
			g.cIdx = 0
		}

		n := copy(p[nw:], g.rChunk.b[g.cIdx:g.rChunk.size])
		g.cIdx += int64(n)
		nw += n
		g.bytesRead += int64(n)

		if g.cIdx >= g.rChunk.size { // chunk complete
			g.sp.Give() <- g.rChunk.b
			g.chunkID++
			g.rChunk = nil
		}
	}
	return nw, nil
}

func (g *Getter) nextChunk() (*chunk, error) {
	for {
		// first check qWait
		c := g.qWait[g.chunkID]
		if c != nil {
			delete(g.qWait, g.chunkID)
			g.cond.L.Lock()
			g.qWaitLen--
			g.cond.L.Unlock()
			g.cond.Signal() // wake up waiting worker goroutine
			if _, err := g.md5.Write(c.b[:c.size]); err != nil {
				return nil, err
			}
			return c, nil
		}
		// if next chunk not in qWait, read from channel
		select {
		case c := <-g.readCh:
			g.qWait[c.id] = c
			g.cond.L.Lock()
			g.qWaitLen++
			g.cond.L.Unlock()
		case <-g.quit:
			return nil, syscall.EINVAL // fatal error, quit.
		case err := <-g.errCh:
			g.err = err
			return nil, err
		}
	}
}

func (g *Getter) Close() error {
	if g.closed {
		return syscall.EINVAL
	}
	g.closed = true
	if g.sp != nil {
		g.sp.Quit()
	}
	close(g.quit)
	g.cond.Broadcast()
	if g.err != nil {
		return g.err
	}
	if g.bytesRead != g.contentLen {
		return fmt.Errorf("read error: %d bytes read. expected: %d", g.bytesRead, g.contentLen)
	}
	return nil
}

func (g *Getter) Sum(b []byte) []byte {
	return g.md5.Sum(b)
}
