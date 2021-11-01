package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
)

//implementation for
//https://developer.basespace.illumina.com/docs/content/documentation/rest-api/history-api-reference#HistoryAPIReference

//support the token is history token with scope=AUDIT USER

type HistoryResp struct {
	Items  []map[string]interface{}
	Paging struct {
		// "TotalCount": 1247,
		TotalCount int64
		//      "DisplayedCount": 10,
		DisplayedCount int64
		//      "Limit": 10,
		Limit int64
		//      "SortBy": "DateCreated",
		SortBy string
		//      "SortDir": "desc",
		SortDir string
		//      "After": "184181146749957033",
		After string
		//      "Before": "184181438804378115"
		Before string
	}
}

func (self *Client) SearchHistory(ctx context.Context, params map[string]string) (*HistoryResp, error) {
	if self.User == nil {
		user, err := self.GetCurrentUser(ctx)
		if err != nil {
			return nil, fmt.Errorf(`GetCurrentUser:%s`, err.Error())
		}
		self.User = user
	}
	_url := self.User.Response.HrefHistory

	if _url == "" {
		return nil, fmt.Errorf(`self.User.Response.HrefHistory is empty`)
	}

	base, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	q := url.Values{}

	if params != nil {
		for k, v := range params {
			q.Add(k, v)
		}
	}

	base.RawQuery = q.Encode()
	fmt.Println(base.String())
	resp, err := self.NewRequestWithContext(ctx, `GET`, base.String(), nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(`bad status code-%d:%s`, resp.StatusCode, string(body))
	}

	ret := new(HistoryResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil

}
