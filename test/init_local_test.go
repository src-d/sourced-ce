// +build integration

package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
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
	req := s.Require()

	r := s.RunCommand("status", "workdirs")
	req.Error(r.Error)

	// Create 2 workdirs, each with a repo
	workdirA := filepath.Join(s.TestDir, "workdir_a")
	workdirB := filepath.Join(s.TestDir, "workdir_b")
	pathA := filepath.Join(workdirA, "repo_a")
	pathB := filepath.Join(workdirB, "repo_b")

	s.initGitRepo(pathA)
	s.initGitRepo(pathB)

	// init with workdir A
	r = s.RunCommand("init", "local", workdirA)
	req.NoError(r.Error, r.Combined())

	r = s.RunCommand("status", "workdirs")
	req.NoError(r.Error, r.Combined())

	req.Equal(fmt.Sprintf("* %v\n", workdirA), r.Stdout())

	r = s.RunCommand("sql", "select * from repositories")
	req.NoError(r.Error, r.Combined())

	req.Contains(r.Stdout(),
		`repository_id
repo_a
`)

	// init with workdir B
	r = s.RunCommand("init", "local", workdirB)
	req.NoError(r.Error, r.Combined())

	r = s.RunCommand("status", "workdirs")
	req.NoError(r.Error, r.Combined())

	req.Equal(fmt.Sprintf("  %v\n* %v\n", workdirA, workdirB), r.Stdout())

	r = s.RunCommand("sql", "select * from repositories")
	req.NoError(r.Error, r.Combined())

	req.Contains(r.Stdout(),
		`repository_id
repo_b
`)

	// Test SQL queries. This should be a different test, but since starting
	// the environment takes a long time, it is bundled together here to speed up
	// the tests
	s.testSQL()

	client, err := newSupersetClient()
	req.NoError(err)

	// Test the list of dashboards created in superset
	s.T().Run("dashboard-list", func(t *testing.T) {
		req := require.New(t)

		links, err := client.dashboards()
		req.NoError(err)

		s.Equal([]string{
			`<a href="/superset/dashboard/1/">Overview</a>`,
			`<a href="/superset/dashboard/welcome-local/">Welcome</a>`,
		}, links)
	})

	// Test gitbase queries through superset
	s.T().Run("superset-gitbase", func(t *testing.T) {
		req := require.New(t)

		rows, err := client.gitbase("select * from repositories")
		req.NoError(err)

		s.Equal([]map[string]interface{}{
			{"repository_id": "repo_b"},
		}, rows)
	})

	// Test bblfsh queries through superset
	s.T().Run("superset-bblfsh", func(t *testing.T) {
		req := require.New(t)

		lang, err := client.bblfsh("hello.js", `console.log("hello");`)
		req.NoError(err)
		req.Equal("javascript", lang)
	})

	// Test gitbase can connect to bblfsh with a SQL query that uses UAST
	s.T().Run("gitbase-bblfsh", func(t *testing.T) {
		req := require.New(t)

		rows, err := client.gitbase(
			`SELECT UAST('console.log("hello");', 'javascript') AS uast`)
		req.NoError(err)

		req.Len(rows, 1)
		req.NotEmpty(rows[0]["uast"])
	})
}

func (s *InitLocalTestSuite) initGitRepo(path string) {
	s.T().Helper()

	err := os.MkdirAll(path, os.ModePerm)
	s.Require().NoError(err)

	cmd := exec.Command("git", "init", path)
	err = cmd.Run()
	s.Require().NoError(err)
}
