# List of `sourced` Commands

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
