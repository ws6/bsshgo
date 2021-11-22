package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

//analysis.go entity model retrieval for analysis(a.k.a appsession)

type ApplicationResp struct {
	Id            string
	Href          string
	Name          string
	VersionNumber string
}

type GeneralPropetyItem struct {
	Type        string
	Href        string
	Name        string
	Description string
	Content     string
}

type AnalysisPropertiesResp struct {
	Items  []*GeneralPropetyItem
	Paging struct {
		DisplayedCount int
		TotalCount     int
		Offset         int
		Limit          int
		SortDir        string
		SortBy         string
	}
}
type GeneralPropertiesResp struct {
	Items  []map[string]interface{}
	Paging struct {
		DisplayedCount int
		TotalCount     int
		Offset         int
		Limit          int
		SortDir        string
		SortBy         string
	}
}

type AnalysisResp struct {
	RunningDuration int64
	Id              string
	Name            string
	Href            string
	Application     *ApplicationResp
	ExecutionStatus string
	QcStatus        string
	StatusSummary   string
	DateCreated     string
	DateModified    string
	DateCompleted   string
	DateStarted     string
	TotalSize       int64
	Properties      struct {
		Items          []*GeneralPropetyItem
		DisplayedCount int
		TotalCount     int
		Href           string
	}

	DeliveryStatus   string
	ContainsComments bool
	HrefComments     string
}

func (self *Client) GetAnalysisById(ctx context.Context, Id string) (*AnalysisResp, error) {
	url := fmt.Sprintf(`/v2/appsessions/%s`, Id)
	body, err := self.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}
	ret := new(AnalysisResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type ChannelResp struct {
	Item interface{}
	Err  error
}

// GeneralPropertiesResp

func (self *Client) GetAnalysisPropertiesChannel(ctx context.Context, analysisId, propertyfilters string) (chan *ChannelResp, error) {
	_url := fmt.Sprintf(`/v2/appsessions/%s/properties`, analysisId)
	params := make(map[string]string)
	limit := 30
	params[`limit`] = fmt.Sprintf(`%d`, limit)
	params[`propertyfilters`] = propertyfilters
	ch, err := self.GetGeneralItemsChannel(ctx, _url, params)
	if err != nil {
		return nil, err
	}
	ch2 := make(chan *ChannelResp, 2*limit)
	toChan := func(item interface{}, err error) {
		fmt.Println(item)
		topush := new(ChannelResp)
		body, err := json.Marshal(item)
		if err != nil {
			topush.Err = err
		}

		if err == nil {
			analysisPropertyItem := new(GeneralPropetyItem)
			err2 := json.Unmarshal(body, analysisPropertyItem)
			if err2 != nil {
				topush.Err = err
			}
			if err2 == nil {
				topush.Item = analysisPropertyItem
			}
		}

		// topush.Item = item
		if err != nil {
			topush.Err = err
		}

		select {
		case <-ctx.Done():
			return
		case ch2 <- topush:
			return

		}
	}
	func() {
		defer close(ch2)
		for found := range ch {
			if found.Err != nil {
				toChan(nil, found.Err)
				return //append err
			}

			toChan(found.Item, nil)

		}
	}()

	return ch2, nil
}

type IsGeneralPropertyItem func(*GeneralPropetyItem) bool

func (self *Client) GetAnalysisPropertyUntil(ctx context.Context, analysisId, propertyfilters string, toSearchFn IsGeneralPropertyItem) (*GeneralPropetyItem, error) {

	childCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	ch, err := self.GetAnalysisPropertiesChannel(childCtx, analysisId, propertyfilters)
	if err != nil {
		return nil, err
	}

	for found := range ch {
		item, ok := found.Item.(*GeneralPropetyItem)
		if !ok {
			return nil, fmt.Errorf(`item is not AnalysisPropetyItem`)
		}

		if found.Err != nil {
			return nil, found.Err
		}

		if toSearchFn(item) {
			return item, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			continue
		}
	}

	return nil, ERR_NOT_FOUND
}

//GetItemsChannel it aims for in general cases a URl with Params and support limit/offset in  Params
func (self *Client) GetGeneralItemsChannel(ctx context.Context, _url string, params map[string]string) (chan *ChannelResp, error) {

	limit := 30
	offset := 0

	limitStr, ok := params[`limit`]
	if ok {
		if n, err := strconv.Atoi(limitStr); err == nil && n > 0 {
			limit = n
		}
	}
	ret := make(chan *ChannelResp, 2*limit)

	params[`limit`] = fmt.Sprintf(`%d`, limit)
	toChan := func(item map[string]interface{}, err error) {
		topush := new(ChannelResp)
		topush.Item = item
		if err != nil {
			topush.Err = err
		}

		select {
		case <-ctx.Done():
			return
		case ret <- topush:
			return

		}
	}
	go func() {
		defer close(ret)
		for {
			base, err := url.Parse(_url)
			if err != nil {
				toChan(nil, err)
				return
			}
			q := url.Values{}

			for k, v := range params {
				q.Add(k, v)
			}

			base.RawQuery = q.Encode()

			body, err := self.GetBytes(ctx, base.String())
			if err != nil {
				toChan(nil, err)
				return
			}
			offset += limit
			params[`offset`] = fmt.Sprintf(`%d`, offset)

			resp := new(GeneralPropertiesResp)
			if err := json.Unmarshal(body, resp); err != nil {
				toChan(nil, err)
				return
			}

			for _, item := range resp.Items {

				toChan(item, nil)
			}
			if resp.Paging.DisplayedCount < limit || resp.Paging.DisplayedCount < resp.Paging.Limit {

				return
			}
		}
	}()

	return ret, nil
}

//GetGeneralPropertyItemUntil the simple goal is to retrieve a Href
func (self *Client) GetGeneralPropertyItemUntil(ctx context.Context, url string, params map[string]string, toSearchFn IsGeneralPropertyItem) (*GeneralPropetyItem, error) {

	childCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	ch, err := self.GetGeneralItemsChannel(childCtx, url, params)
	if err != nil {
		return nil, err
	}

	for found := range ch {

		body, err := json.Marshal(found.Item)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}

		item := new(GeneralPropetyItem)
		if err := json.Unmarshal(body, item); err != nil {
			return nil, fmt.Errorf(`item is not GeneralPropetyItem`)
		}

		if found.Err != nil {
			return nil, found.Err
		}

		if toSearchFn(item) {
			return item, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			continue
		}
	}

	return nil, ERR_NOT_FOUND
}

//BioSampleItemResp
//https://api.basespace.illumina.com/v2/appsessions/486869383/properties/Input.BioSamples/items
type BioSampleItemResp struct {
	Id string

	BioSample struct {
		Id             string
		Href           string
		BioSampleName  string
		DefaultProject struct {
			Id   string
			Href string //most likely v1pre3/projects/:Id
			Name string
		}
		DateModified      string
		DateCreated       string
		Status            string
		LabStatus         string
		ContainerName     string
		ContainerPosition string
	}
}

func (self *Client) GetAnalysisBioSamples(ctx context.Context, analysisId string) (chan *BioSampleItemResp, error) {
	url := fmt.Sprintf(`/v2/appsessions/%s/properties/Input.BioSamples/items`, analysisId)

	params := make(map[string]string)

	limit := 30

	params[`limit`] = fmt.Sprintf(`%d`, limit)

	ch, err := self.GetGeneralItemsChannel(ctx, url, params)
	if err != nil {
		return nil, err
	}

	ret := make(chan *BioSampleItemResp, 2*limit)
	go func() {
		defer close(ret)

		for found := range ch {
			if found.Err != nil {
				return
			}

			body, err := json.Marshal(found.Item)
			if err != nil {
				return
			}
			topush := new(BioSampleItemResp)
			if err := json.Unmarshal(body, topush); err != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case ret <- topush:
				continue
			}
		}

	}()
	return ret, nil
}