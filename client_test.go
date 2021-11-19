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

func TestGetFile(t *testing.T) {
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
