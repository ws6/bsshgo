package sage

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/ws6/bsshgo"
)

func getConfig() map[string]string {
	ret := make(map[string]string)
	ret[bsshgo.AUTH_TOKEN] = os.Getenv("BSSH_TEST_AUTH_TOKEN")

	ret[bsshgo.BASE_URL] = os.Getenv("BSSH_TEST_BASE_URL")

	return ret
}

func getNewClient() *bsshgo.Client {
	ret, err := bsshgo.NewClient(getConfig())
	if err != nil {
		panic(err)
	}
	return ret
}

func TestGetFileAtFolder(t *testing.T) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelFn()
	// GetFileRespByName
	client := getNewClient()

	runIds := []string{
		// `219846639`, //no flowcell_barcode and run_id
		// `219700482`, //no flowcell_barcode and run_id
		// `215662457`,
		// `219101888`, //no flowcell_barcode and run_id
		// `219093876`, //no flowcell_barcode and run_id
		// `212632424`, //a incompleted run from icsl prod
		`225792635`,
	}

	for _, rid := range runIds {

		sageMsg, err := GetSageSpec(ctx, client, rid)

		if err != nil {
			t.Fatal(err.Error())
		}
		j, _ := json.MarshalIndent(sageMsg[`flowcell`], "", "  ")
		t.Log(string(j))
	}

	return

	runId := os.Getenv("BSSH_TEST_RUN_ID")

	runInfo, err := GetFlowcellFromFile(ctx, client, runId)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf("%+v", runInfo)

	return

	txt, err := GetRunInfoXml(ctx, client, runId)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Logf(string(txt))

	if _, err := GetRunFileFromRootIgnoreCase(ctx, client, runId, "there_is_no_such_file"); err == nil {
		t.Fatal(`some thing went wrong`)
	}
}
