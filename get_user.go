package bsshgo

//retrieve more information from user end point /v1pre3/users/current

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type UserResp struct {
	Response struct {
		HrefHistory string
		HrefDomain  string
		IsWorkgroup bool
		Id          string
		Name        string
		DomainName  string
	}
}

func (self *Client) GetCurrentUser(ctx context.Context) (*UserResp, error) {
	url := `/v1pre3/users/current`
	resp, err := self.NewRequestWithContext(ctx, `GET`, url, nil)
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

	ret := new(UserResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, fmt.Errorf(`Unmarshal:%s`, err.Error())
	}
	return ret, nil
}
