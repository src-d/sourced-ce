# Usage

### Init

```
docker-compose run --rm superset ./docker-init.sh
```

to create admin user and initialize superset.

### Run

```
GITBASE_REPOS_DIR=/some/path docker-compose up
```

to start superset, gitbase, bblfsh and other dependencies.

Wait a little and open http://localhost:8088

It would have gitbase datasource already.
