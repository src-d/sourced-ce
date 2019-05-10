# Sandbox CE

## Contents

- [Commands](#commands)
- [Development](#development)

## Commands

### Install

```
GITBASE_REPOS_DIR=/some/path sandbox-ce install
```

This will create admin user and initialize superset.

Currently it would print some exceptions, just ignore them. It will be fixed after [the issue](https://github.com/src-d/gitbase/issues/808) is resolved.

After the initialization, the components superset, gitbase, bblfsh and other dependencies will be started. Wait a little and open http://localhost:8088

It would have gitbase datasource already.


### Start

```
GITBASE_REPOS_DIR=/some/path sandbox-ce start
```

This will start all the components previously installed.

### Stop

```
sandbox-ce stop
```

This will stop all running components.

### Prune

```
sandbox-ce prune
```

This will remove all containers and related resources such as network and volumes.


## Development

### Setup local environment

Run dependencies using docker-compose:
```
docker-compose up gitbase bblfsh-web
```

Update superset directory:

```
make patch-dev
```

Enter into `superset` directory:
```
cd superset
```

Follow original superset instructions for [Flask server](https://github.com/apache/incubator-superset/blob/release--0.32/CONTRIBUTING.md#flask-server) and [Frontend assets](https://github.com/apache/incubator-superset/blob/release--0.32/CONTRIBUTING.md#frontend-assets)


### Build docker image

```
VERSION=latest make superset-build
```

Image name defined in Makefile and matches the one in docker-compose.

### Work with superset upstream

Superset version which we are based on defined in Makefile.

To see which files are patched compare to upstream, run:

```
make superset-diff-stat
```

To see diff with upstream, run:

```
make superset-diff
```


To merge updated upsteam into subdirectory:

```
make superset-merge
```
