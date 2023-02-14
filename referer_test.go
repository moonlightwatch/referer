package referer_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moonlightwatch/referer"

	"gotest.tools/assert"
)

type TestCase struct {
	Referer    string
	Config     *referer.Config
	StatusCode int
}

func TestReferer(t *testing.T) {

	testCases := []TestCase{
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: false, Type: "white"},
			StatusCode: http.StatusOK,
		},
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"*.baidu.com"}, EmptyReferer: false, Type: "white"},
			StatusCode: http.StatusForbidden,
		},
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: true, Type: "white"},
			StatusCode: http.StatusOK,
		},
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"*.baidu.com"}, EmptyReferer: true, Type: "white"},
			StatusCode: http.StatusForbidden,
		},
		{
			Referer:    "",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: true, Type: "white"},
			StatusCode: http.StatusOK,
		},
		{
			Referer:    "",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: false, Type: "white"},
			StatusCode: http.StatusForbidden,
		},

		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: false, Type: "black"},
			StatusCode: http.StatusForbidden,
		},
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"*.baidu.com"}, EmptyReferer: false, Type: "black"},
			StatusCode: http.StatusOK,
		},
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: true, Type: "black"},
			StatusCode: http.StatusForbidden,
		},
		{
			Referer:    "http://baidu.com",
			Config:     &referer.Config{Domains: []string{"*.baidu.com"}, EmptyReferer: true, Type: "black"},
			StatusCode: http.StatusOK,
		},
		{
			Referer:    "",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: true, Type: "black"},
			StatusCode: http.StatusForbidden,
		},
		{
			Referer:    "",
			Config:     &referer.Config{Domains: []string{"baidu.com"}, EmptyReferer: false, Type: "black"},
			StatusCode: http.StatusOK,
		},
	}
	ctx := context.Background()
	for index, testcase := range testCases {

		next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

		handler, err := referer.New(ctx, next, testcase.Config, fmt.Sprintf("testcase_%d", index))

		if err != nil {
			t.Fatalf("New Referer error: %+v\n", err)
		}
		recorder := httptest.NewRecorder()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)

		if err != nil {
			t.Fatal(err)
		}
		req.Header.Add("Referer", testcase.Referer)
		handler.ServeHTTP(recorder, req)
		fmt.Printf("testcase_%d\n", index)
		assert.Equal(t, recorder.Result().StatusCode, testcase.StatusCode)

	}

}
