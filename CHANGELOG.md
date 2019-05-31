# Changelog

## [Unreleased]
<details>
  <summary>
    Changes that have landed in master but are not yet released.
    Click to see more.
  </summary>

### New Features

- Make the gitbase volume for repositories read-only ([#52](https://github.com/src-d/superset-compose/issues/52)).
- Add help messages to the `sandbox-ce` command ([#46](https://github.com/src-d/superset-compose/issues/46)).
- `sandbox-ce install` now starts the containers on detached mode in the background ([#44](https://github.com/src-d/superset-compose/issues/44)).
- New sub command `sandbox-ce web` to open the web UI in the browser ([#17](https://github.com/src-d/superset-compose/issues/17)).
- Add `restart` policy to gitbase and bblfsh containers ([#63](https://github.com/src-d/superset-compose/issues/63)).
- Download `docker-compose.yml` if it doesn't exist ([#51](https://github.com/src-d/superset-compose/pull/51)).
- New sub command `sandbox-ce status` to see components status ([#14](https://github.com/src-d/superset-compose/pull/14)).
- Create default user with default password on install ([#11](https://github.com/src-d/superset-compose/pull/11)).
- New flag `sandbox-ce prune --images` to delete Docker images too ([#69](https://github.com/src-d/superset-compose/issues/69)).
- New sub command `sandbox-ce compose` to download and manage Docker Compose files ([#13](https://github.com/src-d/superset-compose/issues/13)).
- Add spinner for long running commands ([#73](https://github.com/src-d/superset-compose/issues/73)).

### Bug Fixes

- The gitbase indexes are now persistent ([#65](https://github.com/src-d/superset-compose/issues/65)).
- Add proxy for UAST tab ([#47](https://github.com/src-d/superset-compose/issues/47)).
- Download docker-compose version defined in cli ([#24](https://github.com/src-d/superset-compose/issues/24)).

</details>

## v0.0.1 - 2019-05-16

Initial release. It includes a `sandbox-ce` command with the sub commands `install`, `stop`, `start`, `prune`.

This binary is a wrapper for Docker Compose, and requires you to download the `docker-compose.yml` file from this repository.
