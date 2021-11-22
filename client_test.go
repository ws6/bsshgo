package bsshgo

import (
	"context"
	"os"
	"testing"
)

func getConfig() map[string]string {
	ret := make(map[string]string)
	ret[AUTH_TOKEN] = os.Getenv("BSSH_TEST_AUTH_TOKEN")

	ret[BASE_URL] = os.Getenv("BSSH_TEST_BASE_URL")

	return ret
}

func getNewClient() *Client {
	ret, err := NewClient(getConfig())
	if err != nil {
		panic(err)
	}
	return ret
}

//GetAnalysisById
func TestGetAnalysisById(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// 486869383
	_biosamples, err := client.GetAnalysisBioSamples(ctx, `486869383`)
	if err != nil {
		t.Fatal(err.Error())
	}
	for bs := range _biosamples {
		t.Logf("%+v", bs)

	}
	return

	// GetRunSampleSheetLayout

	layouts, err := client.GetRunSampleSheetLayout(ctx, `219700482`)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, lo := range layouts.Applications {
		t.Logf("%+v", lo)
		for _, l2 := range lo.Data {
			t.Logf("%+v", l2)
		}
	}
	return
	// GetBioSampleRuns
	runs, err := client.GetBioSampleRuns(ctx, `470613150`)
	if err != nil {
		t.Fatal(err.Error())
	}
	for found := range runs {
		t.Logf("%+v", found)
	}
	return
	// GetAllBioSampleLibraries
	biosampleIds := []string{
		`470613148`, //has library
		`465598133`, //no libraries
	}
	for _, biosampleId := range biosampleIds {

		foundLibs, err := client.GetAllBioSampleLibraries(ctx, biosampleId)
		if err != nil {
			t.Fatal(err.Error())
		}
		for found := range foundLibs {
			t.Logf("%+v", found.LibraryPrep.LibraryType)
		}
	}

	foundSampleProp, err := client.GetBioSamplePropertyItemUntil(ctx, `465598133`,
		map[string]string{
			`limit`:           `10`,
			`propertyfilters`: `Metadata.SampleCategory`,
		},
		func(i *GeneralPropetyItem) bool {
			return i.Name == `Metadata.SampleCategory`
		})

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf(`%+v`, foundSampleProp)
	return
	sample, err := client.GetBioSample(ctx, `465598133`)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf(`%+v`, sample)

	return

	url := `/v2/appsessions/481413957/properties/Input.BioSamples/items`
	foundItems, err := client.GetGeneralItemsChannel(ctx, url, map[string]string{`limit`: `1`})
	if err != nil {
		t.Fatal(err.Error())
	}
	for item := range foundItems {
		if item.Err != nil {
			t.Fatal(err.Error())
		}
		t.Logf("%+v <<EOF>>\n", item)

	}

	return

	ret, err := client.GetAnalysisById(ctx, `486707225`)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("%+v", ret)
	tosearchPropertyNames := []string{
		`Output.YieldGbp`,
		`Input.BioSamples`,
	}

	for _, p := range tosearchPropertyNames {
		foundProperty, err := client.GetAnalysisPropertyUntil(ctx, `486707225`, p, func(i *GeneralPropetyItem) bool {
			return i.Name == p
		})
		if err != nil {
			t.Fatal(p, err.Error())
		}
		t.Logf("found foundProperty [%s]= %+v", p, foundProperty)
	}

}

func _TestGetFile(t *testing.T) {
	client := getNewClient()
	fileId := `r219846639_24790960379`
	txt, err := client.GetFileBytes(context.Background(), fileId)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(string(txt))
}

func _TestGetUser(t *testing.T) {
	client := getNewClient()

	user, err := client.GetCurrentUser(context.Background())
	if err != nil {
		t.Fatal(err.Error())
	}
	client.User = user
	t.Logf(`%+v`, client.User)
}

func _TestGetHistory(t *testing.T) {
	client := getNewClient()
	params := make(map[string]string)
	params[`SortDir`] = `Asc`
	params[`SortBy`] = `DateCreated`
	// params[`After`] = `179550990792947229`
	// params[`After`] = `179554478061316373`
	// params[`After`] = `179554490569538101`
	// params[`After`] = `179554665299426815`
	// params[`After`] = `179554725731369327`
	params[`After`] = `0`
	params[`Limit`] = `1`
	histories, err := client.SearchHistory(context.Background(), params)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf(`%+v`, histories.Paging)
	for _, item := range histories.Items {
		t.Logf(`%+v`, item)
	}

}
