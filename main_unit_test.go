package main

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

const (
	google              = "https://google.com"
	yandex              = "http://yandex.ru"
	yandexWithoutSchema = "yandex.ru"
	badLink             = "link"
	badLinkWithSchema   = "http://link"

	googleResponseEncodedSum = "3ac7705130ed6a4851f1257f68e35862"
	yandexResponseEncodedSum = "c70a08ead6c51126f46bda4fba28453a"
)

var (
	googleData = []byte("google")
	yandexData = []byte("yandex")

	errBadLink = errors.New("error getting data")
)

func Test(t *testing.T) {
	tests := []struct {
		name string

		rawLinks    []string
		concurrency int

		resultMap map[string]string
	}{
		{
			name: "positive",

			concurrency: defaultGoroutinesNumber,

			rawLinks: []string{google, yandex},

			resultMap: map[string]string{
				google: googleResponseEncodedSum,
				yandex: yandexResponseEncodedSum,
			},
		},
		{
			name: "link without schema",

			concurrency: defaultGoroutinesNumber,

			rawLinks: []string{yandexWithoutSchema},

			resultMap: map[string]string{
				yandex: yandexResponseEncodedSum,
			},
		},
		{
			name: "not existed link",

			concurrency: 1,

			rawLinks: []string{badLink},

			resultMap: map[string]string{},
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lh := New(hasherStub{}, getterStub{}, nil)

			result := lh.GetResponsesHashes(test.rawLinks, test.concurrency)

			isResultEqual := reflect.DeepEqual(result, test.resultMap)
			if !isResultEqual {
				t.Fatal(fmt.Sprintf("%s failed", test.name), result, test.resultMap)
			}
		})
	}
}

type hasherStub struct{}

func (h hasherStub) Sum(data []byte) string {
	switch string(data) {
	case string(googleData):
		return googleResponseEncodedSum
	case string(yandexData):
		return yandexResponseEncodedSum
	default:
		return ""
	}
}

type getterStub struct{}

func (g getterStub) Get(link string) ([]byte, error) {
	switch link {
	case google:
		return googleData, nil
	case yandex:
		return yandexData, nil
	case badLinkWithSchema:
		return nil, errBadLink
	default:
		return nil, nil
	}
}
