package bsshgo

//one way to dig into analysis to get biosamples out.
//it has varies of limitations. we assumed the path way follows:
//analysis->fastq dataset->biosample/run
//be aware this may not the same setup as yours

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

//bssh_analysis_to_matcha_analysis_from_fastq.go generate matcha spec from fastq dataset to align with Fastqc.
//if fastqc failed, it shall not be presented in the analysis by that time.
//be aware fastqc status updates over time.

type DataSetItem struct {
	DataSet struct {
		Id        string
		TotalSize int64
		QcStatus  string
	}
}

type DataSetLayout struct {
	LaneNumber  int
	RunId       string
	SampleId    string
	ProjectName string //default projectname
}

//GetFastqDataSetLayout a dataset may or may not have all layout attributes. it depends on its attributes
// func GetFastqDataSetLayout(ctx context, client *bsshgo.Client, dsId string) ([]*DataSetLayout, error) {
// 	return nil, nil
// }

func (client *Client) GetItems(ctx context.Context, url string, params map[string]string,
	// makeObjFn func() interface{},
	receiveFn func(interface{}) error) error {
	ch, err := client.GetGeneralItemsChannel(
		ctx,
		url,
		params,
	)
	if err != nil {
		return err
	}

	for found := range ch {

		if found.Err != nil {
			return fmt.Errorf(`found.Err:%s`, found.Err.Error())
		}

		if err := receiveFn(found.Item); err != nil {
			return err
		}
	}

	return nil

}

type LaneNumContentMap struct {
	DateModified string
	Type         string
	Href         string
	Name         string
	ContentMap   []struct {
		Key    string
		Values []string
	}
}

func (client *Client) GetFastqDataSetLaneNumbers(ctx context.Context, dsId string) ([]int, error) {
	url := fmt.Sprintf(`/v2/datasets/%s/properties/BaseSpace.Metrics.FastQ`, dsId)
	body, err := client.GetBytes(ctx, url)
	if err != nil {
		return nil, err
	}
	resp := new(LaneNumContentMap)
	if err := json.Unmarshal(body, resp); err != nil {
		return nil, err
	}
	ret := []int{}

	for _, ctm := range resp.ContentMap {
		if ctm.Key != `LaneNumber` { //hacky
			continue
		}
		for _, v := range ctm.Values {
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf(`LaneNumber:%s`, err.Error())
			}
			ret = append(ret, i)
		}
	}

	return ret, nil
}

type BioSampleItem struct {
	BioSample struct {
		Id             string
		BioSampleName  string
		DefaultProject struct {
			Name string
		}
	}
}

func (client *Client) GetFastqDataSetInputBioSamples(ctx context.Context, dsId string) ([]*BioSampleItem, error) {
	url := fmt.Sprintf(`/v2/datasets/%s/properties/Input.BioSamples/Items`, dsId)
	ret := []*BioSampleItem{}
	receiveFn := func(i interface{}) error {
		b, err := json.Marshal(i)
		if err != nil {
			return err
		}
		topush := new(BioSampleItem)
		if err := json.Unmarshal(b, topush); err != nil {
			return err
		}
		ret = append(ret, topush)
		return nil
	}
	if err := client.GetItems(
		ctx,

		url,
		map[string]string{},
		receiveFn,
	); err != nil {
		return nil, fmt.Errorf(`GetItems:%s`, err.Error())
	}
	return ret, nil
}

type RunItem struct {
	Id  string
	Run struct {
		Id   string
		Name string
	}
}

func (client *Client) GetFastqDataSetRuns(ctx context.Context, dsId string) ([]*RunItem, error) {
	url := fmt.Sprintf(`/v2/datasets/%s/properties/Input.Runs/Items`, dsId)
	ret := []*RunItem{}
	receiveFn := func(i interface{}) error {

		b, err := json.Marshal(i)
		if err != nil {
			return err
		}
		topush := new(RunItem)
		if err := json.Unmarshal(b, topush); err != nil {
			return err
		}
		ret = append(ret, topush)
		return nil
	}
	if err := client.GetItems(
		ctx,

		url,
		map[string]string{},
		receiveFn,
	); err != nil {
		return nil, fmt.Errorf(`GetItems:%s`, err.Error())
	}
	return ret, nil
}

type FastqDataSetLayout struct {
	*BioSampleItem
	*RunItem
	LaneNumber int
}

//assume only one Run associate with fastq dataset
//assume only one biosample associate with fastq type dataset.
//it can have multiple lanes
func (client *Client) GetFastqDataSetLayout(ctx context.Context, dsId string) ([]*FastqDataSetLayout, error) {
	biosamples, err := client.GetFastqDataSetInputBioSamples(ctx, dsId)
	if err != nil {
		return nil, fmt.Errorf(`GetFastqDataSetInputBioSamples:%s`, err.Error())
	}
	if len(biosamples) == 0 {
		return nil, fmt.Errorf(`no biosample`)
	}
	runs, err := client.GetFastqDataSetRuns(ctx, dsId)
	if err != nil {
		return nil, fmt.Errorf(`GetFastqDataSetRuns:%s`, err.Error())
	}
	if len(runs) == 0 {
		return nil, fmt.Errorf(`no runs`)
	}
	laneNums, err := client.GetFastqDataSetLaneNumbers(ctx, dsId)
	if err != nil {
		return nil, fmt.Errorf(`GetFastqDataSetLaneNumbers:%s`, err.Error())
	}
	if len(laneNums) == 0 {
		return nil, fmt.Errorf(`no laneNumbers`)
	}
	ret := []*FastqDataSetLayout{}
	for _, ln := range laneNums {
		topush := new(FastqDataSetLayout)
		topush.BioSampleItem = biosamples[0]
		topush.RunItem = runs[0]
		topush.LaneNumber = ln
		ret = append(ret, topush)
	}

	return ret, nil
}

func (self *Client) GetBioSamplesFromAnalysisThroughFastqDatasetUsed(ctx0 context.Context, appsessionId string) ([]*BioSampleItem, error) {
	ctx, cancelFn := context.WithCancel(ctx0) //making sure all early return will cancel the channel
	defer cancelFn()

	dsItems := []*DataSetItem{}

	receiveFn := func(i interface{}) error {
		b, err := json.Marshal(i)
		if err != nil {
			return err
		}
		topush := new(DataSetItem)
		if err := json.Unmarshal(b, topush); err != nil {
			return err
		}
		dsItems = append(dsItems, topush)
		return nil
	}

	if err := self.GetItems(
		ctx,

		fmt.Sprintf(`/v2/appsessions/%s/properties/Input.automation-sample-id.datasets/items`, appsessionId),
		map[string]string{},
		receiveFn,
	); err != nil {
		return nil, fmt.Errorf(`GetItems:%s`, err.Error())
	}
	ret := []*FastqDataSetLayout{}
	for _, fastqDs := range dsItems {

		res, err := self.GetFastqDataSetLayout(ctx, fastqDs.DataSet.Id)
		if err != nil {
			return nil, fmt.Errorf(`GetFastqDataSetLayout:%s`, err.Error())
		}
		ret = append(ret, res...)
	}

	//return unique biosample
	m := make(map[string]*BioSampleItem)
	for _, b := range ret {
		if _, ok := m[b.BioSample.Id]; !ok {
			m[b.BioSample.Id] = b.BioSampleItem
		}
	}
	ret2 := []*BioSampleItem{}
	for _, v := range m {
		ret2 = append(ret2, v)
	}
	return ret2, nil
}
