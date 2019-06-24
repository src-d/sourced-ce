<a href="https://www.sourced.tech">
  <img src="docs/assets/sourced-community-edition.png" alt="source{d} Community Edition (CE)" height="120px" />
</a>

**source{d} Community Edition (CE) is the data platform for your software development life cycle.**

[![GitHub version](https://badge.fury.io/gh/src-d%2Fsourced-ce.svg)](https://github.com/src-d/sourced-ce/releases)
[![Build Status](https://travis-ci.com/src-d/sourced-ce.svg?branch=master)](https://travis-ci.com/src-d/sourced-ce)
![Beta](https://svg-badge.appspot.com/badge/stability/beta?color=D6604A)
[![Go Report Card](https://goreportcard.com/badge/github.com/src-d/sourced-ce)](https://goreportcard.com/report/github.com/src-d/sourced-ce)
[![GoDoc](https://godoc.org/github.com/src-d/sourced-ce?status.svg)](https://godoc.org/github.com/src-d/sourced-ce)

[Website](https://www.sourced.tech) •
[Documentation](https://docs.sourced.tech/community-edition) •
[Blog](https://blog.sourced.tech) •
[Slack](http://bit.ly/src-d-community) •
[Twitter](https://twitter.com/sourcedtech)


### Contents

- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Commands](#commands)
  - [Working With Multiple Data Sets](#working-with-multiple-data-sets)
- [Contributing](#contributing)
- [Community](#community)
- [Code of Conduct](#code-of-conduct)
- [License](#license)

## Quick Start

**source{d} CE** supports Linux, macOS, and Windows.

You will find in the [Quick Start Guide](docs/quickstart/README.md) all the steps to get started with **source{d} CE**, from the installation of its dependencies to running SQL queries to inspect git repositories.


## Usage

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
sourced init orgs --token <token> src-d
```

After a while you may want to analyze the data on another set of repositories. You can run `init` again with a different organization:
```
sourced init orgs --token <token> bblfsh
```

This command will stop any of the currently running containers, create an isolated environment for the new data, and create a new, clean deployment.

Please note that each path will have an isolated deployment. This means that for example any chart or dashboard created for the deployment for `src-d` will not be available to the new deployment for `bblfsh`.

Each isolated environment is persistent (unless you run `prune`). Which means that if you decide to re-deploy **source{d} CE** using the original organization:
```
sourced init orgs --token <token> src-d
```

You will get back to the previous state, and things like charts and dashboards will be restored.

These isolated environments also allow you to deploy **source{d} CE** using a local set of Git repositories. For example, if we wanted a third deployment to analyze repositories already existing in the `~/repos` directory, we just need to run `init` again:

```
sourced init local ~/repos
```

If you are familiar with Docker Compose and you want more control over the underlying resources, you can explore the contents of your `~/.sourced` directory. There you will find a `docker-compose.yml` and `.env` files for each set of repositories used by `sourced init`.


## Contributing

[Contributions](https://github.com/src-d/sourced-ce/issues) are **welcome and very much appreciated** 🙌
Please refer to [our Contribution Guide](docs/CONTRIBUTING.md) for more details.


## Community

source{d} has an amazing community of developers and contributors who are interested in Code As Data and/or Machine Learning on Code. Please join us! 👋

- [Community](https://sourced.tech/community/)
- [Slack](http://bit.ly/src-d-community)
- [Twitter](https://twitter.com/sourcedtech)
- [Email](mailto:hello@sourced.tech)


## Code of Conduct

All activities under source{d} projects are governed by the
[source{d} code of conduct](https://github.com/src-d/guide/blob/master/.github/CODE_OF_CONDUCT.md).


## License

GPL v3.0, see [LICENSE](LICENSE.md).
