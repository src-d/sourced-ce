package release

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/src-d/sourced-ce/cmd/sourced/dir"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	responseTag string // tag returned by github
	current     string
	update      bool
	latest      string
}

func TestFindUpdatesSuccess(t *testing.T) {
	os.RemoveAll(filepath.Join(dir.TmpPath(), "httpcache"))

	cases := []testCase{
		{
			responseTag: "v0.14.0",
			current:     "v0.14.0",
			update:      false,
			latest:      "0.14.0",
		},
		{
			responseTag: "v0.11.0",
			current:     "v0.14.0",
			update:      false,
			latest:      "0.11.0",
		},
		{
			responseTag: "v0.14.0",
			current:     "v0.13.0",
			update:      true,
			latest:      "0.14.0",
		},
		{
			responseTag: "v0.14.0",
			current:     "v0.13.1",
			update:      true,
			latest:      "0.14.0",
		},
	}

	for _, c := range cases {
		name := fmt.Sprintf("%s_to_%s", c.current, c.responseTag)
		t.Run(name, func(t *testing.T) {
			restore := mockGithub(c.responseTag)
			defer restore()

			update, latest, err := FindUpdates(c.current)
			assert.Nil(t, err)
			assert.Equal(t, c.update, update)
			assert.Equal(t, c.latest, latest)
		})
	}
}

func mockGithub(tag string) func() {
	originalTransport := http.DefaultTransport

	http.DefaultTransport = &ghTransport{tag: tag}
	return func() {
		http.DefaultTransport = originalTransport
	}
}

type ghTransport struct {
	tag string
}

func (t *ghTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{"tag_name": "%s"}`, t.tag))),
	}, nil
}
