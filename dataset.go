package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

//dataset.go accessor for datasets

func (self *Client) GetDatasetsFiles(ctx context.Context, dsId string, params map[string]string) (*FileResp, error) {
	_url := fmt.Sprintf(`/v2/datasets/%s/files`, dsId)
	q := url.Values{}

	if params != nil {
		for k, v := range params {
			q.Add(k, v)
		}
	}
	base, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	base.RawQuery = q.Encode()

	body, err := self.GetBytes(ctx, base.String())
	if err != nil {
		return nil, err
	}
	ret := new(FileResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (self *Client) GetDatasetsFilesChan(ctx context.Context, dsId string) chan *FileItem {
	limit := 20
	ret := make(chan *FileItem, 3*limit)
	params := make(map[string]string)
	go func() {
		defer close(ret)
		offset := 0

		for {
			params[`limit`] = fmt.Sprintf(`%d`, limit)
			params[`offset`] = fmt.Sprintf(`%d`, offset)
			offset += limit
			page, err := self.GetDatasetsFiles(ctx, dsId, params)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, found := range page.Items {
				ret <- found
				select {
				case <-ctx.Done():
					return
				default:
					continue
				}
			}

			if page.Paging.DisplayedCount < limit || page.Paging.DisplayedCount == 0 {
				return
			}
		}

	}()
	return ret

}
