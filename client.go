package bsshgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//illumina tss golang client
//http://support-docs.illumina.com/SW/TSSS/TruSight_SW_API/Content/SW/FrontPages/TruSightSoftware_API.htm
type Client struct {
	cfg        map[string]string
	httpclient *http.Client
	User       *UserResp
}

var (
	DEFAULT_TIMEOUT_SECOND = 60
)

const (
	AUTH_TOKEN = `x-access-token`
	// ILMN_DOMAIN    = `X-ILMN-Domain`
	// ILMN_WORKGROUP = `X-ILMN-Workgroup`
	CONTENT_TYPE = `Content-Type`
	BASE_URL     = `BASE_URL`
)

func requiredKeys() []string {
	return []string{
		AUTH_TOKEN,
	}
}

func NewClient(cfg map[string]string) (*Client, error) {
	ret := new(Client)
	ret.cfg = cfg
	if err := ret.ConfigCheck(); err != nil {
		return nil, fmt.Errorf(`ConfigCheck:%s`, err.Error())
	}

	baseUrl := ret.GetBaseUrl()
	if baseUrl == "" {
		return nil, fmt.Errorf(`no base url`)
	}

	ret.httpclient = new(http.Client)
	//you can use GetHttpClient then set the client paramters
	ret.httpclient.Timeout = time.Second * time.Duration(DEFAULT_TIMEOUT_SECOND)

	return ret, nil
}

func (self *Client) GetHttpClient() *http.Client {
	return self.httpclient
}

func (self *Client) ConfigCheck() error {
	for _, rk := range requiredKeys() {
		if _, ok := self.cfg[rk]; !ok {
			return fmt.Errorf(`missing required key:%s`, rk)
		}
	}
	return nil
}

func (self *Client) GetBaseUrl() string {
	//expect https://tss-test.trusight.illumina.com
	return self.cfg[`BASE_URL`]
}

func (self *Client) AttachHeaders(req *http.Request, url string) {
	for _, rk := range requiredKeys() {
		req.Header.Set(rk, self.cfg[rk])
	}

	req.Header.Set(CONTENT_TYPE, `application/json`)

	//for compatible with als domain api
	// for history API https://developer.basespace.illumina.com/docs/content/documentation/rest-api/history-api-reference#HistoryAPIReference
	// req.Header.Set(`Authorization`, "Bearer "+req.Header.Get(AUTH_TOKEN))
	// /v1/feeds/bssh
	//!!!the history API can not share with x-access-token header
	if strings.Contains(url, `/v1/feeds/bssh`) {
		req.Header.Del(`AUTH_TOKEN`)
		req.Header.Set(`Authorization`, "Bearer "+req.Header.Get(AUTH_TOKEN))
	}
}

type ModifyRequest func(*http.Request)

//NewRequestWithContext over write http.NewRequestWithContext wih auth headers
func (self *Client) NewRequestWithContext(ctx context.Context, method, url string, body io.Reader, modifiers ...ModifyRequest) (*http.Response, error) {
	absUrl := self.GetBaseUrl() + url

	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		absUrl = url
	}

	req, err := http.NewRequestWithContext(ctx, method, absUrl, body)
	if err != nil {
		return nil, fmt.Errorf(`NewRequest:%s`, err.Error())
	}
	if len(modifiers) == 0 {
		self.AttachHeaders(req, url)
	}

	if len(modifiers) > 0 {
		for _, modFn := range modifiers {
			modFn(req)
		}
	}

	return self.httpclient.Do(req)
}

//GetBytes a GET method with a return type []byte
func (self *Client) PostBodyReader(ctx context.Context, url string, post map[string]interface{}, modFns ...ModifyRequest) (io.ReadCloser, error) {
	body, err := json.Marshal(post)
	if err != nil {
		return nil, err
	}
	resp, err := self.NewRequestWithContext(ctx, `POST`, url, bytes.NewReader(body), modFns...)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {

		return nil, fmt.Errorf(`bad status code:%d`, resp.StatusCode)
	}
	return resp.Body, nil

}

func (self *Client) PostBytes(ctx context.Context, url string, post map[string]interface{}, modFns ...ModifyRequest) ([]byte, error) {
	reader, err := self.PostBodyReader(ctx, url, post, modFns...)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

//GetBytes a GET method with a return type []byte
func (self *Client) GetBodyReader(ctx context.Context, url string, modFns ...ModifyRequest) (io.ReadCloser, error) {
	resp, err := self.NewRequestWithContext(ctx, `GET`, url, nil, modFns...)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 300 {

		return nil, fmt.Errorf(`bad status code:%d`, resp.StatusCode)
	}
	return resp.Body, nil

}

//GetBytes a GET method with a return type []byte
func (self *Client) GetBytes(ctx context.Context, url string, modFns ...ModifyRequest) ([]byte, error) {
	reader, err := self.GetBodyReader(ctx, url, modFns...)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

//GetMsi a GET method with a return type map[string]interface{}
func (self *Client) GetMsi(ctx context.Context, url string) (map[string]interface{}, error) {
	//!!! no modFns passed in because set Content-Type:applicaon/json
	body, err := self.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	if err := json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

//GetFileBytes expect read entire content into memory.
func (self *Client) GetFileBytes(ctx context.Context, fileId string) ([]byte, error) {
	url := fmt.Sprintf(`/v2/files/%s/content`, fileId)
	removeContentTypeJson := func(r *http.Request) {
		self.AttachHeaders(r, url)
		if t := r.Header.Get(CONTENT_TYPE); t != "" {
			r.Header.Del(CONTENT_TYPE)
		}
	}
	return self.GetBytes(ctx, url, removeContentTypeJson)
}

//GetFileBytes expect read entire content into memory.
func (self *Client) GetFileReader(ctx context.Context, fileId string) (io.ReadCloser, error) {
	url := fmt.Sprintf(`/v2/files/%s/content`, fileId)
	//the file content request will fail with 403 if there is application/json Content-Type
	removeContentTypeJson := func(r *http.Request) {
		self.AttachHeaders(r, url)
		if t := r.Header.Get(CONTENT_TYPE); t != "" {
			r.Header.Del(CONTENT_TYPE)
		}

	}
	return self.GetBodyReader(ctx, url, removeContentTypeJson)
}
