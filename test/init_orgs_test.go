package test

import (
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InitOrgsTestSuite struct {
	IntegrationSuite
}

func TestInitOrgsTestSuite(t *testing.T) {
	itt := InitOrgsTestSuite{}

	if os.Getenv("SOURCED_GITHUB_TOKEN") == "" {
		t.Skip("SOURCED_GITHUB_TOKEN is not set")
		return
	}

	suite.Run(t, &itt)
}

func checkGhsync(require *require.Assertions, repos int) {
	connStr := "user=metadata password=metadata dbname=metadata port=5433 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(err)
	defer db.Close()

	var id int
	var org, entity string
	var done, failed, total int

	// try for 2 minutes
	for i := 0; i < 24; i++ {
		time.Sleep(5 * time.Second)

		row := db.QueryRow("SELECT * FROM status")
		err = row.Scan(&id, &org, &entity, &done, &failed, &total)
		if err == sql.ErrNoRows {
			continue
		}
		require.NoError(err)

		if done == repos {
			break
		}
	}

	require.Equal(repos, done,
		"id = %v, org = %v, entity = %v, done = %v, failed = %v, total = %v",
		id, org, entity, done, failed, total)
}

func checkGitcollector(require *require.Assertions, repos int) {
	connStr := "user=superset password=superset dbname=superset port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(err)
	defer db.Close()

	var org string
	var discovered, downloaded, updated, failed int

	// try for 2 minutes
	for i := 0; i < 24; i++ {
		time.Sleep(5 * time.Second)

		row := db.QueryRow("SELECT * FROM gitcollector_metrics")
		err = row.Scan(&org, &discovered, &downloaded, &updated, &failed)
		if err == sql.ErrNoRows {
			continue
		}
		require.NoError(err)

		if downloaded == repos {
			break
		}
	}

	require.Equal(repos, downloaded,
		"org = %v, discovered = %v, downloaded = %v, updated = %v, failed = %v",
		org, discovered, downloaded, updated, failed)
}

func (s *InitOrgsTestSuite) TestOneOrg() {
	require := s.Require()

	r := s.RunCommand("init", "orgs", "golang-migrate")
	require.NoError(r.Error, r.Combined())

	r = s.RunCommand("workdirs")
	require.NoError(r.Error, r.Combined())

	require.Equal("* golang-migrate\n", r.Stdout())

	checkGhsync(require, 1)
	checkGitcollector(require, 1)

	// Check gitbase can also see the repositories
	r = s.RunCommand("sql", "select * from repositories where repository_id='github.com/golang-migrate/migrate'")
	require.NoError(r.Error, r.Combined())

	require.Contains(r.Stdout(),
		`repository_id
github.com/golang-migrate/migrate
`)
}
