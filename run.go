package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
)

type RunLayoutResp struct {
	Lane string

	Index  string `json:"index"`
	Index2 string `json:"index2"`
	//the below three attributes belong to Name="Cloud" specificly
	Sample_ID   string
	ProjectName string
	LibraryName string
}

//run.go retrieve run models
//RunSampleSheetLayoutResp no pagination support
type RunSampleSheetLayoutResp struct {
	NeedsAttention bool
	FormatVersion  string
	HeaderSettings []struct {
		Key   string
		Value string
	}

	ReadSettings []struct {
		Key   string
		Value string
	}
	Applications []struct {
		Name     string
		Settings []struct {
			Key   string
			Value string
		}
		Data []*RunLayoutResp
		// Data []struct {
		// 	Lane string

		// 	Index  string `json:"index"`
		// 	Index2 string `json:"index2"`
		// 	//the below three attributes belong to Name="Cloud" specificly
		// 	Sample_ID   string
		// 	ProjectName string
		// 	LibraryName string
		// }
	}
}

func (self *Client) GetRunSampleSheetLayout(ctx context.Context, runId string) (*RunSampleSheetLayoutResp, error) {
	url := fmt.Sprintf(`/v2/runs/%s/samplesheet`, runId)
	body, err := self.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}
	ret := new(RunSampleSheetLayoutResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func ConcateRunLayoutWithBCLConvertAndCloudApplicaions(resp *RunSampleSheetLayoutResp) []*RunLayoutResp {
	ret := []*RunLayoutResp{}

	getAppData := func(s string) []*RunLayoutResp {
		for _, app := range resp.Applications {
			if app.Name == s {
				return app.Data
			}
		}
		return nil
	}

	bclConvertData := getAppData(`BCLConvert`)
	cloudData := getAppData(`Cloud`)
	getLayoutMore := func(r *RunLayoutResp) error {
		for _, d := range cloudData {
			if d.Sample_ID == r.Sample_ID {
				r.ProjectName = d.ProjectName
				r.LibraryName = d.LibraryName
			}
		}
		return ERR_NOT_FOUND
	}

	for _, d := range bclConvertData {
		if err := getLayoutMore(d); err != nil {
			continue
		}
		ret = append(ret, d)
	}

	return ret
}
