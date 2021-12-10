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

//FindOneAnalysisByName assumed only one result get hit.
//NB: AnalysisResp is not 100% same spec as it is.
//actual returned
// {
//     "Items": [
//         {
//             "Id": "484415694",
//             "Name": "ILS_DRAGEN_GL_2.0.2 11/08/2021 12:41:34",
//             "Href": "https://api.basespace.illumina.com/v2/appsessions/484415694",
//             "Application": {
//                 "Id": "12760748",
//                 "Href": "v1pre3/applications/12760748",
//                 "Name": "ILS_DRAGEN_GL_2.0.2",
//                 "CompanyName": "2e676b91-a9d1-4518-8f40-b8ec...",
//                 "VersionNumber": "1.0.0",
//                 "ShortDescription": "'ILS",
//                 "DateCreated": "2021-10-29T10:55:16.0000000Z",
//                 "PublishStatus": "Beta",
//                 "IsBillingActivated": false,
//                 "Category": "Workflow",
//                 "Classifications": [],
//                 "AppFamilySlug": "2e676b91-a9d1-4518-8f40-b8ec.ils_dragen_gl_2.0.2",
//                 "AppVersionSlug": "2e676b91-a9d1-4518-8f40-b8ec.ils_dragen_gl_2.0.2.1.0.0",
//                 "Features": [],
//                 "LockStatus": "Unlocked"
//             },
//             "UserCreatedBy": {
//                 "Id": "27392365",
//                 "Href": "https://api.basespace.illumina.com/v2/users/27392365",
//                 "Name": "ILS_SD_NovaSeq_Integration",
//                 "DateCreated": "2021-08-09T21:10:06.0000000Z",
//                 "GravatarUrl": "https://secure.gravatar.com/avatar/8e4e5ab5bf1e5a2eeb564a5f023bb9dc.jpg?s=20&d=mm&r=PG",
//                 "HrefProperties": "https://api.basespace.illumina.com/v2/users/current/properties",
//                 "ExternalDomainId": "YXdzLXVzLXBsYXRmb3JtOjEwMDAwNzgyOjk0ZWFlN2I2LWRjNDItNDAzMS04NGYwLTQ5YTllMmUyOTVhZA"
//             },
//             "ExecutionStatus": "Running",
//             "QcStatus": "Undefined",
//             "StatusSummary": "",
//             "Purpose": "AppTrigger",
//             "DateCreated": "2021-11-08T20:43:44.0000000Z",
//             "DateModified": "2021-11-08T20:43:55.0000000Z",
//             "DateStarted": "2021-11-08T20:43:55.0000000Z",
//             "DeliveryStatus": "None",
//             "ContainsComments": false,
//             "HrefComments": "https://api.basespace.illumina.com/v2/appsessions/484415694/comments"
//         }
//     ],
//     "Paging": {
//         "DisplayedCount": 1,
//         "TotalCount": 1,
//         "Offset": 0,
//         "Limit": 10,
//         "SortDir": "Asc",
//         "SortBy": "Name"
//     }
// }
func (self *Client) FindOneAnalysisByName(ctx context.Context, name string) (*AnalysisResp, error) {
	_url := `/v2/appsessions`
	params := map[string]string{
		`name`:   name,
		`limit`:  "1",
		`offset`: "0",
	}
	base, err := url.Parse(_url)
	if err != nil {

		return nil, err
	}
	q := url.Values{}

	for k, v := range params {
		q.Add(k, v)
	}

	base.RawQuery = q.Encode()

	body, err := self.GetBytes(ctx, base.String())
	if err != nil {

		return nil, err
	}

	resp := new(GeneralPropertiesResp)
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	if resp.Paging.TotalCount == 0 || len(resp.Items) == 0 {
		return nil, fmt.Errorf(`not founds`)
	}
	if resp.Paging.TotalCount > 1 || len(resp.Items) > 1 {
		return nil, fmt.Errorf(`multiple results returned`)
	}

	b, err := json.Marshal(resp.Items[0])
	if err != nil {
		return nil, err
	}
	ret := new(AnalysisResp)
	if err := json.Unmarshal(b, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type DatasetItemResp struct {
	Id      string //dont know what is this
	Dataset struct {
		Id           string //the actual dataset Id
		Href         string
		HrefFiles    string
		Name         string
		DateCreated  string
		DateModified string
		AppSession   struct {
			Id   string
			Name string
			Href string
			ApplicationResp
		}
		Project struct {
			Id   string
			Name string
			Href string
		}
		TotalSize           int64
		IsArchived          bool
		Attributes          map[string]interface{}
		QcStatus            string
		UploadStatus        string
		UploadStatusSummary string
		ValidationStatus    string
		V1pre3Id            string
		HrefComments        string
		ContainsComments    bool
	}
}

func (self *Client) GetAnalysisOutputDatasetChan(ctx context.Context, analysisId string) (chan *DatasetItemResp, error) {
	url := fmt.Sprintf(`/v2/appsessions/%s/properties/Output.Datasets/items`, analysisId)

	params := make(map[string]string)

	limit := 30

	params[`limit`] = fmt.Sprintf(`%d`, limit)

	ch, err := self.GetGeneralItemsChannel(ctx, url, params)
	if err != nil {
		return nil, err
	}

	ret := make(chan *DatasetItemResp, 2*cap(ch))
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
			topush := new(DatasetItemResp)
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

func (self *Client) UpdateAnalysis(ctx context.Context, analysisId string, updates map[string]interface{}) (*AnalysisResp, error) {
	_url := fmt.Sprintf(`/v2/appsessions/%s`, analysisId)
	body, err := self.PostBytes(ctx, _url, updates)
	if err != nil {
		return nil, err
	}
	ret := new(AnalysisResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
