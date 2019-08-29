package workdir

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite

	originSrcdDir string
}

func TestFactorySuite(t *testing.T) {
	suite.Run(t, &FactorySuite{})
}

func (s *FactorySuite) BeforeTest(suiteName, testName string) {
	s.originSrcdDir = os.Getenv("SOURCED_DIR")

	// on macOs os.TempDir returns symlink and tests fails
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	srdPath := path.Join(tmpDir, testName)
	err := os.MkdirAll(srdPath, os.ModePerm)
	s.Nil(err)

	os.Setenv("SOURCED_DIR", srdPath)
}

func (s *FactorySuite) AfterTest(suiteName, testName string) {
	os.RemoveAll(os.Getenv("SOURCED_DIR"))
	os.Setenv("SOURCED_DIR", s.originSrcdDir)
}

func (s *FactorySuite) TestInitLocal() {
	reposdir := "some-dir"
	wd, err := InitLocal(reposdir)
	s.Nil(err)
	s.Equal(Local, wd.Type)
	s.Equal(reposdir, wd.Name)

	// check docker-compose.yml exists
	composeYmlPath := path.Join(wd.Path, "docker-compose.yml")
	_, err = os.Stat(composeYmlPath)
	s.Nil(err)

	// check .env file
	envPath := path.Join(wd.Path, ".env")
	_, err = os.Stat(envPath)
	s.Nil(err)

	envf := envFile{}
	s.Nil(readEnvFile(encodeDirName(reposdir), "local", &envf))

	s.Equal(reposdir, envf.GitbaseVolumeSource)
	s.False(envf.NoForks)
}

func (s *FactorySuite) TestInitOrgs() {
	orgs := []string{"org2", "org1"}
	name := "org1,org2"
	token := "some-token"
	wd, err := InitOrgs(orgs, token, true)
	s.Nil(err)
	s.Equal(Orgs, wd.Type)
	s.Equal(name, wd.Name)

	// check docker-compose.yml exists
	composeYmlPath := path.Join(wd.Path, "docker-compose.yml")
	_, err = os.Stat(composeYmlPath)
	s.Nil(err)

	// check .env file
	envPath := path.Join(wd.Path, ".env")
	_, err = os.Stat(envPath)
	s.Nil(err)

	envf := envFile{}
	s.Nil(readEnvFile(encodeDirName(name), "orgs", &envf))

	s.Equal("gitbase_repositories", envf.GitbaseVolumeSource)
	s.Equal(orgs, envf.GithubOrganizations)
	s.Equal(token, envf.GithubToken)
	s.False(envf.NoForks)
}

func (s *FactorySuite) TestReInitForksOrgs() {
	orgs := []string{"org2", "org1"}
	_, err := InitOrgs(orgs, "", false)
	s.Nil(err)

	_, err = InitOrgs(orgs, "", true)
	s.EqualError(err, "initialization failed: workdir was previously initialized with a different value for forks support")
}
