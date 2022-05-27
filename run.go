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
		if cloudData == nil {
			return fmt.Errorf(`Cloud data is empty`)
		}
		for _, d := range cloudData {
			if d.Sample_ID == r.Sample_ID {
				r.ProjectName = d.ProjectName
				r.LibraryName = d.LibraryName
				return nil
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

//GetRunHref this is being used as GUID for a run entity; it could be used as client side identity
//return like https://api.basespace.illumina.com/v2/runs/219700482
func (self *Client) GetRunHref(runId string) string {
	return fmt.Sprintf(`%s/v2/runs/%s`, self.GetBaseUrl(), runId)
}

type LaneLibraryMappingsItem struct {
	LaneId        string
	LaneNumber    int
	BioSampleName string //finger crossed there is no difference from BioSample.BioSample
	LibraryName   string
	BioSample     BioSampleItem
	Library       LibraryPrepItem
	LibraryPool   struct {
		Id           string
		Href         string
		Name         string
		LibraryCount int
		DateModified string
		DateCreated  string
		Status       string
	}
	DataSetName    string
	DatasetYieldBp int64
	ProjectName    string
}

func (self *Client) GetRunLaneLibraryMappingsItem(ctx context.Context, runId string) ([]*LaneLibraryMappingsItem, error) {
	_url := fmt.Sprintf(`/v2/runs/%s/lanelibrarymappings`, runId)
	ch, err := self.GetGeneralItemsChannel(ctx, _url, nil)
	if err != nil {
		return nil, err
	}
	ret := []*LaneLibraryMappingsItem{}
	for item := range ch {
		if item.Err != nil {
			return nil, err
		}
		body, err := json.Marshal(item.Item)
		if err != nil {
			return nil, err
		}
		topush := new(LaneLibraryMappingsItem)
		if err := json.Unmarshal(body, topush); err != nil {
			return nil, err
		}
		ret = append(ret, topush)
	}

	return ret, nil
}
func (self *Client) GetRunLayoutFromLanelibrarymappings(ctx context.Context, runId string) ([]*RunLayoutResp, error) {

	getLibPool := func() func(poolId string) ([]*LibraryPoolItem, error) {
		cache := make(map[string][]*LibraryPoolItem)
		return func(poolId string) ([]*LibraryPoolItem, error) {
			found, ok := cache[poolId]
			if ok {
				return found, nil
			}
			ret, err := self.GetLibraryPoolInfo(ctx, poolId)
			if err != nil {
				return nil, err
			}
			cache[poolId] = ret

			return ret, nil
		}
	}()
	getLibraryItem := func(poolId, libraryName string) (*LibraryPoolItem, error) {
		pool, err := getLibPool(poolId)
		if err != nil {
			return nil, fmt.Errorf(`no pool`)
		}
		for _, item := range pool {
			if item.Name == libraryName {
				return item, nil
			}
		}
		return nil, fmt.Errorf(`no pool item`)
	}

	laneMappings, err := self.GetRunLaneLibraryMappingsItem(ctx, runId)

	if err != nil {
		return nil, fmt.Errorf(`GetRunLaneLibraryMappingsItem:%s`, err.Error())
	}
	if len(laneMappings) == 0 {
		return nil, fmt.Errorf(`no laneMappings`)
	}
	ret := []*RunLayoutResp{}
	for _, item := range laneMappings {
		topush := new(RunLayoutResp)
		topush.Sample_ID = item.BioSampleName
		topush.ProjectName = item.ProjectName
		topush.LibraryName = item.LibraryName
		libpoolitem, err := getLibraryItem(item.LibraryPool.Id, topush.LibraryName)
		if err != nil {
			return nil, fmt.Errorf(`getLibraryItem:%s`, err.Error())
		}
		topush.Lane = fmt.Sprintf("%d", item.LaneNumber)
		topush.Index = libpoolitem.Index1Sequence
		topush.Index2 = libpoolitem.Index2Sequence
		ret = append(ret, topush)
	}

	return ret, nil
}

func (self *Client) GetRunLayoutFromV2SampleSheetAPI(ctx context.Context, runId string) ([]*RunLayoutResp, error) {
	lay1, err := self.GetRunSampleSheetLayout(ctx, runId)
	if err != nil {
		return nil, fmt.Errorf(`GetRunSampleSheetLayout:%s`, err.Error())
	}
	return ConcateRunLayoutWithBCLConvertAndCloudApplicaions(lay1), nil
}

func (self *Client) GetRunLayout(ctx context.Context, runId string) ([]*RunLayoutResp, error) {
	ret1, err1 := self.GetRunLayoutFromV2SampleSheetAPI(ctx, runId)
	if err1 != nil {
		fmt.Println(`GetRunLayoutFromV2SampleSheetAPI not working. trying GetRunLayoutFromLanelibrarymappings`)
		return self.GetRunLayoutFromLanelibrarymappings(ctx, runId)
	}
	return ret1, nil
}

type RunDetailsResp struct {
	Id                  string
	ExperimentName      string
	DateCreated         string
	DateModified        string
	Status              string
	InstrumentRunStatus string
	FlowcellPosition    string
	LaneAndQcStatus     string
	Workflow            string
	V1Pre3Id            string
	Instrument          struct {
		Id           string
		Name         string
		Type         string
		PlatformName string
	}

	UploadStatus        string
	DateUploadStarted   string
	DateUploadCompleted string
}

func (self *Client) GetRun(ctx context.Context, runId string) (*RunDetailsResp, error) {
	_url := fmt.Sprintf(`/v2/runs/%s`, runId)
	body, err := self.GetBytes(ctx, _url)
	if err != nil {
		return nil, err
	}
	ret := new(RunDetailsResp)
	if err := json.Unmarshal(body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
