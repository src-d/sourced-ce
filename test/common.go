// +build integration

package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TODO (carlosms) this could be build/bin, workaround for https://github.com/src-d/ci/issues/97
var srcdBin = fmt.Sprintf("../build/sourced-ce_%s_%s/sourced", runtime.GOOS, runtime.GOARCH)

func init() {
	if os.Getenv("SOURCED_BIN") != "" {
		srcdBin = os.Getenv("SOURCED_BIN")
	}
}

type IntegrationSuite struct {
	suite.Suite
	*Commander
	TestDir string
}

func (s *IntegrationSuite) SetupTest() {
	testDir, err := ioutil.TempDir("", strings.Replace(s.T().Name(), "/", "_", -1))
	if err != nil {
		log.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		testDir, err = filepath.EvalSymlinks(testDir)
		if err != nil {
			log.Fatal(err)
		}
	}

	s.TestDir = testDir
	s.Commander = &Commander{bin: srcdBin, sourcedDir: filepath.Join(s.TestDir, "sourced")}

	// Instead of downloading the compose file, copy the local file
	err = os.MkdirAll(filepath.Join(s.sourcedDir, "compose-files", "local"), os.ModePerm)
	s.Require().NoError(err)

	p, _ := filepath.Abs(filepath.FromSlash("../docker-compose.yml"))
	bytes, err := ioutil.ReadFile(p)
	s.Require().NoError(err)
	err = ioutil.WriteFile(filepath.Join(s.sourcedDir, "compose-files", "local", "docker-compose.yml"), bytes, 0644)
	s.Require().NoError(err)

	r := s.RunCommand("compose", "set", "local")
	s.Require().NoError(r.Error, r.Combined())
}

func (s *IntegrationSuite) TearDownTest() {
	// don't run prune on failed test to help debug. But stop the containers
	// to avoid port conflicts in the next test
	if s.T().Failed() {
		s.RunCommand("stop")
		s.T().Logf("Test failed. sourced data dir left in %s", s.TestDir)
		s.T().Logf("Probably there are also docker volumes left untouched")
		return
	}

	s.RunCommand("prune", "--all")

	os.RemoveAll(s.TestDir)
}

func (s *IntegrationSuite) testSQL() {
	testCases := []string{
		"show tables",
		"show tables;",
		" show tables ; ",
		"/* comment */ show tables;",
		`/* multi line
			comment */
			show tables;`,
	}

	showTablesOutput :=
		`Table
blobs
commit_blobs
commit_files
commit_trees
commits
files
ref_commits
refs
remotes
repositories
tree_entries
`

	for _, query := range testCases {
		s.T().Run(query, func(t *testing.T) {
			assert := assert.New(t)

			r := s.RunCommand("sql", query)
			assert.NoError(r.Error, r.Combined())

			assert.Contains(r.Stdout(), showTablesOutput)
		})
	}
}
