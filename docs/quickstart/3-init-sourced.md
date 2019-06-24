# Initialize source{d} Community Edition

_For the full list of the sub-commands offered by `sourced`, please take a look
into [the `sourced` sub-commands inventory](../usage/commands.md)._

**source{d} CE** can be initialized from 2 different data sources: GitHub organizations, or local Git repositories.

Please note that you have to choose one data source to initialize **source{d} CE**, but you can have more than one isolated environment, and they can have different sources. 

**source{d} CE** will download and install Docker images on demand. Therefore, the first time you run some of these commands, they might take a bit of time to start up. Subsequent runs will be faster.


#### From GitHub Organizations

When using GitHub organizations to populate the **source{d} CE** database you only need to provide a list of organization names and a [GitHub personal access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/). If no scope is granted to the user token, only public
data will be fetched. To let **source{d} CE** to access to also private repos and hidden users, the token should
have the following scopes enabled:

- `repo` Full control of private repositories
- `read:org` Read org and team membership, read org projects


Use this command to initialize, e.g.

```shell
$ sourced orgs init --token <token> src-d,bblfsh
```

It will also download, on background, the repositories of the passed GitHub organizations, and its metadata: pull requests, issues, users...


#### From Local Repositories

```shell
$ sourced init </path/to/repositories>
```

It will initialize **source{d} CE** to analyze the git repositories under the passed path, or under the current directory if no one is passed. The repositories will be found recursively.

**Note for macOS:**
Docker for Mac [requires enabling file sharing](https://docs.docker.com/docker-for-mac/troubleshoot/#volume-mounting-requires-file-sharing-for-any-project-directories-outside-of-users) for any path outside of `/Users`.

**Note for Windows:** Docker for Windows [requires shared drives](https://docs.docker.com/docker-for-windows/#shared-drives). Other than that, it's important to use a working directory that doesn't include any sub-directory whose access is not readable by the user running `sourced`. For example, using `C:\Users` as workdir, will most probably not work. For more details see [this issue](https://github.com/src-d/engine/issues/250).
