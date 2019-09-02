package workdir

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// default limits depend on host system and can't be used in tests
func setResourceLimits(f *envFile) {
	f.GitcollectorLimitCPU = 2.0
	f.GitbaseLimitCPU = 1.5
	f.GitbaseLimitMem = 100
}

const localContent = `COMPOSE_PROJECT_NAME=srcd-dir-name
GITBASE_VOLUME_TYPE=bind
GITBASE_VOLUME_SOURCE=repo-dir
GITBASE_SIVA=
GITHUB_ORGANIZATIONS=
GITHUB_TOKEN=
NO_FORKS=
GITBASE_LIMIT_CPU=1.5
GITCOLLECTOR_LIMIT_CPU=2
GITBASE_LIMIT_MEM=100
`

const orgsContent = `COMPOSE_PROJECT_NAME=srcd-dir-name
GITBASE_VOLUME_TYPE=volume
GITBASE_VOLUME_SOURCE=gitbase_repositories
GITBASE_SIVA=true
GITHUB_ORGANIZATIONS=org1,org2
GITHUB_TOKEN=token
NO_FORKS=true
GITBASE_LIMIT_CPU=1.5
GITCOLLECTOR_LIMIT_CPU=2
GITBASE_LIMIT_MEM=100
`

const emptyContent = `COMPOSE_PROJECT_NAME=
GITBASE_VOLUME_TYPE=
GITBASE_VOLUME_SOURCE=
GITBASE_SIVA=
GITHUB_ORGANIZATIONS=
GITHUB_TOKEN=
NO_FORKS=
GITBASE_LIMIT_CPU=0
GITCOLLECTOR_LIMIT_CPU=0
GITBASE_LIMIT_MEM=0
`

func TestEnvMarshal(t *testing.T) {
	assert := assert.New(t)

	f := newLocalEnvFile("dir-name", "repo-dir")
	setResourceLimits(&f)
	b, err := f.MarshalEnv()
	assert.Nil(err)
	assert.Equal(localContent, strings.ReplaceAll(string(b), "\r\n", "\n"))

	f = newOrgEnvFile("dir-name", []string{"org1", "org2"}, "token", false)
	setResourceLimits(&f)
	b, err = f.MarshalEnv()
	assert.Nil(err)
	assert.Equal(orgsContent, strings.ReplaceAll(string(b), "\r\n", "\n"))

	f = envFile{}
	b, err = f.MarshalEnv()
	assert.Nil(err)
	assert.Equal(emptyContent, strings.ReplaceAll(string(b), "\r\n", "\n"))
}

func TestEnvUnmarshal(t *testing.T) {
	assert := assert.New(t)

	b := []byte(localContent)
	f := envFile{}
	assert.Nil(f.UnmarshalEnv(b))
	assert.Equal(envFile{
		ComposeProjectName:  "srcd-dir-name",
		GitbaseVolumeType:   "bind",
		GitbaseVolumeSource: "repo-dir",

		GitcollectorLimitCPU: 2.0,
		GitbaseLimitCPU:      1.5,
		GitbaseLimitMem:      100,
	}, f)

	b = []byte(orgsContent)
	f = envFile{}
	assert.Nil(f.UnmarshalEnv(b))
	assert.Equal(envFile{
		ComposeProjectName:  "srcd-dir-name",
		GitbaseVolumeType:   "volume",
		GitbaseVolumeSource: "gitbase_repositories",
		GitbaseSiva:         true,
		GithubOrganizations: []string{"org1", "org2"},
		GithubToken:         "token",
		NoForks:             true,

		GitcollectorLimitCPU: 2.0,
		GitbaseLimitCPU:      1.5,
		GitbaseLimitMem:      100,
	}, f)

	b = []byte("")
	f = envFile{}
	assert.Nil(f.UnmarshalEnv(b))

	b = []byte(" COMPOSE_PROJECT_NAME=srcd-dir-name  \n\n  GITBASE_VOLUME_TYPE=volume  ")
	f = envFile{}
	assert.Nil(f.UnmarshalEnv(b))
	assert.Equal(envFile{
		ComposeProjectName: "srcd-dir-name",
		GitbaseVolumeType:  "volume",
	}, f)

	b = []byte("UNKNOWN=1")
	f = envFile{}
	assert.Nil(f.UnmarshalEnv(b))
}
