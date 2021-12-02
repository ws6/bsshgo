package s3io

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// See Amazon S3 Developer Guide for explanation
// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
var paramsToSign = map[string]bool{
	"acl":                          true,
	"location":                     true,
	"logging":                      true,
	"notification":                 true,
	"partNumber":                   true,
	"policy":                       true,
	"requestPayment":               true,
	"torrent":                      true,
	"uploadId":                     true,
	"uploads":                      true,
	"versionId":                    true,
	"versioning":                   true,
	"versions":                     true,
	"response-content-type":        true,
	"response-content-language":    true,
	"response-expires":             true,
	"response-cache-control":       true,
	"response-content-disposition": true,
	"response-content-encoding":    true,
}

type v2Signer struct {
	hm   io.Writer
	req  *http.Request
	Keys Keys
}

type v2 struct{}

func V2Signer() Signer {
	return v2{}
}

// Sign an AWS request using V2 signature method
func (v2) Sign(b *Bucket, req *http.Request) {
	s := v2Signer{req: req, Keys: b.S3.Keys}
	s.Sign()
}

// Signer that signs and AWS request using V2 signature methods detailed here:
// http://docs.aws.amazon.com/AmazonS3/latest/dev/RESTAuthentication.html
func (s *v2Signer) Sign() error {
	if dateHeader := s.req.Header.Get("Date"); dateHeader == "" {
		s.req.Header.Set("Date", time.Now().UTC().Format(http.TimeFormat))
	}
	if s.Keys.SecurityToken != "" {
		s.req.Header.Set("X-Amz-Security-Token", s.Keys.SecurityToken)
	}

	hm := hmac.New(sha1.New, []byte(s.Keys.SecretKey))
	s.hm = hm
	s.writeSignature()
	s.writeAmzHeaders()
	s.writeResourceHeaders()

	signature := make([]byte, base64.StdEncoding.EncodedLen(hm.Size()))
	base64.StdEncoding.Encode(signature, hm.Sum(nil))
	s.req.Header.Set("Authorization", "AWS "+s.Keys.AccessKey+":"+string(signature))

	return nil
}

func (s *v2Signer) writeSignature() {
	var n []byte = []byte{'\n'}
	s.hm.Write([]byte(s.req.Method))
	s.hm.Write(n)
	s.hm.Write([]byte(s.req.Header.Get("content-md5")))
	s.hm.Write(n)
	s.hm.Write([]byte(s.req.Header.Get("content-type")))
	s.hm.Write(n)
	if _, ok := s.req.Header["X-Amz-Date"]; !ok {
		s.hm.Write([]byte(s.req.Header.Get("date")))
	}
	s.hm.Write(n)
}

func (s *v2Signer) writeAmzHeaders() {
	var amzHeaders []string

	for h := range s.req.Header {
		if strings.HasPrefix(strings.ToLower(h), "x-amz-") {
			amzHeaders = append(amzHeaders, h)
		}
	}
	sort.Strings(amzHeaders)
	for _, h := range amzHeaders {
		v := s.req.Header[h]
		s.hm.Write([]byte(strings.ToLower(h)))
		s.hm.Write([]byte(":"))
		s.hm.Write([]byte(strings.Join(v, ",")))
		s.hm.Write([]byte("\n"))
	}
}

func (s *v2Signer) writeResourceHeaders() {
	// Currently only supports path-style resource
	u, _ := url.Parse(s.req.URL.String())
	s.hm.Write([]byte(u.EscapedPath()))

	var sr []string
	for k, vs := range s.req.URL.Query() {
		if paramsToSign[k] {
			for _, v := range vs {
				if v == "" {
					sr = append(sr, k)
				} else {
					sr = append(sr, k+"="+v)
				}
			}

		}
	}
	sort.Strings(sr)
	var q byte = '?'
	for _, str := range sr {
		s.hm.Write([]byte{q})
		s.hm.Write([]byte(str))
		q = '&'
	}

}
