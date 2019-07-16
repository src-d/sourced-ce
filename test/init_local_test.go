package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type InitLocalTestSuite struct {
	IntegrationSuite
}

func TestInitLocalTestSuite(t *testing.T) {
	itt := InitLocalTestSuite{}
	suite.Run(t, &itt)
}

func (s *InitLocalTestSuite) TestWithInvalidWorkdir() {
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

func (s *InitLocalTestSuite) TestChangeWorkdir() {
	require := s.Require()

	// TODO will need to change with https://github.com/src-d/sourced-ce/issues/144
	r := s.RunCommand("workdirs")
	require.Error(r.Error)

	// Create 2 workdirs, each with a repo
	workdirA := filepath.Join(s.TestDir, "workdir_a")
	workdirB := filepath.Join(s.TestDir, "workdir_b")
	pathA := filepath.Join(workdirA, "repo_a")
	pathB := filepath.Join(workdirB, "repo_b")

	s.initGitRepo(pathA)
	s.initGitRepo(pathB)

	// init with workdir A
	r = s.RunCommand("init", "local", workdirA)
	require.NoError(r.Error, r.Combined())

	r = s.RunCommand("workdirs")
	require.NoError(r.Error, r.Combined())

	require.Equal(fmt.Sprintf("* %v\n", workdirA), r.Stdout())

	r = s.RunCommand("sql", "select * from repositories")
	require.NoError(r.Error, r.Combined())

	require.Contains(r.Stdout(),
		`repository_id
repo_a
`)

	// init with workdir B
	r = s.RunCommand("init", "local", workdirB)
	require.NoError(r.Error, r.Combined())

	r = s.RunCommand("workdirs")
	require.NoError(r.Error, r.Combined())

	require.Equal(fmt.Sprintf("  %v\n* %v\n", workdirA, workdirB), r.Stdout())

	r = s.RunCommand("sql", "select * from repositories")
	require.NoError(r.Error, r.Combined())

	require.Contains(r.Stdout(),
		`repository_id
repo_b
`)

	// Test SQL queries. This should be a different test, but since starting
	// the environment takes a long time, it is bundled together here to speed up
	// the tests
	s.testSQL()
}

func (s *InitLocalTestSuite) initGitRepo(path string) {
	s.T().Helper()

	err := os.MkdirAll(path, os.ModePerm)
	s.Require().NoError(err)

	cmd := exec.Command("git", "init", path)
	err = cmd.Run()
	s.Require().NoError(err)
}
