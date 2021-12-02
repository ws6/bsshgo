package download

import (
	"context"
	"fmt"
	"os"

	"io"
	"path/filepath"

	"net/url"
	"strings"

	"github.com/ws6/bsshgo"
	"github.com/ws6/bsshgo/s3io"
)

//download_file.go download bssh files
type Options struct {
	DestinationPrefix string
	NumWorkers        int
}

type s3Loc struct {
	bucket string
	key    string
}

func getBucketAndKeyFromPresignedS3URL(_url string) (*s3Loc, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}

	ret := new(s3Loc)
	path := strings.SplitN(u.Path, "/", 3)
	if len(path) < 3 {
		return nil, fmt.Errorf(`malformatted URL:%s`, _url)
	}
	ret.bucket = path[1]
	ret.key = path[2]
	return ret, nil
}

func getPresignedS3UrlReader(_url string) (io.ReadCloser, error) {
	u, err := url.Parse(_url)
	if err != nil {

		return nil, err
	}

	r, _, err := s3io.PreSigned(*u, nil)

	return r, err
}

func Download(ctx context.Context, ch chan *bsshgo.FileS3PresignedUrlResp, opts *Options) error {

	for file := range ch {

		loc, err := getBucketAndKeyFromPresignedS3URL(file.HrefContent)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		outputFile := filepath.Join(opts.DestinationPrefix, loc.key)
		fmt.Println(`downloading`, file, `to`, outputFile)
		rc, err := getPresignedS3UrlReader(file.HrefContent)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		ofh, err := os.Create(outputFile)
		n, err := io.Copy(ofh, rc)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		ofh.Close()
		rc.Close()
		fmt.Println(outputFile, n, `bytes copied`)

	}
	return nil
}
