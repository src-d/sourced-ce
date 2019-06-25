# Changelog

## [Unreleased]
<details>
  <summary>
    Changes that have landed in master but are not yet released.
    Click to see more.
  </summary>

Initial release. It includes a `docker-compose.yml` file to deploy source{d} CE locally, and a `sourced` installer command.

The `sourced` binary is a wrapper for Docker Compose that downloads the `docker-compose.yml` file from this repository, and includes the following sub commands:

- `init`: Initialize source{d} to work on local or Github orgs datasets
  - `local`: Initialize source{d} to analyze local repositories
  - `orgs`: Initialize source{d} to analyze GitHub organizations
- `status`: Show the status of all components
- `stop`: Stop any running components
- `start`: Start any stopped components
- `web`: Open the web interface in your browser
- `sql`: Open a MySQL client connected to a SQL interface for Git
- `prune`: Stop and remove components and resources
- `workdirs` List all working directories
- `compose`: Manage source{d} docker compose files
  - `download`: Download docker compose files
  - `list`: List the downloaded docker compose files
  - `set`: Set the active docker compose file

</details>
