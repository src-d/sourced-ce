# Usage

### Init

```
GITBASE_REPOS_DIR=/some/path  docker-compose run --rm superset ./docker-init.sh
```

to create admin user and initialize superset.

Currently it would print some exceptions, just ignore them. It will be fixed after [the issue](https://github.com/src-d/gitbase/issues/808) is resolved.

### Run

```
GITBASE_REPOS_DIR=/some/path docker-compose up
```

to start superset, gitbase, bblfsh and other dependencies.

Wait a little and open http://localhost:8088

It would have gitbase datasource already.


# Development

## Setup local environment

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


## Build docker image

```
make build
```

Image name defined in Makefile and matches the one in docker-compose.

## Work with superset upstream

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
