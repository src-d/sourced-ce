// +build integration

package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ComposeTestSuite struct {
	IntegrationSuite
}

func TestComposeTestSuite(t *testing.T) {
	itt := ComposeTestSuite{}
	suite.Run(t, &itt)
}

func (s *ComposeTestSuite) TestListComposeFiles() {
	r := s.RunCommand("compose", "list")
	s.Contains(r.Stdout(), "[0]* local")
}

func (s *ComposeTestSuite) TestSetComposeFile() {
	r := s.RunCommand("compose", "set", "0")
	s.Contains(r.Stdout(), "Active docker compose file was changed successfully")

	r = s.RunCommand("compose", "list")
	s.Contains(r.Stdout(), "[0]* local")
}

func (s *ComposeTestSuite) TestSetComposeFilIndexOutOfRange() {
	r := s.RunCommand("compose", "set", "5")
	s.Contains(r.Stderr(), "Index is out of range, please check the output of 'sourced compose list'")
}

func (s *ComposeTestSuite) TestSetComposeNotFound() {
	r := s.RunCommand("compose", "set", "NotFound")
	s.Error(r.Error)
}

func (s *ComposeTestSuite) TestSetComposeFilesWithStringIndex() {
	r := s.RunCommand("compose", "set", "local")
	s.Contains(r.Stdout(), "Active docker compose file was changed successfully")

	r = s.RunCommand("compose", "list")
	s.Contains(r.Stdout(), "[0]* local")
}
