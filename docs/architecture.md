#  source{d} Community Editon Architecture

**source{d} Community Editon** provides a frictionless experience for trying
source{d} for Code Analysis.


## Technical Architecture

The `sourced` binary, a single CLI binary [written in Go](../../cmd/sourced/main.go),
is the user's main interaction mechanism with **source{d} CE**.
It is also the only piece (other than Docker) that the user will need to explicitly
download on their machine to get started.

The `sourced` binary manages the different installed environments and their
configurations, acting as a wrapper of Docker Compose.

The whole architecture is based on Docker containers, orchestrated by Docker Compose
and managed by `sourced`.


## Components of source{d}

**source{d} CE** relies on different components to handle different use cases
and to cover different functionalities. Each component is implemented as a running
Docker container.

- `bblfsh`: parses source code into UASTs using [Babelfish](https://docs.sourced.tech/babelfish/);
you can learn more about it in our [Babelfish UAST guide](usage/bblfsh.md)
- `gitbase`: runs [gitbase](https://docs.sourced.tech/gitbase), a SQL database
interface to Git repositories.
- `gitcollector`: is responsible for fetching repositories from the organizations
used to initialize **source{d} CE**. It uses [gitcollector](https://github.com/src-d/gitcollector).
- `ghsync`: is responsible for fetching repository metadata from the organizations
used to initialize **source{d} CE**. It uses [ghsync](https://github.com/src-d/ghsync)
- `metadatadb`: runs the PostgreSQL database that stores the repositories
metadata (users, pull requests, issues...) extracted by `ghsync`.
- `postgres`: runs the PostgreSQL database that stores the state of the UI
(charts, dashboards, users, saved queries and such).
- `sourced-ui`: runs the **source{d} CE** Web Interface. This component queries
data from `bblfsh`, `gitbase`, `metadatadb` and `postgres`.

Some of these components can be accessed from the outside as described by
[Docker Networking section](#docker-networking).


## Docker Set Up

In order to make this work in the easiest way, some design decisions were made:

### Isolated Environments.

_Read more in [Working With Multiple Data Sets](../usage/multiple-datasets.md)_

Each dataset runs in an isolated environment, and only one environment can run
at the same time.
Each environment is defined by one `docker-compose.yml` and one `.env`, stored
in `~/.sourced`.

### Docker Naming

All the Docker containers from the same environment share its prefix:
`srcd-<HASH>_` followed by the name of the service running inside, e.g
`srcd-c3jjlwq_gitbase_1` and `srcd-c3jjlwq_bblfsh_1` will contain gitbase and
babelfish for the same environment.

### Docker Networking

In order to provide communication between the multiple containers started, all of
them are attached to the same single bridge network. The network name also has
the same prefix than the containers inside the same environment, e.g.
`srcd-c3jjlwq_default`.

Some environment services can be accessed from the outside, using their exposed
port and connection values:
- `bblfsh`:
    - port: `9432`
- `gitbase`:
    - port: `3306`
    - database: `gitbase`
    - user: `root`
- `metadatadb`:
    - port: `5433`
    - database: `metadata`
    - user: `metadata`
    - password: `metadata`
- `sourced-ui`:
    - port: `8088`

### Persistence

To prevent losing data when restarting services, or upgrading containers, the data
is stored in volumes. These volumes also share the same prefix with the containers
in the same environment, e.g. `srcd-c3jjlwq_gitbase_repositories`.

These are the most relevant volumes:
- `gitbase_repositories`, stores the repositories to be analyzed
- `gitbase_indexes`, stores the gitbases indexes
- `metadata`, stores the metadata from GitHub pull requests, issues, users...
- `postgres`, stores the dashboards and charts used by the web interface
