# Contribution Guidelines

As all source{d} projects, this project follows the
[source{d} Contributing Guidelines](https://github.com/src-d/guide/blob/master/engineering/documents/CONTRIBUTING.md).


# Additional Contribution Guidelines

In addition to the [source{d} Contributing Guidelines](https://github.com/src-d/guide/blob/master/engineering/documents/CONTRIBUTING.md), this project follows the following guidelines.


## Changelog

This project lists the important changes between releases in the [`CHANGELOG.md`](../CHANGELOG.md) file.

If you open a PR, you should also add a brief summary in the `CHANGELOG.md` mentioning the new feature, change or bugfix that you proposed.


## How To Restore Dashboards and Charts to Defaults

The official way to restore **source{d} CE** to its initial state, is to remove the running components with
`sourced prune --all`, and then init again with `sourced init`.

In some circumstances you need to restore only the state modified from the UI (charts, dashboards, saved queries, users,
roles, etcetera), using the default ones for the version of **source{d} CE** that you're currently using, and preserve
the repositories and metadata fetched from GitHub organizations.

To do so, you only need to delete the docker volume containing the PostgreSQL database, and restart **source{d} CE**.
It can be done following these steps if you already have [Docker Compose](https://docs.docker.com/compose/) installed:

```shell
$ cd ~/.sourced/workdirs/__active__
$ source .env
$ docker-compose stop postgres
$ docker-compose rm -f postgres
$ ENV_PREFIX=`awk '{print tolower($0)}' <<< ${COMPOSE_PROJECT_NAME}`
$ docker volume rm ${ENV_PREFIX}_postgres
$ docker-compose up -d postgres
$ docker-compose exec -u superset sourced-ui bash -c 'sleep 10s && python bootstrap.py'
```
