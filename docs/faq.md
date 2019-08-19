# FAQ

_For tips and advices to deal with unexpected errors, please refer to [Troubleshooting guide](./troubleshooting.md)_

## Index

- [Where Can I Find More Assistance to Run source{d} or Notify You About Any Issue or Suggestion?](#where-can-i-find-more-assistance-to-run-source-d-or-notify-you-about-any-issue-or-suggestion)
- [How Can I Update My Version Of **source{d} CE**?](#how-can-i-update-my-version-of-source-d-ce)
- [How To Restore Dashboards and Charts to Defaults](#how-to-restore-dashboards-and-charts-to-defaults)
- [How to Update the Data from the Organizations That I'm Analyzing](#how-to-update-the-data-from-the-organizations-being-analyzed)
- [Can I Query Gitbase or Babelfish with External Tools?](#can-i-query-gitbase-or-babelfish-with-external-tools)
- [Where Can I Read More About the Web Interface?](#where-can-i-read-more-about-the-web-interface)
- [I Get IOError Permission denied](#i-get-ioerror-permission-denied)

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

1. download the new binary for your OS
1. move the binary to the right place ([see Quickstart. Install](quickstart/2-install-sourced.md))
1. run `sourced compose download`
1. run `sourced restart` to apply the new configuration.

This process will reinstall **source{d} CE** with the new components, but it will
keep your data (repositories, charts, dashboards, etc). If you want to replace
all the current dashboards with the ones from the release that you just installed,
you have two alternatives:

- run `sourced prune --all` before running `sourced init`; but if you were using
repositories downloaded from a GitHub organization, they will be deleted, and
downloaded again.
- drop only the dashboards, and load the new ones, following also this other
instructions: [How To Restore Dashboards and Charts to Defaults](#how-to-restore-dashboards-and-charts-to-defaults)


## How To Restore Dashboards and Charts to Defaults

In some circumstances, you may want to restore the UI dashboard and charts to
its defaults. Currently, there is no clean way of doing it without deleting all
the state of the UI.

Following these steps, all the state modified from the UI (charts, dashboards,
saved queries, users, roles, etcetera) will be replaced by the default ones for
the version of **source{d} CE** that you're currently using. If you're using
repositories from a GitHub organization, all its data will be preserved, and only
charts and dashboards will be restarted.

To do so, you only need to delete the docker volume containing the PostgreSQL
database, and restart **source{d} CE**. It can be done following these steps if
you already have [Docker Compose](https://docs.docker.com/compose/) installed:

```shell
$ cd ~/.sourced/workdirs/__active__
$ source .env
$ docker-compose stop postgres
$ docker-compose rm -f postgres
$ ENV_PREFIX=`awk '{print tolower($0)}' <<< ${COMPOSE_PROJECT_NAME}`
$ docker volume rm ${ENV_PREFIX}_postgres
$ sourced restart
```


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
