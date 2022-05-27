package bsshgo

import (
	"context"
	"encoding/json"
	"fmt"
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

//GetAnalysisChan
func TestGetAnalysisChan(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	c := client.GetAnalysisChan(ctx)
	prop := make(map[string]map[string]int)
	for a1 := range c {

		k := a1.Application.Name + "/" + a1.Application.VersionNumber
		if _, ok := prop[k]; !ok {
			prop[k] = make(map[string]int)
		}
		m := prop[k]

		a2, err := client.GetAnalysisById(ctx, a1.Id)
		if err != nil {
			t.Fatal(err.Error())
		}
		for _, pi := range a2.Properties.Items {
			m[pi.Name]++
		}
	}

	for k1, v1 := range prop {
		t.Log(k1)
		for k2, v2 := range v1 {
			t.Log(k2, v2)
		}
	}
}
func _TestGetRunLayout(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	runId := `221597439`
	res, err := client.GetRunLayout(
		ctx,
		runId,
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, found := range res {

		t.Logf(`%+v`, found)
	}

}

// GetRunLayoutFromLanelibrarymappings
func _TestGetGetRunLayoutFromLanelibrarymappings(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	runId := `221597439`
	res, err := client.GetRunLayoutFromLanelibrarymappings(
		ctx,
		runId,
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, found := range res {

		t.Logf(`%+v`, found)
	}

}

// GetLibraryPoolInfo
func _TestGetLibraryPoolInfo(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	poolId := `199512489`
	res, err := client.GetLibraryPoolInfo(
		ctx,
		poolId,
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, found := range res {

		t.Logf(`%+v`, found)
	}

}
func _TestGetGetBioSamplesFromAnalysisThroughFastqDatasetUsed(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	appsessionId := `494909450`
	res, err := client.GetBioSamplesFromAnalysisThroughFastqDatasetUsed(
		ctx,
		appsessionId,
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, found := range res {

		t.Logf(`%+v`, found.BioSample)
	}

}

func _TestGetAnalysisGeneralItems(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	appsessionId := `511567056`
	ch, err := client.GetGeneralItemsChannel(
		ctx,
		fmt.Sprintf(`/v2/appsessions/%s/properties/Input.automation-sample-id.datasets/items`, appsessionId),
		map[string]string{},
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	for found := range ch {
		if found.Err != nil {
			t.Fatal(found.Err.Error())
		}
		t.Logf(`%+v`, found.Item)
	}

}

func _TestUpdateAnalysis(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	appsessionId := `486707226`
	updates := map[string]interface{}{
		`DeliveryStatus`: `Delivered`,
	}
	res, err := client.UpdateAnalysis(ctx, appsessionId, updates)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(`%+v`, res)
}

// GetAnalysisOutputDatasetChan

func _TestGetAnalysisOutputDatasetChan(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	appsessionId := `486707227`
	ch, err := client.GetAnalysisOutputDatasetChan(ctx, appsessionId)
	if err != nil {
		t.Fatal(err.Error())
	}
	found := 0
	for item := range ch {
		t.Logf("%+v\n", item)
		found++
	}
	t.Logf("found total dataset %d\n", found)
}

func _TestFindOneAnalysisByName(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout
	name := `ILS_DRAGEN_GL_2.0.2 11/08/2021 12:41:34`
	found, err := client.FindOneAnalysisByName(ctx, name)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("%+v\n", found)
}

//GetDatasetsFilesChan
func _TestGetDatasetsFilesChan(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
	// GetRunSampleSheetLayout

	ch := client.GetDatasetsFilesChan(ctx, `ds.64c26093cba14416997cf8d48fd8b046`)
	found := 0
	for item := range ch {
		t.Log(item)
		found++
	}
	t.Logf("found total files %d\n", found)
}

//GetAnalysisById
func _TestGetAnalysisById(t *testing.T) {
	client := getNewClient()
	ctx := context.Background()
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
	// 486869383
	_biosamples, err := client.GetAnalysisBioSamples(ctx, `486869383`)
	if err != nil {
		t.Fatal(err.Error())
	}
	for bs := range _biosamples {
		t.Logf("%+v", bs)

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

var samplesheetresp = `{
	"NeedsAttention": false,
	"FormatVersion": "2",
	"HeaderSettings": [{
		"Key": "Date",
		"Value": "12/11/2021"
	}, {
		"Key": "FileFormatVersion",
		"Value": "2"
	}, {
		"Key": "InstrumentPlatform",
		"Value": "NovaSeq6000"
	}, {
		"Key": "InstrumentType",
		"Value": "NovaSeq6000"
	}, {
		"Key": "Investigator Name",
		"Value": "Bradley Durham"
	}, {
		"Key": "RunName",
		"Value": "TechnicalReport_TEST"
	}],
	"ReadSettings": [{
		"Key": "Index1Cycles",
		"Value": "8"
	}, {
		"Key": "Index2Cycles",
		"Value": "8"
	}, {
		"Key": "Read1Cycles",
		"Value": "151"
	}, {
		"Key": "Read2Cycles",
		"Value": "151"
	}],
	"SequencingSettings": [],
	"Applications": [{
		"Name": "BCLConvert",
		"Settings": [{
			"Key": "AdapterRead1",
			"Value": "CTGTCTCTTATACACATCT"
		}, {
			"Key": "AdapterRead2",
			"Value": "CTGTCTCTTATACACATCT"
		}, {
			"Key": "SoftwareVersion",
			"Value": "3.5.8"
		}],
		"Data": [{
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A01",
			"index": "GACCTGAA",
			"index2": "TTGGTGAG"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A02",
			"index": "ATATGGAT",
			"index2": "CTGTATTA"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A03",
			"index": "GCGCAAGC",
			"index2": "TCACGCCG"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A04",
			"index": "AAGATACT",
			"index2": "ACTTACAT"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A05",
			"index": "GGAGCGTC",
			"index2": "GTCCGTGC"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A06",
			"index": "ATGGCATG",
			"index2": "AAGGTACC"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A07",
			"index": "GCAATGCA",
			"index2": "GGAACGTT"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A08",
			"index": "ACCTTGGC",
			"index2": "GGCCTCAT"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A09",
			"index": "ATATCTCG",
			"index2": "ATCTTAGT"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A10",
			"index": "GCGCTCTA",
			"index2": "GCTCCGAC"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A11",
			"index": "GGTGAACC",
			"index2": "GCGTTGGA"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_A12",
			"index": "CAACAATG",
			"index2": "CTTCACGG"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B01",
			"index": "TGGTGGCA",
			"index2": "TCCTGTAA"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B02",
			"index": "AGGCAGAG",
			"index2": "AGAATGCC"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B03",
			"index": "GAATGAGA",
			"index2": "GAGGCATT"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B04",
			"index": "TGCGGCGT",
			"index2": "CCTCGGTA"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B05",
			"index": "CATAATAC",
			"index2": "TTCTAACG"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B06",
			"index": "GATCTATC",
			"index2": "ATGAGGCT"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B07",
			"index": "AGCTCGCT",
			"index2": "GCAGAATC"
		}, {
			"Lane": "1",
			"Sample_ID": "LP9000124-DNA_B08",
			"index": "CGGAACTG",
			"index2": "CACTACGA"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A01",
			"index": "GACCTGAA",
			"index2": "TTGGTGAG"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A02",
			"index": "ATATGGAT",
			"index2": "CTGTATTA"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A03",
			"index": "GCGCAAGC",
			"index2": "TCACGCCG"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A04",
			"index": "AAGATACT",
			"index2": "ACTTACAT"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A05",
			"index": "GGAGCGTC",
			"index2": "GTCCGTGC"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A06",
			"index": "ATGGCATG",
			"index2": "AAGGTACC"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A07",
			"index": "GCAATGCA",
			"index2": "GGAACGTT"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A08",
			"index": "ACCTTGGC",
			"index2": "GGCCTCAT"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A09",
			"index": "ATATCTCG",
			"index2": "ATCTTAGT"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A10",
			"index": "GCGCTCTA",
			"index2": "GCTCCGAC"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A11",
			"index": "GGTGAACC",
			"index2": "GCGTTGGA"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_A12",
			"index": "CAACAATG",
			"index2": "CTTCACGG"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B01",
			"index": "TGGTGGCA",
			"index2": "TCCTGTAA"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B02",
			"index": "AGGCAGAG",
			"index2": "AGAATGCC"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B03",
			"index": "GAATGAGA",
			"index2": "GAGGCATT"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B04",
			"index": "TGCGGCGT",
			"index2": "CCTCGGTA"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B05",
			"index": "CATAATAC",
			"index2": "TTCTAACG"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B06",
			"index": "GATCTATC",
			"index2": "ATGAGGCT"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B07",
			"index": "AGCTCGCT",
			"index2": "GCAGAATC"
		}, {
			"Lane": "2",
			"Sample_ID": "LP9000124-DNA_B08",
			"index": "CGGAACTG",
			"index2": "CACTACGA"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A01",
			"index": "GACCTGAA",
			"index2": "TTGGTGAG"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A02",
			"index": "ATATGGAT",
			"index2": "CTGTATTA"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A03",
			"index": "GCGCAAGC",
			"index2": "TCACGCCG"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A04",
			"index": "AAGATACT",
			"index2": "ACTTACAT"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A05",
			"index": "GGAGCGTC",
			"index2": "GTCCGTGC"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A06",
			"index": "ATGGCATG",
			"index2": "AAGGTACC"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A07",
			"index": "GCAATGCA",
			"index2": "GGAACGTT"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A08",
			"index": "ACCTTGGC",
			"index2": "GGCCTCAT"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A09",
			"index": "ATATCTCG",
			"index2": "ATCTTAGT"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A10",
			"index": "GCGCTCTA",
			"index2": "GCTCCGAC"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A11",
			"index": "GGTGAACC",
			"index2": "GCGTTGGA"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_A12",
			"index": "CAACAATG",
			"index2": "CTTCACGG"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B01",
			"index": "TGGTGGCA",
			"index2": "TCCTGTAA"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B02",
			"index": "AGGCAGAG",
			"index2": "AGAATGCC"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B03",
			"index": "GAATGAGA",
			"index2": "GAGGCATT"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B04",
			"index": "TGCGGCGT",
			"index2": "CCTCGGTA"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B05",
			"index": "CATAATAC",
			"index2": "TTCTAACG"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B06",
			"index": "GATCTATC",
			"index2": "ATGAGGCT"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B07",
			"index": "AGCTCGCT",
			"index2": "GCAGAATC"
		}, {
			"Lane": "3",
			"Sample_ID": "LP9000124-DNA_B08",
			"index": "CGGAACTG",
			"index2": "CACTACGA"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A01",
			"index": "GACCTGAA",
			"index2": "TTGGTGAG"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A02",
			"index": "ATATGGAT",
			"index2": "CTGTATTA"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A03",
			"index": "GCGCAAGC",
			"index2": "TCACGCCG"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A04",
			"index": "AAGATACT",
			"index2": "ACTTACAT"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A05",
			"index": "GGAGCGTC",
			"index2": "GTCCGTGC"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A06",
			"index": "ATGGCATG",
			"index2": "AAGGTACC"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A07",
			"index": "GCAATGCA",
			"index2": "GGAACGTT"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A08",
			"index": "ACCTTGGC",
			"index2": "GGCCTCAT"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A09",
			"index": "ATATCTCG",
			"index2": "ATCTTAGT"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A10",
			"index": "GCGCTCTA",
			"index2": "GCTCCGAC"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A11",
			"index": "GGTGAACC",
			"index2": "GCGTTGGA"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_A12",
			"index": "CAACAATG",
			"index2": "CTTCACGG"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B01",
			"index": "TGGTGGCA",
			"index2": "TCCTGTAA"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B02",
			"index": "AGGCAGAG",
			"index2": "AGAATGCC"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B03",
			"index": "GAATGAGA",
			"index2": "GAGGCATT"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B04",
			"index": "TGCGGCGT",
			"index2": "CCTCGGTA"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B05",
			"index": "CATAATAC",
			"index2": "TTCTAACG"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B06",
			"index": "GATCTATC",
			"index2": "ATGAGGCT"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B07",
			"index": "AGCTCGCT",
			"index2": "GCAGAATC"
		}, {
			"Lane": "4",
			"Sample_ID": "LP9000124-DNA_B08",
			"index": "CGGAACTG",
			"index2": "CACTACGA"
		}]
	}, {
		"Name": "Cloud",
		"Settings": [{
			"Key": "BsshApp",
			"Value": "illumina-inc.bcl-convert.1.3.0"
		}],
		"Data": [{
			"Sample_ID": "LP9000124-DNA_A01",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_A08_805b2603-3381-4a14-9626-486f4aa26ad4"
		}, {
			"Sample_ID": "LP9000124-DNA_A02",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_A09_1ad9e529-609d-45fb-b4ef-5ad969a976dd"
		}, {
			"Sample_ID": "LP9000124-DNA_A03",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_B09_a11d25c4-cc0a-486c-af73-73238b2ab594"
		}, {
			"Sample_ID": "LP9000124-DNA_A04",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_B09_c3ece952-099a-4a2e-9a1e-8bf662608797"
		}, {
			"Sample_ID": "LP9000124-DNA_A05",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_C07_40d066a2-fa90-494e-a9b5-5d56a07be16a"
		}, {
			"Sample_ID": "LP9000124-DNA_A06",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_C08_3f528def-426c-4302-b1f7-1056c0f31056"
		}, {
			"Sample_ID": "LP9000124-DNA_A07",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_C09_05d17dd0-b700-4f7e-88bf-f473ec6814eb"
		}, {
			"Sample_ID": "LP9000124-DNA_A08",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_D07_48f06831-c766-4d6f-ab4d-9b130cb58421"
		}, {
			"Sample_ID": "LP9000124-DNA_A09",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_D08_e200i009-d222-4d6f-ab4d-9b130cb58422"
		}, {
			"Sample_ID": "LP9000124-DNA_A10",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_D09_49f06831-c766-4d6f-ab4d-9b130cb58423"
		}, {
			"Sample_ID": "LP9000124-DNA_A11",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_E07_41f06831-c766-4d6f-ab4d-9b130cb58424"
		}, {
			"Sample_ID": "LP9000124-DNA_A12",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_E08_42f06831-c766-4d6f-ab4d-9b130cb58425"
		}, {
			"Sample_ID": "LP9000124-DNA_B01",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_E09_43f06831-c766-4d6f-ab4d-9b130cb58426"
		}, {
			"Sample_ID": "LP9000124-DNA_B02",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_F07_44f06831-c766-4d6f-ab4d-9b130cb58427"
		}, {
			"Sample_ID": "LP9000124-DNA_B03",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_F08_45f06831-c766-4d6f-ab4d-9b130cb58428"
		}, {
			"Sample_ID": "LP9000124-DNA_B04",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_F09_51f06831-c766-4d6f-ab4d-9b130cb58429"
		}, {
			"Sample_ID": "LP9000124-DNA_B05",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_G07_52f06831-c766-4d6f-ab4d-9b130cb58430"
		}, {
			"Sample_ID": "LP9000124-DNA_B06",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_G08_53f06831-c766-4d6f-ab4d-9b130cb58431"
		}, {
			"Sample_ID": "LP9000124-DNA_B07",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_H07_54f06831-c766-4d6f-ab4d-9b130cb58432"
		}, {
			"Sample_ID": "LP9000124-DNA_B08",
			"ProjectName": "jingtao-test",
			"LibraryName": "LP9000124-NTP_H08_55f06831-c766-4d6f-ab4d-9b130cb58433"
		}]
	}]
}`

func _TestConcateSampleSheet(t *testing.T) {
	tocat := new(RunSampleSheetLayoutResp)
	if err := json.Unmarshal([]byte(samplesheetresp), tocat); err != nil {
		t.Fatal(err.Error())
	}
	combo := ConcateRunLayoutWithBCLConvertAndCloudApplicaions(tocat)

	for _, layout := range combo {
		if layout.Sample_ID != `LP9000124-DNA_A01` {
			continue
		}
		t.Logf("%+v\n", layout)
	}
}
