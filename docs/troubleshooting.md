
# Troubleshooting:

_For commonly asked questions and their answers, you can refer to the [FAQ](./faq.md)_

Currently, **source{d} CE** does not expose nor log all errors directly into the
UI. In the current stage of **source{d} CE**, following these steps is the
better way to know if something is failing, why, and to know how to recover the
app from some problems. The first two steps use to be always mandatory:

1. **[To see if any component is broken](#how-can-i-see-the-status-of-source-d-ce-components)**
1. **[To see the logs of the running components](#how-can-i-see-logs-of-the-running-components)**
1. [To know if scrapers finished their job](#how-can-i-see-what-happened-with-the-scrapers)
  - [To restart one scraper](#how-can-i-restart-one-scraper)
1. [To restart o initialize **source{d} CE** again](#how-to-restart-source-d-ce)
1. [To ask for help if the issue could not be solved](./faq.md#where-can-i-find-more-assistance-to-run-source-d-or-notify-you-about-any-issue-or-suggestion)

Other issues that we detected, and which are strictly related to the UI are:

- [When I Try to Create a Chart from a Query, Nothing Happens.](#when-i-try-to-create-a-chart-from-a-query-nothing-happens)
- [When I Try to Export a Dashboard, Nothing Happens.](#when-i-try-to-export-a-dashboard-nothing-happens)
- [The Dashboard Takes a Long to Load, and the UI Freezes.](#the-dashboard-takes-a-long-to-load-and-the-ui-freezes)


## source{d} CE Fails During Its Initialization

The initialization can fail fast if there is any port conflict, or missing config
file, etcetera; those errors are clearly logged in the terminal when they appear.

If when initializing **source{d} CE**, all the required components appear as created,
but the loading spinner keeps spinning forever (more than 1 minute can be symptomatic),
there can be an underlying problem causing the UI not to be opened. In this
situation you should:

1. **[See if any component is broken](#how-can-i-see-the-status-of-source-d-ce-components)**
1. **[See app logs or certain component logs](#how-can-i-see-logs-of-the-running-components)**
1. [Restart o initialize **source{d} CE** again](#how-to-restart-source-d-ce)
1. [To ask for help if the issue could not be solved](./faq.md#where-can-i-find-more-assistance-to-run-source-d-or-notify-you-about-any-issue-or-suggestion)


## How Can I See the Status of source{d} CE Components?

To see the status of **source{d} CE** components, just run:

```
$ sourced status

Name                      Command                   State         Ports
------------------------------------------------------------------------------
srcd-xxx_sourced-ui_1    /entrypoint.sh             Up (healthy)  :8088->8088
srcd-xxx_gitbase_1       ./init.sh                  Up            :3306->3306
srcd-xxx_bblfsh_1        /tini -- bblfshd           Up            :9432->9432
srcd-xxx_bblfsh-web_1    /bin/bblfsh-web -addr ...  Up            :9999->8080
srcd-xxx_metadatadb_1    docker-entrypoint.sh  ...  Up            :5433->5432
srcd-xxx_postgres_1      docker-entrypoint.sh  ...  Up            :5432->5432
srcd-xxx_redis_1         docker-entrypoint.sh  ...  Up            :6379->6379
srcd-xxx_ghsync_1        /bin/sh -c sleep 10s  ...  Exit 0
srcd-xxx_gitcollector_1  /bin/dumb-init -- /bi ...  Exit 0

```

It will report the status of all **source{d} CE** component. All components should
be `Up`, but the scrapers: `ghsync` and `gitcollector`; these exceptions are
explanined in [How Can I See What Happened with the Scrapers?](#how-can-i-see-what-happened-with-the-scrapers)

If any component is not `Up` (but the scrapers), here are some key points to
understand what might be happening:

- All the components (but the scrapers) are restarted by Docker Compose
automatically &mdash;process that can take some seconds&mdash;; if the component
enters in a restart loop, something wrong is happening.
- When any component is failing, or died, you should
[see its logs to understand what is happening](#how-can-i-see-logs-of-the-running-components)

When one of the required components fails, it uses to print an error in the UI,

e.g. `lost connection to mysql server during query` while running a query might
mean that `gitbase` went down. 

e.g. `unable to establish a connection with the bblfsh server: deadline exceeded`
in SQL Lab might mean that `bblfsh` went down.

If the failing component is not successfully restarted in a few seconds, or if it
goes down when running certain queries, it could be a good idea to [open an issue](https://github.com/src-d/sourced-ce/issues)
describing the problem.


## How Can I See Logs of The Running Components?

```shell
$ sourced logs [-f] [components...]
```

Adding `-f` will keep the connection opened, and the logs will appear as they
come instead of exiting after the last logged one.

You can pass a space-separated list of component names to see only their logs
(i.e. `sourced-ui`, `gitbase`, `bblfsh`, `gitcollector`, `ghsync`, `metadatadb`, `postgres`, `redis`).
If you do not pass any component name, there will appear the logs of all of them.

Currently, there is no way to filter by error level, so you could try with `grep`,
e.g. 

```shell
sourced logs gitcollector | grep error
```

will output only log lines where `error` word appears.


## How Can I See What Happened with the Scrapers?

_When **souece{d} CE** is initialized with `sourced init local`, the scrapers are
not relevant because the repositories to analyze comes from your local data, so
`ghsync` and `gitcollector` status is not relevant in this case._

When running **souece{d} CE** to analyze data from a list of GitHub organizations,
`gitcollector` component is in charge of fetching  GitHub repositories and `ghsync`
component is in charge of fetching GitHub metadata (issues, pull requests...)

Once the UI is opened, you can see the progress of the importation in the welcome
dashboard, reporting the data imported, skipped, failed and completed. The process
can take many minutes if the organization is big, so be patient. You can manually
refresh both charts to confirm that the process is progressing, and it is not stuck.
If you believe that there can be any problem during the process, the better way
to find what is happening is:

- **[check the components status](#how-can-i-see-the-status-of-source-d-ce-components)
with `sourced status`**; `gitcollector` and `ghsync` should be `Up` (the process
didn't finish yet), or `Exit 0` (the process finished succesfully). They are
independent components, so they can finish on different order depending on how
many repositories or metadata is needed to process.

- **[check the logs](#how-can-i-see-logs-of-the-running-components) of the failing component with `sourced logs [-f] {gitcollector,ghsync}`**
to get more info about the errors found.


## How Can I Restart One Scraper?

_Restarting a scraper should be done to recover from temporal problems like
connectivity loss, or lack of space in disc, not
[to update the data you're analyzing](./faq.md#how-to-update-the-data-from-the-organizations-being-analyzed)_

**source{d} CE** does not provide way to start only one scraper. The recommended way
to restart them would be [to restart the whole **source{d} CE**](#how-to-restart-source-d-ce),
which is fast and safe for your data. In order to restart **source{d} CE**, run:

```shell
$ sourced restart
```

_Read more about [which data will be imported after restarting a scraper](./faq.md#how-to-update-the-data-from-the-organizations-being-Analyzed)_

If you feel comfortable enough with Docker Compose, you could also try restarting
each scraper separatelly, running:

```shell
$ cd ~/.sourced/workdirs/__active__
$ docker-compose run gitcollector # to restart gitcollector
$ docker-compose run ghsync       # to restart ghsync
```


## How to Restart source{d} CE

Restarting **source{d} CE**, can fix some errors and is also the official way to
restart the scrapers. It is also needed after downloading a new config (by running
`sourced compose download`). **source{d} CE** is restarted with the command:

```shell
$ sourced restart
```

It only recreates the component containers, keeping all your data, like charts,
dashboards, repositories, and GitHub metadata.


## When I Try to Create a Chart from a Query, Nothing Happens.

The charts can be created from the SQL Lab, using the `Explore` button once you
run a query. If nothing happens, the browser may be blocking the new window that
should be opened to edit the new chart. You should configure your browser to let
source{d} UI to open pop-ups (e.g. in Chrome it is done allowing `127.0.0.1:8088`
to handle `pop-ups and redirects` from the `Site Settings` menu).


## When I Try to Export a Dashboard, Nothing Happens.

If nothing happens when pressing the `Export` button from the dashboard list, then
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
