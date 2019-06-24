#  source{d} Community Editon Architecture

**source{d} Community Editon** provides a frictionless experience for trying
source{d} for Code Analysis.


## Technical Architecture

The `sourced` binary, a single CLI binary [written in Go](../../cmd/sourced/main.go),
is the user's main interaction mechanism with **source{d} CE**.
It is also the only piece (other than Docker) that the user will need to explicitly
download on their machine to get started.

The `sourced` binary manages the different installed environments and its
configurations, acting as a wrapper of Docker Compose.

The whole architecture is based on Docker containers, orchestrated by Docker Compose
and managed by `sourced`.


## Docker Set Up

In order to make this work in the easiest way, there was made some design decisions:

### Isolated environments.

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

Some environment services can be accessed from the outside, using the port exposed
for this purpose:
- `bblfsh`: `9432`
- `gitbase`: `3306`
- `metadatadb`: `5433`
- `sourced-ui`: `8088`

## Persistence

To prevent losing data when restarting services, or upgrading containers, its data
is stored in volumes. These volumes also share the same prefix with the containers
in the same environment, e.g. `srcd-c3jjlwq_gitbase_repositories`.

These are the most relevant volumes:
- `gitbase_repositories`, stores the repositories to be analyzed,
- `gitbase_indexes`, stores the gitbases indexes,
- `metadata`, stores the metadata from GitHub pull requests, issues, users...
- `postgres`, stores the dashboards and charts used by the web interface.
