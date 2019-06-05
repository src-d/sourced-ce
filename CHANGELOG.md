# Changelog

## [Unreleased]
<details>
  <summary>
    Changes that have landed in master but are not yet released.
    Click to see more.
  </summary>

Initial release. It includes a `docker-compose.yml` file to deploy source{d} CE locally, and a `sourced` installer command.

The `sourced` binary is a wrapper for Docker Compose that downloads the `docker-compose.yml` file from this repository, and includes the following sub commands:

- `install`: Install and initialize containers
- `status`: Shows status of the components
- `stop`: Stop running containers
- `start`: Start stopped containers
- `web`: Open the web interface in your browser
- `prune`: Stop and remove containers and resources
- `workdirs` List working directories
- `compose`: Manage docker compose files
  - `download`: Download docker compose files
  - `list`: List the downloaded docker compose files
  - `set`: Set the active docker compose file

</details>
