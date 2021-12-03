package download

import (
	"context"
	"os"
	"path/filepath"
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

func TestDownloadDataset(t *testing.T) {
	ds := `ds.79a1fc0a462844b4baa031fb75a1b817` //https://api.basespace.illumina.com/v2/appsessions/481414979/properties/Output.Datasets
	client := getNewClient()
	// ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*45)
	// defer cancelFn()
	opts := &Options{
		DestinationPrefix: filepath.Join(`.tmp`, ds),
		NumWorkers:        15,
	}
	if err := DownloadDataset(context.Background(), client, ds, opts); err != nil {
		t.Fatal(err.Error())
	}
}

func _TestDownload(t *testing.T) {
	// ds := `ds.64c26093cba14416997cf8d48fd8b046`
	// ds := `ds.6bce605e8a24499087af34ede5230a3f` //tow fastq.gz files
	ds := `ds.79a1fc0a462844b4baa031fb75a1b817` //https://api.basespace.illumina.com/v2/appsessions/481414979/properties/Output.Datasets
	client := getNewClient()
	// ctx, cancelFn := context.WithTimeout(context.Background(), time.Second*45)
	// defer cancelFn()
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	ch := client.GetDatasetsFilesChan(ctx, ds)
	ch2 := make(chan *bsshgo.FileS3PresignedUrlResp)
	go func() {
		defer close(ch2)
		n := 0
		for f := range ch {
			n++
			if n >= 21 && false {
				return
			}
			found, err := client.GetFileS3PresignedUrlResp(ctx, f.Id)
			if err != nil {
				t.Log(`GetFileS3PresignedUrlResp`, err.Error())
				return
			}
			ch2 <- found
		}
	}()
	opts := &Options{
		DestinationPrefix: filepath.Join(`.tmp`, ds),
		NumWorkers:        15,
	}
	if _, err := Download(ctx, ch2, opts); err != nil {
		t.Fatal(`Download err`, err.Error())
	}
	t.Log(`everyting went well`)

}
