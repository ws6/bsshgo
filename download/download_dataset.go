package download

import (
	"context"
	"fmt"

	"github.com/ws6/bsshgo"
)

//download_dataset.go a wrapper to download dataset

func DownloadDataset(ctx0 context.Context, client *bsshgo.Client, ds string, opts *Options) error {
	ctx, cancelFn := context.WithCancel(ctx0)

	defer cancelFn()
	ch := client.GetDatasetsFilesChan(ctx, ds)
	ch2 := make(chan *bsshgo.FileS3PresignedUrlResp, cap(ch))
	errCh := make(chan error)
	done := make(chan struct{})
	var todownload int64
	go func() {
		defer close(ch2)

		for f := range ch {

			found, err := client.GetFileS3PresignedUrlResp(ctx, f.Id)
			if err != nil { //what to do with this err?
				go func() {
					errCh <- fmt.Errorf(`GetFileS3PresignedUrlResp:%s`, err.Error())
				}()
				return
			}

			if opts.IsSkipFn != nil {
				skip, err := opts.IsSkipFn(found)
				if err != nil {
					go func() {
						errCh <- fmt.Errorf(`IsSkipFn:%s`, err.Error())

					}()
					return
				}
				if skip {
					continue
				}

			}
			todownload++

			ch2 <- found
		}
	}()
	var downloaded int64
	go func() {
		_downloaded, err := Download(ctx, ch2, opts)
		downloaded = _downloaded
		if err != nil {
			go func() {
				errCh <- fmt.Errorf(`Download err:%s`, err.Error())
			}()

			return
		}
		done <- struct{}{}
	}()
	select {
	case errRet := <-errCh:
		return errRet
	case <-done:
		if todownload != downloaded {
			return fmt.Errorf(`todownload[%d]!=downloaded[%d]`, todownload, downloaded)
		}

		return nil
	case <-ctx0.Done():
		return ctx.Err()
	}
	panic(`it should never reach here`)
	return nil

}
