package main

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	t.Run("simplest handler", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_ = afero.WriteFile(fs, CustomBodyExampleFileName, []byte(CustomBodyExampleFileContent), 0664)
		ts := httptest.NewServer(http.HandlerFunc(newHandleFunc(fs, Response{
			Status:   200,
			Method:   "GET",
			Path:     "/test",
			Latency:  "0ms",
			JsonBody: CustomBodyExampleFileName,
		})))
		defer ts.Close()

		res, err := http.Get(ts.URL)
		if err != nil {
			t.FailNow()
		}
		greeting, err := ioutil.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			t.FailNow()
		}

		assert.Equal(t, 200, res.StatusCode)
		assert.Equal(t, CustomBodyExampleFileContent, string(greeting))
	})
}
