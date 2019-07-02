# FAQ and Troubleshooting

## Index

- [Where Can I Find More Assistance to Run source{d} or Notify You About Any Issue or Suggestion?](#where-can-i-find-more-assistance-to-run-source-d-or-notify-you-about-any-issue-or-suggestion)
- [How Can I Update My Version Of source{d} CE?](#how-can-i-update-my-version-of-source-d-ce)
- [How To Restore Dashboards and Charts to Defaults](#how-to-restore-dashboards-and-charts-to-defaults)
- [How Can I See Logs of Running Components?](#how-can-i-see-logs-of-running-components)
- [Can I Query Gitbase or Babelfish with External Tools?](#can-i-query-gitbase-or-babelfish-with-external-tools)
- [Where Can I Read More About the Web Interface?](#where-can-i-read-more-about-the-web-interface)
- [When I Try to Create a Chart from a Query, Nothing Happens.](#when-i-try-to-create-a-chart-from-a-query-nothing-happens)
- [When I Try to Export a Dashboard, Nothing Happens.](#when-i-try-to-export-a-dashboard-nothing-happens)
- [The Dashboard Takes a Long to Load, and the UI Freezes.](#the-dashboard-takes-a-long-to-load-and-the-ui-freezes)


## Where Can I Find More Assistance to Run source{d} or Notify You About Any Issue or Suggestion?

If this documentation was not enough, you could also try:

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


## How Can I See Logs of Running Components?

```shell
$ cd ~/.sourced/workdirs/__active__
$ docker-compose logs -f [components...]
```

Where `-f` will keep the connection opened, and the logs will appear as they come
instead of exiting after the last logged one.

Where you can pass a space separated list of component names to see only their
logs (i.e. `sourced-ui`, `gitbase`, `bblfsh`, `gitcollector`, `ghsync`, `metadatadb`, `postgres`, `redis`).
If you do not pass any component name, there will appear the logs of all of them.


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


## When I Try to Create a Chart from a Query, Nothing Happens.

The charts can be created from the SQL Lab, using the `Explore` button once you
run a query. If nothing happens, the browser may be blocking the new window that
should be opened to edit the new chart. You should configure your browser to let
source{d} UI to open pop-ups (e.g. in Chrome it is done allowing `127.0.0.1:8088`
to handle `pop-ups and redirects` from the `Site Settings` menu).


## When I Try to Export a Dashboard, Nothing Happens.

If nothing happens when pressing `Export` button from the dashboard list, then
you should configure your browser to let source{d} UI to open pop-ups (e.g. in
Chrome it is done allowing `127.0.0.1:8088` to handle `pop-ups and redirects`
from the `Site Settings` menu)


## The Dashboard Takes a Long to Load and the UI Freezes.

_This is a known issue that we're trying to address, but here is more info about it._

In some circumstances, loading the data for the dashboards can take some time,
and the UI can be frozen in the meanwhile. It can happen &mdash;on big datasets&mdash;,
the first time you access the dashboards, or when they are refreshed.

There are some limitations with how Apache Superset handles long-running SQL
queries, which may affect the dashboard charts. Since most of the charts of the
Overview dashboard loads its data from gitbase, its queries can take more time
than the expected for the UI.

When it happens, the UI can be frozen, or you can get this message in some charts:
>_Query timeout - visualization queries are set to timeout at 300 seconds.
Perhaps your data has grown, your database is under unusual load, or you are
simply querying a data source that is too large to be processed within the timeout
range. If that is the case, we recommend that you summarize your data further._

When it occurs, you should wait till the UI is responsive again, and separately
refresh each failing chart with its `force refresh` option (on its top-right corner).
With some big datasets, it took 3 refreshes and 15 minutes to get data for all charts.
