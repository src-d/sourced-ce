package dir

import (
	"os"
	"path"
	"testing"

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
