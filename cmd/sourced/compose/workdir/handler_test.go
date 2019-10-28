package workdir

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type HandlerSuite struct {
	suite.Suite

	h             *Handler
	originSrcdDir string
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, &HandlerSuite{})
}

func (s *HandlerSuite) BeforeTest(suiteName, testName string) {
	s.originSrcdDir = os.Getenv("SOURCED_DIR")

	// on macOs os.TempDir returns symlink and tests fails
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	srdPath := path.Join(tmpDir, testName)
	err := os.MkdirAll(srdPath, os.ModePerm)
	s.Nil(err)

	os.Setenv("SOURCED_DIR", srdPath)

	s.h, err = NewHandler()
	s.Nil(err)
}

func (s *HandlerSuite) AfterTest(suiteName, testName string) {
	os.RemoveAll(filepath.Dir(s.h.workdirsPath))
	os.Setenv("SOURCED_DIR", s.originSrcdDir)
}

// This tests only public interface without checking implementation (filesystem) details
func (s *HandlerSuite) TestSuccessFlow() {
	wd := s.createWd("flow")

	s.Nil(s.h.Validate(wd))
	s.Nil(s.h.SetActive(wd))

	active, err := s.h.Active()
	s.Nil(err)
	s.Equal(wd, active)

	s.Nil(s.h.UnsetActive())

	_, err = s.h.Active()
	s.True(ErrMalformed.Is(err))

	wds, err := s.h.List()
	s.Nil(err)
	s.Len(wds, 1)
	s.Equal(wd, wds[0])

	s.Nil(s.h.Remove(wd))

	wds, err = s.h.List()
	s.Nil(err)
	s.Len(wds, 0)
}

// All tests below rely on implementation details to check error cases

func (s *HandlerSuite) TestSetActiveOk() {
	wd := s.createWd("some")

	// non-active before
	s.Nil(s.h.SetActive(wd))
	// re-activation should also work
	s.Nil(s.h.SetActive(wd))

	// validate link points correctly
	target, err := filepath.EvalSymlinks(path.Join(s.h.workdirsPath, activeDir))
	s.Nil(err)
	s.Equal(wd.Path, target)
}

func (s *HandlerSuite) TestSetActiveError() {
	wd := s.createWd("some")

	// break active path by making it dir with files
	activePath := path.Join(s.h.workdirsPath, activeDir)
	s.Nil(os.MkdirAll(activePath, os.ModePerm))
	_, err := os.Create(path.Join(activePath, "some-file"))
	s.Nil(err)

	s.Error(s.h.SetActive(wd))
}

func (s *HandlerSuite) TestUnsetActiveOk() {
	activePath := path.Join(s.h.workdirsPath, activeDir)
	s.Nil(os.MkdirAll(s.h.workdirsPath, os.ModePerm))
	_, err := os.Create(activePath)
	s.Nil(err)

	s.Nil(s.h.UnsetActive())
	// unset without active dir
	s.Nil(s.h.UnsetActive())

	// validate we deleted the file
	_, err = os.Stat(activePath)
	s.True(os.IsNotExist(err))
}

func (s *HandlerSuite) TestUnsetActiveError() {
	// break active path by making it dir with files
	activePath := path.Join(s.h.workdirsPath, activeDir)
	s.Nil(os.MkdirAll(activePath, os.ModePerm))
	_, err := os.Create(path.Join(activePath, "some-file"))
	s.Nil(err)

	s.Error(s.h.UnsetActive())
}

func (s *HandlerSuite) TestValidateError() {
	// dir doesn't exist
	wd, err := s.h.builder.Build(path.Join(s.h.workdirsPath, "local", "some"))
	s.Nil(err)
	err = s.h.Validate(wd)
	s.True(ErrMalformed.Is(err))

	// dir is a file
	s.Nil(os.MkdirAll(path.Join(s.h.workdirsPath, "local"), os.ModePerm))
	_, err = os.Create(wd.Path)
	s.Nil(err)

	err = s.h.Validate(wd)
	s.True(ErrMalformed.Is(err))
	s.Nil(os.RemoveAll(wd.Path))

	// dir is empty
	s.Nil(os.MkdirAll(wd.Path, os.ModePerm))
	err = s.h.Validate(wd)
	s.True(ErrMalformed.Is(err))
}

func (s *HandlerSuite) TestListOk() {
	s.Nil(os.MkdirAll(s.h.workdirsPath, os.ModePerm))

	// empty results
	wds, err := s.h.List()
	s.Nil(err)
	s.Len(wds, 0)

	// multiple results
	s.createWd("one")
	s.createWd("two")

	wds, err = s.h.List()
	s.Nil(err)
	s.Len(wds, 2)

	s.Equal("one", wds[0].Name)
	s.Equal("two", wds[1].Name)

	// incorrect directory should be skipped
	wd, err := s.h.builder.Build(path.Join(s.h.workdirsPath, "local", "some"))
	s.Nil(err)
	s.Nil(os.MkdirAll(wd.Path, os.ModePerm))

	wds, err = s.h.List()
	s.Nil(err)
	s.Len(wds, 2)
}

func (s *HandlerSuite) TestListError() {
	// workdirs dir doesn't exist
	_, err := s.h.List()
	s.True(ErrMalformed.Is(err))
}

func (s *HandlerSuite) TestRemoveOk() {
	// local
	wd, err := InitLocal("local")
	s.Nil(err)
	s.Nil(s.h.Remove(wd))
	_, err = os.Stat(wd.Path)
	s.True(os.IsNotExist(err))

	// org
	wd, err = InitOrgs([]string{"some-org"}, "token", false)
	s.Nil(err)
	s.Nil(s.h.Remove(wd))
	_, err = os.Stat(wd.Path)
	s.True(os.IsNotExist(err))

	// skip deleting dir with extra files
	wd = s.createWd("some")
	_, err = os.Create(path.Join(wd.Path, "some-file"))
	s.Nil(err)

	s.Nil(s.h.Remove(wd))
	_, err = os.Stat(wd.Path)
	s.Nil(err)
}

func (s *HandlerSuite) createWd(name string) *Workdir {
	wd, err := InitLocal(name)
	s.Nil(err)
	return wd
}
