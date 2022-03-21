package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
)

//biosample.go biosample entity retrieval

type BioSampleResp struct {
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
	ContainerName     string
	ContainerPosition string
	Status            string
	LabStatus         string
	Properties        struct {
		Items        []*GeneralPropetyItem
		DisplayCount int
		TotalCount   int
		Href         string
	}
}

func (self *Client) GetBioSample(ctx context.Context, Id string) (*BioSampleResp, error) {
	body, err := self.GetBytes(ctx,
		fmt.Sprintf(`/v2/biosamples/%s`, Id),
	)

	if err != nil {
		return nil, err
	}
	ret := new(BioSampleResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// https://api.basespace.illumina.com/v2/biosamples/465598133/properties?propertyfilters=Metadata.Concentration
func (self *Client) GetAllBioSampleProperties(ctx context.Context, sampleId, propertyfilters string) (chan *ChannelResp, error) {

	params := make(map[string]string)
	if propertyfilters != "" {
		params[`propertyfilters`] = propertyfilters
	}
	limit := 30

	params[`limit`] = fmt.Sprintf(`%d`, limit)
	url := fmt.Sprintf(`/v2/biosamples/%s/properties`, sampleId)
	return self.GetGeneralItemsChannel(ctx, url, params)
}

// GetGeneralPropertyItemUntil
func (self *Client) GetBioSamplePropertyItemUntil(ctx context.Context, sampleId string, params map[string]string, fn IsGeneralPropertyItem) (*GeneralPropetyItem, error) {

	limit := 30

	params[`limit`] = fmt.Sprintf(`%d`, limit)
	url := fmt.Sprintf(`/v2/biosamples/%s/properties`, sampleId)
	return self.GetGeneralPropertyItemUntil(ctx, url, params, fn)
}

type LibraryPrepItem struct {
	Id                      string
	Href                    string
	Name                    string
	ValidIndexingStrategies string
	ValidReadTypes          string
	NumIndexCycles          int
	AdapterSequenceRead1    string
	AdapterSequenceRead2    string
	DateModified            string
	State                   string
	DefaultRead1Cycles      int
	DefaultRead2Cycles      int
	LibraryType             string
}
type LibraryItemResp struct {
	Id           string
	Href         string
	Name         string
	DateCreated  string
	DateModified string
	Status       string
	LibraryPrep  LibraryPrepItem
	Pools        []struct {
		Id   string
		Href string

		UserPoolId   string
		LibraryCount int
		DateModified string
		DateCreated  string
		Status       string
		LibraryPrep  []string
	}
	Biomolecule string
}

func (self *Client) GetAllBioSampleLibraries(ctx context.Context, sampleId string) (chan *LibraryItemResp, error) {

	params := make(map[string]string)

	limit := 30

	params[`limit`] = fmt.Sprintf(`%d`, limit)
	url := fmt.Sprintf(`/v2/biosamples/%s/libraries`, sampleId)
	ch, err := self.GetGeneralItemsChannel(ctx, url, params)
	if err != nil {
		return nil, err
	}

	ret := make(chan *LibraryItemResp, 2*limit)
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
			topush := new(LibraryItemResp)
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

type RunResp struct {
	RunId           string
	LaneNumber      int
	LibraryName     string
	RunStatus       string
	Workflow        string
	FlowcellBarcode string
	FastQStatus     string
	ExperimentName  string
	DateModified    string
}

// https://api.basespace.illumina.com/v2/biosamples/465598133/runlanesummaries?offset=0&limit=25&sortDir=Desc&sortBy=DateCreated

//getBiosample runs
func (self *Client) GetBioSampleRuns(ctx context.Context, sampleId string) (chan *RunResp, error) {
	params := make(map[string]string)

	limit := 30

	params[`limit`] = fmt.Sprintf(`%d`, limit)
	url := fmt.Sprintf(`/v2/biosamples/%s/runlanesummaries`, sampleId)
	ch, err := self.GetGeneralItemsChannel(ctx, url, params)
	if err != nil {
		return nil, err
	}

	ret := make(chan *RunResp, 2*limit)
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
			topush := new(RunResp)
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
