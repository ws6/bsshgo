package s3io

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type S3DebugHandler struct {
}

func (S3DebugHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if _, ok := r.Header["Range"]; ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	for i := 0; i < 1000; i++ {
		rw.Write([]byte("ok"))
	}
}

// Tests that a broken chunk endpoint can fail and exit
func TestGetter_Retry(t *testing.T) {
	s := httptest.NewServer(S3DebugHandler{})
	defer s.Close()

	DefaultConfig.NTry = 3
	DefaultConfig.PartSize = 10

	logger.debug = true
	logger.SetOutput(os.Stderr)
	u, _ := url.Parse(s.URL)
	g, _, err := newGetter(*u, DefaultConfig, &Bucket{
		Config: DefaultConfig,
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(ioutil.Discard, g)
	if err == nil {
		t.Fatal("expected an error in io.Copy")
	}
	_, ok := err.(*RespError)
	if !ok {
		t.Fatalf("expected a response error, got %T", err)
	}
}
