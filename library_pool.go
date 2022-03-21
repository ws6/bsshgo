package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
)

type LibraryPoolItem struct {
	BioSampleName string
	ProjectName   string
	Id            string
	Href          string
	Name          string
	DataCreated   string
	DateModified  string
	Status        string
	BioSample     BioSampleItem
	LibraryPrep   LibraryItemResp
	Project       struct {
		Id   string
		Name string
		Href string
	}
	Biomolecule    string
	Index1Sequence string
	Index2Sequence string
}

//library_pool.go fetch library pool entities
// https://developer.basespace.illumina.com/docs/content/documentation/rest-api/api-reference#operation--librarypools--id--libraries-get
// https://api.basespace.illumina.com/v2/librarypools/199512489/libraries?Include=LibraryIndex&limit=25&offset=0&sortBy=DateCreated&sortDir=desc

func (self *Client) GetLibraryPoolInfo(ctx context.Context, libraryPoolId string) ([]*LibraryPoolItem, error) {
	_url := fmt.Sprintf(`/v2/librarypools/%s/libraries`, libraryPoolId)
	params := map[string]string{
		`Include`: `LibraryIndex`,
	}
	ch, err := self.GetGeneralItemsChannel(ctx, _url, params)
	if err != nil {
		return nil, err
	}
	ret := []*LibraryPoolItem{}
	for item := range ch {
		if item.Err != nil {
			return nil, item.Err
		}
		body, err := json.Marshal(item.Item)
		if err != nil {
			return nil, err
		}
		topush := new(LibraryPoolItem)

		if err := json.Unmarshal(body, topush); err != nil {
			return nil, err
		}
		ret = append(ret, topush)

	}
	return ret, nil
}
