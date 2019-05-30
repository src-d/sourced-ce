# Changelog

## [Unreleased]
<details>
  <summary>
    Changes that have landed in master but are not yet released.
    Click to see more.
  </summary>

### sandbox-ce

#### New Features

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

#### Bug Fixes

- The gitbase indexes are now persistent ([#65](https://github.com/src-d/superset-compose/issues/65)).
- Add proxy for UAST tab ([#47](https://github.com/src-d/superset-compose/issues/47)).
- Download docker-compose version defined in cli ([#24](https://github.com/src-d/superset-compose/issues/24)).

### `srcd/superset` Docker Image

#### New Features

- Loading of default dashboards on bootstrap ([#71](https://github.com/src-d/superset-compose/issues/71)).
- Update gitbase to v0.20.0 ([#81](https://github.com/src-d/superset-compose/issues/81)).
- Update bblfsh to v2.14.0 ([#81](https://github.com/src-d/superset-compose/issues/81)).
- Use sparksql instead of mysql for enterprise version ([#79](https://github.com/src-d/superset-compose/issues/79)).
- Cancel query to database on stop ([#35](https://github.com/src-d/superset-compose/issues/35)).


</details>

## [v0.0.1](https://github.com/src-d/superset-compose/releases/tag/v0.0.1) - 2019-05-16

### sandbox-ce

Initial release. It includes a `sandbox-ce` command with the sub commands `install`, `stop`, `start`, `prune`.

This binary is a wrapper for Docker Compose, and requires you to download the `docker-compose.yml` file from this repository.

### `srcd/superset` Docker Image

The `srcd/superset` docker image is based on Superset 0.32, and contains the following additions:
- an extra tab, UAST, to explore bblfsh parsing results.
- SQL Lab contains a modal dialog to visualize columns that contain UAST.
- source{d} branding.
- SQLAlchemy dependency upgraded to 1.3, for compatibility with gitbase ([#18](https://github.com/src-d/superset-compose/issues/18)).
- Backport an upstream fix for Hive Database connection ([#21](https://github.com/src-d/superset-compose/issues/21)).
