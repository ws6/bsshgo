package s3io

import (
	"bytes"

	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// convenience multipliers
const (
	_        = iota
	kb int64 = 1 << (10 * iota)
	mb
	gb
	tb
	pb
	eb
)

const BackOffRate time.Duration = 600

func sleep(i int) {
	time.Sleep(time.Duration(math.Exp2(float64(i))) * BackOffRate * time.Millisecond)
}

// Min and Max functions
func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func newRespError(r *http.Response) *RespError {
	defer r.Body.Close()
	e := new(RespError)
	e.StatusCode = r.StatusCode
	b, _ := ioutil.ReadAll(r.Body)
	err := xml.NewDecoder(bytes.NewReader(b)).Decode(e) // parse error from response
	if err != nil {
		e.Code = "GENERIC_UNKNOWN"
		e.Message = string(b)
	}
	if e.Message == "" {
		e.Message = fmt.Sprintf("HTTP %d", e.StatusCode)
	}
	return e
}

// RespError represents an http error response
// http://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html
type RespError struct {
	Code       string
	Message    string
	Resource   string
	RequestID  string `xml:"RequestId"`
	StatusCode int
}

func (e *RespError) Error() string {
	return fmt.Sprintf(
		"S3 [%s]:  %s",
		e.Code,
		e.Message,
	)
}

func checkClose(c io.Closer, err error) {
	if c != nil {
		cerr := c.Close()
		if err == nil {
			err = cerr
		}
	}

}

func checkRequest(i int, resp *http.Response, err error) (e error) {
	if resp == nil {
		if err != nil {
			e = err
		} else {
			e = fmt.Errorf("unexpected nil response")
		}
	} else {
		if resp.StatusCode >= 300 {
			e = newRespError(resp)
		} else {
			e = err
		}
	}

	if e != nil {
		logger.Printf("Error on attempt %d: Response: %s", i, e)
		sleep(i)
	}
	return
}

// Derive the part size for a file of known size
func GetPartSize(partsize, filesize int64) int64 {
	// while parts needed is less than max number of parts,
	// up the part size
	for partsize*10000 < filesize {
		partsize += 5 * mb
	}
	return partsize
}
