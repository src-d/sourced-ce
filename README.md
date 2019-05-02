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
