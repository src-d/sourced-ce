# Frequently Asked Questions

_For tips and advices to deal with unexpected errors, please refer to [Troubleshooting guide](./troubleshooting.md)_

## Index

- [Where Can I Find More Assistance to Run source{d} or Notify You About Any Issue or Suggestion?](#where-can-i-find-more-assistance-to-run-source-d-or-notify-you-about-any-issue-or-suggestion)
- [How Can I Update My Version Of **source{d} CE**?](#how-can-i-update-my-version-of-source-d-ce)
- [How to Update the Data from the Organizations That I'm Analyzing](#how-to-update-the-data-from-the-organizations-being-analyzed)
- [Can I Query Gitbase or Babelfish with External Tools?](#can-i-query-gitbase-or-babelfish-with-external-tools)
- [Where Can I Read More About the Web Interface?](#where-can-i-read-more-about-the-web-interface)
- [I Get IOError Permission denied](#i-get-ioerror-permission-denied)
- [Why Do I Need Internet Connection?](#why-do-i-need-internet-connection)


## Where Can I Find More Assistance to Run source{d} or Notify You About Any Issue or Suggestion?

_If you're dealing with an error or something that you think that can be caused
by an unexpected error, please refer to our [Troubleshooting guide](./troubleshooting.md).
With the info that you can obtain following those steps, you could fix the problem
or you will be able to explain it better in the following channels:_

* [open an issue](https://github.com/src-d/sourced-ce/issues), if you want to
suggest a new feature, if you need assistance with a contribution, or if you
found any bug.
* [Visit the source{d} Forum](https://forum.sourced.tech) where users and community
members discuss anything source{d} related. You will find there some common questions
from other source{d} users, or ask yours.
* [join our community on Slack](https://sourced-community.slack.com/join/shared_invite/enQtMjc4Njk5MzEyNzM2LTFjNzY4NjEwZGEwMzRiNTM4MzRlMzQ4MmIzZjkwZmZlM2NjODUxZmJjNDI1OTcxNDAyMmZlNmFjODZlNTg0YWM),
and talk with some of our engineers.


## How Can I Update My Version Of source{d} CE?

When there is a new release of **source{d} CE**, it is noticed every time a `sourced`
command is called. When it happens you can download the new version from
[src-d/sourced-ce/releases/latest](https://github.com/src-d/sourced-ce/releases/latest),
and proceed as it follows:

(You can also follow these steps if you want to update to any beta version, to
downgrade, or to use your own built version of **source{d} CE**)

1. replace your current version of `sourced` from its current location with the
one you're installing ([see Quickstart. Install](quickstart/2-install-sourced.md)),
and confirm it was done by running `sourced version`.
1. run `sourced compose download` to download the new configuration.
1. run `sourced restart` to apply the new configuration.

This process will reinstall **source{d} CE** with the new components, but it will
keep your current data (repositories, metadata, charts, dashboards, etc) of your
existent workdirs.

If you want to replace all your current customizations &mdash;including charts and
dashboards&mdash;, with the ones from the release that you just installed, the
official way to proceed is to `prune` the running workdirs, and `init` them again.

_**disclaimer:** pruning a workdir will delete all its data: its saved queries
and charts, and if you were using repositories and metadata downloaded from a
GitHub organization, they will be deleted, and downloaded again._

1. `sourced status workdirs` to get the list of your current workdirs
1. Prune the workdirs you need, or prune all of them at once running
`sourced prune --all`
1. `sourced init [local|orgs] ...` for each workdir again, to initialize them with
the new configuration.


## How to Update the Data from the Organizations Being Analyzed

There is no way to update imported data, and
[when a scraper is restarted](./troubleshooting.md#how-can-i-restart-one-scraper),
it procedes as it follows:

### gitcollector

Organizations and repositories are downloaded independently, so if they fail,
the process is not stopped until all the organizations and repositories have been
iterated.

If `gitcollector` is restarted, it will download more repositories, but it won’t
update any of the already existent ones. You can see the progress of the new process
in the welcome dashboard; since already existent repositories won't be updated,
those will appear as `failed` in progress status.

### ghsync

The way how metadata is imported by `ghsync` is a bit different, and it is done
sequentially per each organization, so if any step fails, the whole importation
will fail.

Pull requests, issues, and users of the same organization, are imported in that
order in separate transaction each one, and if one transaction fails, the process
will be stopped so the next ones won't be processed.

Once the three different entities have been imported, the organization will be
considered as "done", and restarting `ghsync` won't cause to update its data.

If `ghsync` is restarted, it will only import data from organizations that could
not be finished considering the rules explained above. The process of `ghsync`
will be updated in the welcome dashboard and if an organization was already
imported, it will appear as "nothing imported" in the status chart.


## Can I Query Gitbase or Babelfish with External Tools?

Yes, as explained in our docs about [**source{d} CE** Architecture](./architecture.md#docker-networking),
these and other components are exposed to the host machine, to be used by third
party tools like [Jupyter Notebook](https://jupyter.org/),
[gitbase clients](https://docs.sourced.tech/gitbase/using-gitbase/supported-clients)
and [Babelfish clients](https://docs.sourced.tech/babelfish/using-babelfish/clients).

The connection values that you should use to connect to these components, are
defined in the [`docker-compose.yml`](../docker-compose.yml), and sumarized in
the [Architecture documentation](./architecture.md#docker-networking)


## Where Can I Read More About the Web Interface?

The user interface is based in the open-sourced [Apache Superset](http://superset.apache.org),
so you can also refer to [Superset tutorials](http://superset.apache.org/tutorial.html)
for advanced usage of the web interface.


## I Get IOError Permission denied

If you get this error message:

```
IOError: [Errno 13] Permission denied: u'./.env'
```

This may happen if you have installed Docker from a snap package. This installation mode is not supported, please install it following [the official documentation](./quickstart/1-install-requirements.md#install-docker) (See [#78](https://github.com/src-d/sourced-ce/issues/78)).


## Why Do I Need Internet Connection?

source{d} CE automatically fetches some resources from the Internet when they are not found locally:

- the source{d} CE configuration is fetched automatically when initializing it for the first time, using the proper version for the current version of `sourced`, e.g. if using `v0.16.0` it will automatically fetch `https://raw.githubusercontent.com/src-d/sourced-ce/v0.16.0/docker-compose.yml`.
- to download the docker images of the source{d} CE components when initializing source{d} for the first time, or when initializing it after changing its configuration.
- to download repositories and its metadata from GitHub when you initialize source{d} CE with `sourced init orgs`.
- to download and install [Docker Compose alternative](#docker-compose) if there is no local installation of Docker Compose.

If your connection to the network does not let source{d} CE to access to Internet, you should manually provide all these dependencies.
