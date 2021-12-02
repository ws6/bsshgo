package download

import (
	"context"
	"os"
	"testing"

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

func TestDownload(t *testing.T) {
	ds := `ds.64c26093cba14416997cf8d48fd8b046`
	client := getNewClient()
	ctx := context.Background()
	ch := client.GetDatasetsFilesChan(ctx, ds)
	ch2 := make(chan *bsshgo.FileS3PresignedUrlResp)
	go func() {
		defer close(ch2)
		for f := range ch {
			t.Logf("%+v\n", f)
			found, err := client.GetFileS3PresignedUrlResp(ctx, f.Id)
			if err != nil {
				t.Fatal(`GetFileS3PresignedUrlResp`, err.Error())
				return
			}
			ch2 <- found
		}
	}()
	opts := &Options{
		DestinationPrefix: `.`,
	}
	if err := Download(ctx, ch2, opts); err != nil {
		t.Fatal(err.Error())
	}
}
