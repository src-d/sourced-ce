package workdir

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FactorySuite struct {
	suite.Suite
}

func TestFactorySuite(t *testing.T) {
	suite.Run(t, &FactorySuite{})
}

func (s *FactorySuite) TestEmptyEnv() {
	// Test for https://github.com/src-d/sourced-ce/pull/212,
	// the change from ${VARIABLE:-default} to ${VARIABLE-default} means that
	// the env file should not contain any empty `VAR=\n`
	envf := envFile{}

	contents, err := envf.MarshalEnv()
	s.NoError(err)
	s.NotNil(contents)
	s.NotContains(string(contents), `=
`)
}
