# source{d} Community Edition (CE)

## Contents

- [Installation](#installation)
  - [Install Docker](#install-docker)
  - [Install source{d} Community Edition](#install-source-d-community-edition)
- [Usage](#usage)
  - [Defaults](#defaults)
  - [Initialization](#initialization)
  - [Commands](#commands)
  - [Working With Multiple Data Sets](#working-with-multiple-data-sets)
- [Docker Compose](#docker-compose)

## Installation

### Install Docker

_Please note that Docker Toolbox is not supported neither for Windows nor for macOS. In case that you're running Docker Toolbox, please consider updating to newer Docker Desktop for Mac or Docker Desktop for Windows._

Follow the instructions based on your OS:

- [Docker for Ubuntu Linux](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-docker-ce-1)
- [Docker for Arch Linux](https://wiki.archlinux.org/index.php/Docker#Installation)
- [Docker for macOS](https://store.docker.com/editions/community/docker-ce-desktop-mac)
- [Docker Desktop for Windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows). Make sure to read the [system requirements for Docker on Windows](https://docs.docker.com/docker-for-windows/install/).


### Install source{d} Community Edition

Download the **[latest release](https://github.com/src-d/sourced-ce/releases/latest)** for your Linux, macOS (Darwin) or Windows.

#### on Linux or macOS

Extract `sourced` binary from the release you downloaded, and move it into your bin folder to make it executable from any directory:

```bash
$ tar -xvf path/to/sourced-ce_REPLACE-VERSION_REPLACE-OS_amd64.tar.gz
$ sudo mv path/to/sourced-ce_REPLACE-OS_amd64/sourced /usr/local/bin/
```

#### on Windows

*Please note that from now on we assume that the commands are executed in `powershell` and not in `cmd`.*

Create a directory for `sourced.exe` and add it to your `$PATH`, running these commands in a powershell as administrator:
```powershell
mkdir 'C:\Program Files\sourced'
# Add the directory to the `%path%` to make it available from anywhere
setx /M PATH "$($env:path);C:\Program Files\sourced"
# Now open a new powershell to apply the changes
```

Extract the `sourced.exe` executable from the release you downloaded, and copy it into the directory you created in the previous step:
```powershell
mv \path\to\sourced-ce_windows_amd64\sourced.exe 'C:\Program Files\sourced'
```

## Usage

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
sourced orgs init --token <token> src-d,bblfsh
```

It will automatically open the web UI. Use login: `admin` and password `admin` to access it.

If the UI wasn't opened automatically, use `sourced web` or visit http://localhost:8088.

#### From Local Repositories

```
sourced init /path/to/repositories
```

This will initialize **source{d} CE** to analyze the given Git repositories.

The argument must point to a directory containing one or more Git repositories. The repositories will be found recursively. If no argument is given, the current directory will be used.

It will automatically open the web UI. Use login: `admin` and password `admin` to access it.

If the UI wasn't opened automatically, use `sourced web` or visit http://localhost:8088.

### Commands

#### Start

```
sourced start
```

This will start all the components previously initialized with `init`.

#### Stop

```
sourced stop
```

This will stop all running containers without removing them. They can be started again with 'start'.

#### Prune

```
sourced prune
```

Stops containers and removes containers, networks, and volumes created by `init`.
Images are not deleted unless you specify the `--images` flag.

If you want to completely uninstall `sourced` you may want to delete the `~/.sourced` directory.

#### Status

```
sourced status
```

Shows status of the components.

#### Web

```
sourced web
```

Opens the web interface in your browser.

#### SQL

```
sourced sql
```

Open a MySQL client connected to gitbase.

#### Workdirs

```
sourced workdirs
```

Lists previously initialized working directories.

#### Compose

```
sourced compose
```

Manage docker compose files in the `~/.sourced` directory with the following sub commands:

##### Download

```
sourced compose download
sourced compose download v0.0.1
sourced compose download master
sourced compose download https://raw.githubusercontent.com/src-d/sourced-ce/master/docker-compose.yml
```

Download docker compose files. By default the command downloads the file in `master`.

Use the `version` argument to choose a specific revision from the https://github.com/src-d/sourced-ce repository, or to set a URL to a `docker-compose.yml` file.

##### List

```
sourced compose
```

List the downloaded docker compose files.

##### Set

```
sourced compose set
```

Set the active docker compose file.

### Working With Multiple Data Sets

You can deploy more than one **source{d} CE** instance with different sets of organizations, or repositories, to analyze.

For example you may have initially started **source{d} CE** with the repositories in the `src-d` organization, with the command:
```
sourced orgs init --token <token> src-d
```

After a while you may want to analyze the data on another set of repositories. You can run `init` again with a different organization:
```
sourced orgs init --token <token> bblfsh
```

This command will stop any of the currently running containers, create an isolated environment for the new data, and create a new, clean deployment.

Please note that each path will have an isolated deployment. This means that for example any chart or dashboard created for the deployment for `src-d` will not be available to the new deployment for `bblfsh`.

Each isolated environment is persistent (unless you run `prune`). Which means that if you decide to re-deploy **source{d} CE** using the original organization:
```
sourced orgs init --token <token> src-d
```

You will get back to the previous state, and things like charts and dashboards will be restored.

These isolated environments also allow you to deploy **source{d} CE** using a local set of Git repositories. For example, if we wanted a third deployment to analyze repositories already existing in the `~/repos` directory, we just need to run `init` again:

```
sourced init ~/repos
```

If you are familiar with Docker Compose and you want more control over the underlying resources, you can explore the contents of your `~/.sourced` directory. There you will find a `docker-compose.yml` and `.env` files for each set of repositories used by `sourced init`.

## Docker Compose

As an alternative to `sourced` you can download the compose file and use the `docker-compose` command. Go to the [releases page](https://github.com/src-d/sourced-ce/releases) to download the `docker-compose.yml` file included in the release assets.

Then you can start the containers like follows:

```shell
GITBASE_REPOS_DIR=/some/path docker-compose up
```
