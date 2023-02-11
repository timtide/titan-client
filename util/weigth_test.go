package util

import (
	"github.com/linguohua/titan/api"
	"testing"
)

var testData = []*api.DownloadInfoResult{
	{
		URL:    "first",
		Weight: 1,
	},
	{
		URL:    "two",
		Weight: 2,
	},
	{
		URL:    "three",
		Weight: 3,
	},
	{
		URL:    "four",
		Weight: 4,
	},
}

func TestRandomDraw(t *testing.T) {
	ch, err := NewChooser(testData...)
	if err != nil {
		t.Error(err)
		return
	}
	res := make(map[*api.DownloadInfoResult]int)
	for _, v := range testData {
		res[v] = 0
	}
	for i := 0; i < 10000; i++ {
		item := ch.Pick()
		res[item] += 1
	}

	for k, v := range res {
		t.Logf("%s %s%d,%s%d", k.URL, "weigth:", k.Weight, "total times:", v)
	}
}
