package dir

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert := assert.New(t)

	err := validate("/does/not/exists")
	assert.True(ErrNotExist.Is(err))
	assert.EqualError(err, "/does/not/exists does not exist")

	// with a file
	tmpFile := path.Join(os.TempDir(), "tmp-file")
	f, err := os.Create(tmpFile)
	assert.Nil(err)
	assert.Nil(f.Close())
	defer func() {
		os.RemoveAll(tmpFile)
	}()

	err = validate(tmpFile)
	assert.True(ErrNotValid.Is(err))
	assert.EqualError(err, tmpFile+" is not a valid config directory: it is not a directory")

	// with a dir
	tmpDir := path.Join(os.TempDir(), "tmp-dir")
	assert.Nil(os.Mkdir(tmpDir, os.ModePerm))
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	err = validate(tmpDir)
	assert.Nil(err)

	// read only
	assert.Nil(os.Chmod(tmpDir, 0444))
	err = validate(tmpDir)
	assert.True(ErrNotValid.Is(err))
	assert.EqualError(err, tmpDir+" is not a valid config directory: it has no read-write access")

	// write only
	assert.Nil(os.Chmod(tmpDir, 0222))
	err = validate(tmpDir)
	assert.True(ErrNotValid.Is(err))
	assert.EqualError(err, tmpDir+" is not a valid config directory: it has no read-write access")
}

func TestPrepare(t *testing.T) {
	assert := assert.New(t)

	tmpDir := path.Join(os.TempDir(), "tmp-dir")
	assert.Nil(os.Mkdir(tmpDir, os.ModePerm))
	defer func() {
		os.RemoveAll(tmpDir)
	}()

	originSrcdDir := os.Getenv("SOURCED_DIR")
	defer func() {
		os.Setenv("SOURCED_DIR", originSrcdDir)
	}()

	os.Setenv("SOURCED_DIR", tmpDir)
	assert.Nil(Prepare())

	toCreateDir := path.Join(os.TempDir(), "to-create-dir")
	defer func() {
		os.RemoveAll(toCreateDir)
	}()
	_, err := os.Stat(toCreateDir)
	assert.True(os.IsNotExist(err))

	os.Setenv("SOURCED_DIR", toCreateDir)
	assert.Nil(Prepare())
	_, err = os.Stat(toCreateDir)
	assert.Nil(err)
}

func TestDownloadURL(t *testing.T) {
	assert := assert.New(t)

	// success
	fileContext := []byte("hello")
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fileContext)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	dirPath := path.Join(os.TempDir(), "some-dir")
	filePath := path.Join(dirPath, "file-to-download")
	defer func() {
		os.RemoveAll(dirPath)
	}()

	assert.Nil(DownloadURL(server.URL, filePath))
	_, err := os.Stat(filePath)
	assert.Nil(err)

	b, err := ioutil.ReadFile(filePath)
	assert.Nil(err)
	assert.Equal(fileContext, b)

	// error
	handler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}
	server = httptest.NewServer(http.HandlerFunc(handler))
	err = DownloadURL(server.URL, "/dev/null")
	errExpected := errors.Wrapf(
		fmt.Errorf("HTTP status %v", "404 Not Found"),
		"network error downloading %s", server.URL,
	)

	assert.EqualError(err, errExpected.Error())
}
