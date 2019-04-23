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

### Add gitbase:

Go to Sources -> Databases, create a new database:

```
name: gitbase
uri: mysql://root@gitbase:3306/gitbase
need to check Allow DML for complex queries
need to check Asynchronous Query Execution or queries will timeout
```
