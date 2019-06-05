// Package release deals with versioning and releases
package release

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/src-d/sourced-ce/cmd/sourced/dir"

	"github.com/blang/semver"
	"github.com/google/go-github/v25/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

// FindUpdates calls the GitHub API to check the latest release tag. It returns
// true if the latest stable release is newer than the current tag, and also
// that latest tag name.
func FindUpdates(current string) (update bool, latest string, err error) {
	currentV, err := semver.ParseTolerant(current)
	if err != nil {
		return false, "", err
	}

	diskcachePath := filepath.Join(dir.TmpPath(), "httpcache")
	err = os.MkdirAll(diskcachePath, os.ModePerm)
	if err != nil {
		return false, "", err
	}

	cache := diskcache.New(diskcachePath)
	client := github.NewClient(&http.Client{Transport: httpcache.NewTransport(cache)})

	rel, _, err := client.Repositories.GetLatestRelease(context.Background(), "src-d", "sourced-ce")
	if err != nil {
		return false, "", err
	}

	latestV, err := semver.ParseTolerant(rel.GetTagName())
	if err != nil {
		return false, "", err
	}

	update = latestV.GT(currentV)
	latest = latestV.String()

	return update, latest, nil
}
