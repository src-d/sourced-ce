# Changelog

## [Unreleased]
<details>
  <summary>
    Changes that have landed in master but are not yet released.
    Click to see more.
  </summary>

</details>

## [v0.14.0](https://github.com/src-d/sourced-ce/releases/tag/v0.14.0) - 2019-08-07

Initial release of **source{d} Community Edition (CE)**, the data platform for your software development life cycle.

The `sourced` binary is a wrapper for Docker Compose that downloads the `docker-compose.yml` file from this repository, and includes the following sub commands:

- `init`: Initialize source{d} to work on local or GitHub organization datasets
  - `local`: Initialize source{d} to analyze local repositories
  - `orgs`: Initialize source{d} to analyze GitHub organizations
- `status`: Show the status of all components
- `stop`: Stop any running components
- `start`: Start any stopped components
- `logs`: Show logs from components
- `web`: Open the web interface in your browser
- `sql`: Open a MySQL client connected to a SQL interface for Git
- `prune`: Stop and remove components and resources
- `workdirs` List all working directories
- `compose`: Manage source{d} docker compose files
  - `download`: Download docker compose files
  - `list`: List the downloaded docker compose files
  - `set`: Set the active docker compose file
- `restart`: Update current installation according to the active docker compose file

### Known Issues

- On Windows, if you use `sourced init local` on a directory with a long path, you may encounter the following error:
  ```
  Can't find a suitable configuration file in this directory or any
  parent. Are you in the right directory?
  ```

  This is caused by the [`MAX_PATH` limitation on windows](https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#maximum-path-length-limitation). The only workaround is to move the target directory to a shorter path, closer to the root of your drive ([#191](https://github.com/src-d/sourced-ce/issues/191)).

- Linux only: Docker installed from snap packages is not supported, please install it following [the official documentation](https://docs.docker.com/install/) ([#78](https://github.com/src-d/sourced-ce/issues/78)).

### Upgrading

For internal releases we don't support upgrading. If you have a previous `sourced-ce` pre-release version installed, clean up all your data **before** downloading this release. This will delete everything, including the UI data for dashboards, charts, users, etc:

```shell
sourced prune --all
rm -rf ~/.sourced
```
