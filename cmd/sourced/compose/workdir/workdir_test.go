package workdir

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	assert := assert.New(t)

	workdirsPath := path.Join(os.TempDir(), "builder")
	defer func() {
		os.RemoveAll(workdirsPath)
	}()

	b := builder{workdirsPath: workdirsPath}

	// incorrect: not in workdirsPath
	_, err := b.Build("/not/in/workdirs")
	assert.EqualError(err, "invalid workdir type for path /not/in/workdirs")

	// incorrect: unknown type
	unknownDir := path.Join(workdirsPath, "unknown")
	_, err = b.Build(unknownDir)
	assert.EqualError(err, "invalid workdir type for path "+unknownDir)

	// local
	name := "some"
	localDir := path.Join(workdirsPath, "local", encodeDirName(name))
	wd, err := b.Build(localDir)
	assert.Nil(err)
	assert.Equal(Local, wd.Type)
	assert.Equal(name, wd.Name)
	assert.Equal(localDir, wd.Path)

	// org
	orgDir := path.Join(workdirsPath, "orgs", encodeDirName(name))
	wd, err = b.Build(orgDir)
	assert.Nil(err)
	assert.Equal(Orgs, wd.Type)
	assert.Equal(name, wd.Name)
	assert.Equal(orgDir, wd.Path)
}

func TestIsEmptyFile(t *testing.T) {
	assert := assert.New(t)

	// not exist
	ok, err := isEmptyFile("/does/not/exist")
	assert.Nil(err)
	assert.True(ok)

	// empty
	emptyPath := path.Join(os.TempDir(), "empty")
	defer func() {
		os.RemoveAll(emptyPath)
	}()
	f, err := os.Create(emptyPath)
	assert.Nil(err)
	assert.Nil(f.Close())

	ok, err = isEmptyFile(emptyPath)
	assert.Nil(err)
	assert.True(ok)

	// not empty
	nonEmptyPath := path.Join(os.TempDir(), "non-empty")
	defer func() {
		os.RemoveAll(nonEmptyPath)
	}()
	f, err = os.Create(nonEmptyPath)
	assert.Nil(err)
	_, err = f.Write([]byte("some content"))
	assert.Nil(err)
	assert.Nil(f.Close())

	ok, err = isEmptyFile(nonEmptyPath)
	assert.Nil(err)
	assert.False(ok)
}
