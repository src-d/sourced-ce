# Initialize source{d} Community Edition

**source{d} Community Edition (CE)** is deployed as Docker containers, using Docker Compose.

This tool is a wrapper for Docker Compose to manage the compose files and containers easily. Moreover, `sourced` does not require a local installation of Docker Compose, if it is not found it will be deployed inside a container.

You may also choose to manage the containers yourself with the `docker-compose.yml` file included in this repository.

### Defaults

- Default login: `admin`
- Default password: `admin`


### Initialization

**source{d} CE** can be initialized from 2 different data sources: local Git repositories, or GitHub organizations.

Please note that you have to choose one data source to initialize **source{d} CE**, but you can have more than one isolated environment, and they can have different sources. See the [Working With Multiple Data Sets](#working-with-multiple-data-sets) section below for more details.

#### From GitHub Organizations

When using GitHub organizations to populate the **source{d} CE** database you only need to provide a list of organization names, and a [GitHub personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/). The token should have the following scopes enabled:

- [x] `repo`  Full control of private repositories
- [ ] `admin:org`  Full control of orgs and teams, read and write org projects
  - [ ] `write:org`  Read and write org and team membership, read and write org projects
  - [x] `read:org`  Read org and team membership, read org projects

Use this command to initialize:

```shell
sourced init orgs --token <token> src-d,bblfsh
```

It will automatically open the web UI. Use login: `admin` and password `admin` to access it.

If the UI wasn't opened automatically, use `sourced web` or visit http://localhost:8088.

#### From Local Repositories

```
sourced init local /path/to/repositories
```

This will initialize **source{d} CE** to analyze the given Git repositories.

The argument must point to a directory containing one or more Git repositories. The repositories will be found recursively. If no argument is given, the current directory will be used.

It will automatically open the web UI. Use login: `admin` and password `admin` to access it.

If the UI wasn't opened automatically, use `sourced web` or visit http://localhost:8088.
