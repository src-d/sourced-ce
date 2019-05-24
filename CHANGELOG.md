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

#### Bug Fixes

- The gitbase indexes are now persistent ([#65](https://github.com/src-d/superset-compose/issues/65)).

</details>

## [v0.0.1](https://github.com/src-d/superset-compose/releases/tag/v0.0.1) - 2019-05-16

### sandbox-ce

Initial release. It includes a `sandbox-ce` command with the sub commands `install`, `stop`, `start`, `prune`.

This binary is a wrapper for Docker Compose, and requires you to download the `docker-compose.yml` file from this repository.

### `srcd/superset` docker image

The `srcd/superset` docker image is based on Superset 0.32, and contains the following additions:
- an extra tab, UAST, to explore bblfsh parsing results.
- SQL Lab contains a modal dialog to visualize columns that contain UAST.
- source{d} branding.
- SQLAlchemy dependency upgraded to 1.3, for compatibility with gitbase ([#18](https://github.com/src-d/superset-compose/issues/18)).
- Backport an upstream fix for Hive Database connection ([#21](https://github.com/src-d/superset-compose/issues/21)).
