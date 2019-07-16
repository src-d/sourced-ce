package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InitTestSuite struct {
	IntegrationSuite
}

func TestInitTestSuite(t *testing.T) {
	itt := InitTestSuite{}
	suite.Run(t, &itt)
}

func (s *InitTestSuite) TestWithValidWorkdir() {
	require := s.Require()

	validWorkDir := filepath.Join(s.TestDir, "valid-workdir")
	s.initGitRepo(filepath.Join(validWorkDir, "repo_a"))

	r := s.RunCommand("init", "local", validWorkDir)
	require.NoError(r.Error, r.Combined())

	expectedMsg := [2]string{
		fmt.Sprintf("docker-compose working directory set to %s/workdirs%s", s.sourcedDir, validWorkDir),
		"Initializing source{d}",
	}

	for _, exp := range expectedMsg {
		require.Contains(r.Combined(), exp)
	}
}

func (s *InitTestSuite) TestWithInvalidWorkdir() {
	require := s.Require()

	invalidWorkDir := filepath.Join(s.TestDir, "invalid-workdir")

	_, err := os.Create(invalidWorkDir)
	if err != nil {
		s.T().Fatal(err)
	}

	r := s.RunCommand("init", "local", invalidWorkDir)
	require.Error(r.Error)

	require.Equal(
		fmt.Sprintf("path '%s' is not a valid directory\n", invalidWorkDir),
		r.Stderr(),
	)
}

func (s *InitTestSuite) initGitRepo(path string) {
	s.T().Helper()

	err := os.MkdirAll(path, os.ModePerm)
	s.Require().NoError(err)

	cmd := exec.Command("git", "init", path)
	err = cmd.Run()
	s.Require().NoError(err)
}
