package download

import (
	"context"
	"fmt"
	"os"

	"io"
	"path/filepath"

	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ws6/bsshgo"
	"github.com/ws6/bsshgo/s3io"
)

//download_file.go download bssh files
type Options struct {
	DestinationPrefix string
	NumWorkers        int
	//a SkipFn removes the unwanted download job. if returns an error, it will exist the worker
	IsSkipFn func(*bsshgo.FileS3PresignedUrlResp) (bool, error)
}

type s3Loc struct {
	bucket string
	key    string
}

func getBucketAndKeyFromPresignedS3URL(_url string) (*s3Loc, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}

	ret := new(s3Loc)
	path := strings.SplitN(u.Path, "/", 3)
	if len(path) < 3 {
		return nil, fmt.Errorf(`malformatted URL:%s`, _url)
	}
	ret.bucket = path[1]
	ret.key = path[2]
	return ret, nil
}

func getPresignedS3UrlReader(_url string) (io.ReadCloser, error) {
	u, err := url.Parse(_url)
	if err != nil {

		return nil, err
	}

	r, _, err := s3io.PreSigned(*u, nil)

	return r, err
}

func downloadWorker(ctx context.Context, file *bsshgo.FileS3PresignedUrlResp, opts *Options) error {

	done := make(chan struct{})

	defer func() {
		done <- struct{}{}
	}()

	loc, err := getBucketAndKeyFromPresignedS3URL(file.HrefContent)
	if err != nil {
		return err
	}

	outputFile := filepath.Join(opts.DestinationPrefix, loc.key)
	rcClosed := false
	rc, err := getPresignedS3UrlReader(file.HrefContent)
	if err != nil {
		return err
	}

	defer func() {
		if !rcClosed { //not atomic
			rc.Close()
			rcClosed = true
		}
	}()

	if _, err := os.Stat(filepath.Dir(outputFile)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(outputFile), 0644); err != nil {
			if err != nil {
				return fmt.Errorf(`os.MkdirAll:%s`, err.Error())
			}
		}
	}
	ofh, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer ofh.Close()

	go func() {
		select {
		case <-ctx.Done():
			//early terminate the io.Copy

			if !rcClosed {
				rc.Close()
				rcClosed = true
			}
			return
		case <-done:
			return //as expected
		}

	}()

	if _, err := io.Copy(ofh, rc); err != nil {
		//if invalidate argument, it is most likely from the rc.Close()
		return fmt.Errorf(`io.Copy:%s`, err.Error())
	}

	return nil

}

func Download(ctx1 context.Context, ch <-chan *bsshgo.FileS3PresignedUrlResp, opts *Options) (int64, error) {
	var downloaded int64
	if opts.NumWorkers <= 0 {
		opts.NumWorkers = 1
	}
	ctx, cancelFn := context.WithCancel(ctx1)
	if _, err := os.Stat(opts.DestinationPrefix); os.IsNotExist(err) {
		os.MkdirAll(opts.DestinationPrefix, 0644)
	}
	defer cancelFn()
	var wg sync.WaitGroup
	errChan := make(chan error, opts.NumWorkers)
	defer close(errChan)
	for i := 0; i < opts.NumWorkers; i++ {
		wg.Add(1)
		go func(idx int) {

			defer wg.Done()
			for job := range ch {

				if err := downloadWorker(ctx, job, opts); err != nil {

					go func() {
						errChan <- err

					}()
					return
				}
				atomic.AddInt64(&downloaded, 1)
				select {
				case <-ctx.Done(): //they shall align with same ctx

					return
				default:
					continue
				}
			}
		}(i)
	}
	allDone := make(chan struct{})
	go func() {
		wg.Wait()
		allDone <- struct{}{}
	}()
	select {
	case <-allDone:

		return downloaded, ctx.Err()
	case errRet := <-errChan:

		return downloaded, errRet
	case <-ctx.Done():
		return downloaded, ctx.Err()
	}

	return downloaded, nil
}
