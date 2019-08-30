# sourced Command Reference

`sourced` binary offers you different kinds of sub-commands:

* [to manage their containers](commands.md#manage-containers)
* [to manage **source{d} CE** configuration](commands.md#manage-configuration)
* [to open interfaces to access its data](commands.md#open-interfaces)
* [show info about the command](commands.md#others)

Here is the list of all these commands and its description; you can get more info about each one adding `--help` when you run it.

## Manage Containers

### sourced init

_There is a dedicated section to document this command in the quickstart about_ [_how to initialize **source{d} CE**_](../quickstart/3-init-sourced.md)

This command installs and initializes **source{d} CE** docker containers, networks, and volumes, downloading its docker images if needed.

It can work over a local repository or a list of GitHub organizations.

**source{d} CE** will download and install Docker images on demand. Therefore, the first time you run some of these commands, they might take a bit of time to start up. Subsequent runs will be faster.

Once **source{d} CE** has been initialized, it will automatically open the web UI. If the UI is not opened automatically, you can use [`sourced web`](commands.md#sourced-web) command, or visit [http://127.0.0.1:8088](http://127.0.0.1:8088).

Use login: `admin` and password: `admin`, to access the web interface.

#### sourced init orgs

```text
$ sourced init orgs --token=_USER_TOKEN_ org1,org2...
```

Installs and initializes **source{d} CE** for a list of GitHub organizations, downloading their repositories and metadata: Users, PullRequests, Issues...

The `orgs` argument must be a comma-separated list of GitHub organizations.

The `--token` must contain a valid GitHub user token for the given organizations. It should be granted with 'repo' and'read:org' scopes.

#### sourced init local

```text
$ sourced init local [/path/to/repos]
```

Installs and initializes **source{d} CE** using a local directory containing the git repositories to be processed by **source{d} CE**. If the local path to the `workdir` is not provided, the current working directory will be used.

### sourced start

Starts all the components that were initialized with `init` and then stopped with `stop`.

### sourced stop

Stops all running containers without removing them. They can be started again with `start`.

### sourced prune

Stops containers and removes containers, networks, volumes, and configurations created by `init` for the current working directory.

To delete resources for all the installed working directories, add the `--all` flag.

Container images are not deleted unless you specify the `--images` flag.

If you want to completely uninstall `sourced` you must also delete the `~/.sourced` directory.

### sourced logs

Show logs from source{d} components.

If `--follow` is used the logs are shown as they are logged until you exit with `Ctrl+C`.

You can optionally pass component names to see only their logs.

```text
$ sourced logs
$ sourced logs --follow
$ sourced logs --follow gitbase bblfsh
```

## Manage Configuration

### sourced status

Shows the status of **source{d} CE** components.

### sourced workdirs

Lists previously initialized working directories.

### sourced compose

Manages Docker Compose files in the `~/.sourced` directory with the following subcommands:

### sourced compose download

Download the `docker-compose.yml` file to define **source{d} CE** services. By default, the command downloads the file for this binary version, but you can also download other version or any other custom one using its URL.

Examples:

```text
$ sourced compose download
$ sourced compose download v0.0.1
$ sourced compose download master
$ sourced compose download https://raw.githubusercontent.com/src-d/sourced-ce/master/docker-compose.yml
```

### sourced compose list

Lists the available `docker-compose.yml` files, and shows which one is active. You can activate any other with `compose set`.

### sourced compose set

Sets the active `docker-compose.yml` file.

#### sourced restart

Updates current installation according to the active docker compose file.

It only recreates the component containers, keeping all your data, as charts, dashboards, repositories and GitHub metadata.

## Open Interfaces

### sourced sql

Opens a MySQL client connected to gitbase.

You can also pass a SQL query to be run by gitbase instead of opening the REPL, e.g.

```text
$ sourced sql "show databases"

+----------+
| Database |
+----------+
| gitbase  |
+----------+
```

**source{d} CE** SQL supports a [UAST](commands.md#babelfish-uast) function that returns a Universal AST for the selected source text. UAST values are returned as binary blobs and are best visualized in the [SQL Lab, from the web interface](https://github.com/src-d/sourced-ce/tree/89a607e16c0163b71eb79d83bad0e507aaffa168/docs/usage/.../quickstart/4-explore-sourced.md#sql-lab-querying-code) rather than the CLI where are seen as binary data.

### sourced web

Opens the web interface in your browser.

Use login: `admin` and password: `admin`, to access the web interface.

## Others

### sourced version

Shows the version of the `sourced` command being used.

### sourced completion

Prints a bash completion script for sourced; you can place its output in `/etc/bash_completion.d/sourced`, or add it to your `.bashrc` running:

```text
$ echo "source <(sourced completion)" >> ~/.bashrc
```

