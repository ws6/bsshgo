package bsshgo

//get_run_file.go an enhancement of trying to reach out actual file contents

// https://api.basespace.illumina.com/v2/runs/219846639/files?recursive=false&filehrefcontentresolution=false&sortdir=Asc&sortby=Name&limit=2000&directory=/

//ref https://developer.basespace.illumina.com/docs/content/documentation/rest-api/api-reference#operation--runs--id--files-get

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

var ERR_NOT_FOUND = fmt.Errorf(`not found`)

type FileItem struct {
	Id           string
	Href         string
	HrefContent  string
	Name         string
	ContentType  string
	Size         int64
	Path         string
	IsArchived   bool
	DateCreated  string
	DateModified string
	ETag         string
}

type FileResp struct {
	Items  []*FileItem
	Paging struct {
		DisplayedCount int
		Offset         int
		Limit          int
		SortDir        string
		SortBy         string
	}
}

type FileFilterFn func(*FileItem) bool

func (self *Client) GetFileFromDir(pctx context.Context,

	runId string,
	params map[string]string,
	dir string,
	fn FileFilterFn,
) (*FileItem, error) {
	ctx, cancelfn := context.WithCancel(pctx)
	defer cancelfn()
	ch := self.GetAllFilesFromDir(ctx, runId, params, dir)
	for fileitem := range ch {
		if fn(fileitem) {
			return fileitem, nil
		}
	}
	return nil, ERR_NOT_FOUND
}

func (self *Client) GetAllFilesFromDir(ctx context.Context,

	runId string,
	params map[string]string,
	dir string) chan *FileItem {

	ret := make(chan *FileItem)
	go func() {
		defer close(ret)
		limit := 20
		offset := 0
		params[`limit`] = fmt.Sprintf(`%d`, limit)
		params[`offset`] = fmt.Sprintf(`%d`, offset)
		for {
			page, err := self.GetFileRespFromDir(ctx, runId, params, dir)
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

//GetFileRespByName seach one page
func (self *Client) GetFileRespFromDir(ctx context.Context,

	runId string,
	params map[string]string,
	dir string) (*FileResp, error) {
	// /v2/runs/219846639/files?recursive=false&filehrefcontentresolution=false&sortdir=Asc&sortby=Name&limit=2000&directory=/&
	// extensions=.xml
	_url := fmt.Sprintf(`/v2/runs/%s/files`, runId)
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
