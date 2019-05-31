# source{d} Community Edition (CE)

## Contents

- [Usage](#usage)
  - [Defaults](#defaults)
  - [Commands](#commands)
  - [Working With Multiple Data Sets](#working-with-multiple-data-sets)
- [Docker Compose](#docker-compose)

## Usage

**source{d} Community Edition (CE)** is deployed as Docker containers, using Docker Compose.

This repository provides the `sandbox-ce` binary as a wrapper to manage the Docker Compose files and containers easily. Moreover, `sandbox-ce` does not require a local installation of Docker Compose, if it is not found it will be deployed inside a container.

You may also choose to manage the containers yourself with the `docker-compose.yml` file included in this repository.

### Defaults

- Default login: `admin`
- Default password: `admin`

### Commands

Go to the [releases page](https://github.com/src-d/sourced-ce/releases) and download the `sandbox-ce` binary for your system. You will also need to download the `docker-compose.yml` file included in the release assets.

Please make sure you run `sandbox-ce` commands in the same directory where you placed the `docker-compose.yml` file.

#### Install

```
sandbox-ce install /path/to/repositories
```

This will initialize **source{d} CE** to analyze the given Git repositories.

The argument must point to a directory containing one or more Git repositories. The repositories will be found recursively. If no argument is given, the current directory will be used.

It will automatically open the web UI. Use login: `admin` and password `admin` to access it.

If the UI wasn't opened automatically, use `sandbox-ce web` or visit http://localhost:8088.

#### Start

```
sandbox-ce start
```

This will start all the components previously installed.

#### Stop

```
sandbox-ce stop
```

This will stop all running components.

#### Prune

```
sandbox-ce prune
```

Stops containers and removes containers, networks, and volumes created by `install`.
Images are not deleted unless you specify the `--images` flag.

If you want to completely uninstall `sandbox-ce` you may want to delete the `~/.srcd` directory.

#### Status

```
sandbox-ce status
```

Shows status of the components.

#### Web

```
sandbox-ce web
```

Opens the web interface in your browser.

#### Compose

```
sandbox-ce compose
```

Manage docker compose files in the `~/.srcd` directory with the following sub commands:

##### Download

```
sandbox-ce compose download
sandbox-ce compose download v0.0.1
sandbox-ce compose download master
sandbox-ce compose download https://raw.githubusercontent.com/src-d/sourced-ce/master/docker-compose.yml
```

Download docker compose files. By default the command downloads the file in `master`.

Use the `version` argument to choose a specific revision from the https://github.com/src-d/sourced-ce repository, or to set a URL to a docker-compose.yml file.

##### List

```
sandbox-ce compose
```

List the downloaded docker compose files.

##### Set

```
sandbox-ce compose set
```

Set the active docker compose file.

### Working With Multiple Data Sets

You can deploy more than one **source{d} CE** instance with different sets of repositories to analyze.

For example you may have initially started **source{d} CE** with the repositories in `~/repos`, with the command:
```
sandbox-ce install ~/repos
```

After a while you may want to analyze the data on another set of repositories. You can run `install` again with a different path:
```
sandbox-ce install ~/go/src/github.com/src-d
```

This command will stop any of the currently running containers, create an isolated environment for the new repositories path, and create a new, clean deployment.

Please note that each path will have an isolated deployment. This means that for example any chart or dashboard created for the deployment using `~/repos` will not be available to the new deployment for `~/go/src/github.com/src-d`.

Each isolated environment is persistent (unless you run `prune`). Which means that if you decide to re-deploy **source{d} CE** using the original set of repositories:
```
sandbox-ce install ~/repos
```

You will get back to the previous state, and things like charts and dashboards will be restored.

If you are familiar with Docker Compose and you want more control over the underlying resources, you can explore the contents of your `~/.srcd` directory. There you will find a `docker-compose.yml` and `.env` files for each set of repositories used by `sandbox-ce install`.

## Docker Compose

As an alternative to `sandbox-ce` you can download the compose file and use the `docker-compose` command. Go to the [releases page](https://github.com/src-d/sourced-ce/releases) to download the `docker-compose.yml` file included in the release assets.

Then you can start the containers like follows:

```shell
GITBASE_REPOS_DIR=/some/path docker-compose up
```
