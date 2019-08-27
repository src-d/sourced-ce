# Changelog

## [Unreleased]

## [v0.15.1](https://github.com/src-d/sourced-ce/releases/tag/v0.15.1) - 2019-08-27

### Fixed

- Fix incompatibility of empty resource limits ([#227](https://github.com/src-d/sourced-ce/issues/227)).
- Fix incorrect value for `GITCOLLECTOR_LIMIT_CPU` in some cases ([#225](https://github.com/src-d/sourced-ce/issues/225)).
- Fix gitbase `LOG_LEVEL` environment variable in the compose file ([#228](https://github.com/src-d/sourced-ce/issues/228)).

### Removed

- Remove the `completion` sub-command on Windows, as it only works for bash ([#169](https://github.com/src-d/sourced-ce/issues/169)).

### Upgrading

Install the new `v0.15.1` binary, then run `sourced compose download`.

For an upgrade from `v0.15.0`, you just need to run `sourced restart` to re-deploy the containers.

For an upgrade from `v0.14.0`, please see the upgrade instructions in the release notes for `v0.15.0`.


## [v0.15.0](https://github.com/src-d/sourced-ce/releases/tag/v0.15.0) - 2019-08-21

### Components

- `srcd/sourced-ui` has been updated to [v0.5.0](https://github.com/src-d/sourced-ui/releases/tag/v0.5.0).
- `srcd/ghsync` has been updated to [v0.2.0](https://github.com/src-d/ghsync/releases/tag/v0.2.0).

### Added

- Add a monitoring of containers state while waiting for the web UI to open during initialization ([#147](https://github.com/src-d/sourced-ce/issues/147)).
- Exclude forks by default in `sourced init orgs`, adding a new flag `--with-forks` to include them if needed ([#109](https://github.com/src-d/sourced-ce/issues/109)).

### Changed

- Refactor of the `status` command ([#203](https://github.com/src-d/sourced-ce/issues/203)):
  - `sourced status components` shows the previous output of `sourced status`
  - `sourced status workdirs` replaces `sourced workdirs`
  - `sourced status config` shows the contents of the Docker Compose environment variables. This is useful, for example, to check if the active working directory was configured to include or skip forks when downloading the data from GitHub
  - `sourced status all` shows all of the above

### Upgrading

Install the new `v0.15.0` binary, then run `sourced compose download`. If you had a deployment running, you can re-deploy the containers with `sourced restart`.

Please note: `sourced-ui` contains fixes for the default dashboard charts that will only be visible when you run `sourced init local/org` with a new path or organization.
If you want to apply the new default dashboards over your existing deployment, you will need to run `sourced prune` (or `sourced prune --all`) and `sourced init local/org` again.

Important: running `prune` will delete all your current data and customizations, including charts or dashboards. You can choose to not `prune` your existing deployments, keeping you previous default dashboards and charts.

```shell
$ sourced version
sourced version v0.15.0 build 08-21-2019_08_30_24

$ sourced compose download
Docker compose file successfully downloaded to your ~/.sourced/compose-files directory. It is now the active compose file.
To update your current installation use `sourced restart`

$ sourced status workdirs
  bblfsh
* src-d

$ sourced prune --all
$ sourced init orgs src-d
$ sourced init orgs bblfsh
```

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
