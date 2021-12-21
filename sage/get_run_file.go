package sage

//get_run_file.go an enhancement of trying to reach out actual file contents

// https://api.basespace.illumina.com/v2/runs/219846639/files?recursive=false&filehrefcontentresolution=false&sortdir=Asc&sortby=Name&limit=2000&directory=/

//ref https://developer.basespace.illumina.com/docs/content/documentation/rest-api/api-reference#operation--runs--id--files-get

import (
	"context"

	"strings"

	"github.com/ws6/bsshgo"
	"github.com/ws6/interop/fcinfo"
)

func GetFlowcellFromFile(ctx context.Context, client *bsshgo.Client, runId string) (*fcinfo.Flowcell, error) {
	runInfoByte, err := GetRunInfoXml(ctx, client, runId)
	if err != nil {
		return nil, err
	}
	runParamsByte, err := GetRunParametersXml(ctx, client, runId)
	if err != nil {
		return nil, err
	}

	return fcinfo.ParseFlowcell(string(runInfoByte), string(runParamsByte))
}

func GetRunInfoXml(ctx context.Context, client *bsshgo.Client, runId string,
) ([]byte, error) {
	return GetRunFileFromRootIgnoreCase(ctx, client, runId, `RunInfo.xml`)
}

func GetRunParametersXml(ctx context.Context, client *bsshgo.Client, runId string,
) ([]byte, error) {
	return GetRunFileFromRootIgnoreCase(ctx, client, runId, `RunParameters.xml`)
}

// SampleSheet.csv
func GetRunSampleSheetCsv(ctx context.Context, client *bsshgo.Client, runId string,
) ([]byte, error) {
	return GetRunFileFromRootIgnoreCase(ctx, client, runId, `SampleSheet.csv`)
}

func GetRunFileFromRootIgnoreCase(ctx context.Context, client *bsshgo.Client,

	runId, fn string,
) ([]byte, error) {
	params := map[string]string{
		`extensions`: ".xml",
		`directory`:  "/",
		`recursive:`: "false",
	}
	searchRunInfoFn := func(f *bsshgo.FileItem) bool {
		return strings.ToLower(f.Name) == strings.ToLower(fn)
	}
	found, err := client.GetFileFromDir(ctx, runId, params, "/", searchRunInfoFn)
	if err != nil {
		return nil, err
	}
	return client.GetFileBytes(ctx, found.Id)
}
