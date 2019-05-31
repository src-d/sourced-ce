# source{d} CE

## Contents

- [Usage](#usage)
  - [Commands](#commands)
  - [Docker Compose](#docker-compose)

## Usage

The source{d} Sandbox CE is deployed as Docker containers, using Docker Compose.

This repository provides the `sandbox-ce` binary as a wrapper to manage the Docker Compose files and containers easily. Moreover, `sandbox-ce` does not require a local installation of Docker Compose, if it is not found it will be deployed inside a container.

You may also choose to manage the containers yourself with the `docker-compose.yml` file included in this repository.

## Defaults

Default login: `admin`

Default password: `admin`

### Commands

Go to the [releases page](https://github.com/src-d/sourced-ce/releases) and download the `sandbox-ce` binary for your system. You will also need to download the `docker-compose.yml` file included in the release assets.

Please make sure you run `sandbox-ce` commands in the same directory where you placed the `docker-compose.yml` file.

#### Install

```
GITBASE_REPOS_DIR=/some/path sandbox-ce install
```

This will create admin user and initialize superset.

Currently it would print some exceptions, just ignore them. It will be fixed after [the issue](https://github.com/src-d/gitbase/issues/808) is resolved.

After the initialization, the components superset, gitbase, bblfsh and other dependencies will be started.

It will automatically open WebUI. Use login: `admin` and password `admin` to access it.

If the UI wasn't opened automatically, you can access it going to http://localhost:8088


#### Start

```
GITBASE_REPOS_DIR=/some/path sandbox-ce start
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

### Docker Compose

As an alternative to `sandbox-ce` you can download the compose file and use the `docker-compose` command. Go to the [releases page](https://github.com/src-d/sourced-ce/releases) to download the `docker-compose.yml` file included in the release assets.

Then you can start the containers like follows:

```shell
GITBASE_REPOS_DIR=/some/path docker-compose up
```
