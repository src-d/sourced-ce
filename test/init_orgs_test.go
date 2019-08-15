// +build integration

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
	req := s.Require()

	// TODO will need to change with https://github.com/src-d/sourced-ce/issues/144
	r := s.RunCommand("status", "workdirs")
	req.Error(r.Error)

	r = s.RunCommand("init", "orgs", "golang-migrate")
	req.NoError(r.Error, r.Combined())

	r = s.RunCommand("status", "workdirs")
	req.NoError(r.Error, r.Combined())

	req.Equal("* golang-migrate\n", r.Stdout())

	checkGhsync(req, 1)
	checkGitcollector(req, 1)

	// Check gitbase can also see the repositories
	r = s.RunCommand("sql", "select * from repositories where repository_id='github.com/golang-migrate/migrate'")
	req.NoError(r.Error, r.Combined())

	req.Contains(r.Stdout(),
		`repository_id
github.com/golang-migrate/migrate
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
			`<a href="/superset/dashboard/2/">Welcome</a>`,
			`<a href="/superset/dashboard/3/">Collaboration</a>`,
		}, links)
	})

	// Test gitbase queries through superset
	s.T().Run("superset-gitbase", func(t *testing.T) {
		req := require.New(t)

		rows, err := client.gitbase("select * from repositories")
		req.NoError(err)

		s.Equal([]map[string]interface{}{
			{"repository_id": "github.com/golang-migrate/migrate"},
		}, rows)
	})

	// Test metadata queries through superset
	s.T().Run("superset-metadata", func(t *testing.T) {
		req := require.New(t)

		rows, err := client.metadata("select * from organizations")
		req.NoError(err)
		req.Len(rows, 1)

		s.Equal("golang-migrate", rows[0]["login"])
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
