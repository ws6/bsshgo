package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
)

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
		Data []struct {
			Lane     string
			SampleID string
			Index    string `json:"index"`
			Index2   string `json:"index2"`
			//the below three attributes belong to Name="Cloud" specificly
			Sample_ID   string
			ProjectName string
			LibraryName string
		}
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
