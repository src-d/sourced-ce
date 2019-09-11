# Explore source{d} CE Web Interface

_If you have any problem running **source{d} CE** you can take a look to our [Troubleshooting](../learn-more/troubleshooting.md) section, and to our [source{d} Forum](https://forum.sourced.tech), where you can also ask for help when using **source{d} CE**. If you spotted a bug, or have a feature request, please [open an issue](https://github.com/src-d/sourced-ce/issues) to let us know abut it._

_In some circumstances, loading the data for the dashboards can take some time, and the UI can be frozen in the meanwhile. It can happen &mdash;on big datasets&mdash;, the first time you access the dashboards, or when they are refreshed. Please, take a look to our
[Troubleshooting](../learn-more/troubleshooting.md#the-dashboard-takes-a-long-to-load-and-the-ui-freezes)
to get more info about this exact issue._

Once **source{d} CE** has been [initialized with `sourced init`](./3-init-sourced.md), it will automatically open the web UI. If the UI is not automatically opened, you can use `sourced web` command, or visit http://127.0.0.1:8088.

Use login: `admin` and password: `admin`, to access the web interface.

If you [initialized **source{d} CE** from GitHub Organizations](./3-init-sourced.md#from-github-oganizations), its repositories and metadata will be downloaded on background, and it will be available graduatelly. You will find more info in the welcome dashboard once you log in.


## Sections

The most relevant features that **source{d} CE** Web Interface offers are:
- **[SQL Lab](#sql-lab-querying-code-and-metadata)**, to query your repositories and its GitHub metadata.
- **[Babelfish web](#uast-parsing-code)**, web interface to parse files into UAST.
- **[Dashboards](#dashboards)**, to aggregate charts for exploring and visualizing your data.
- **Charts**, to see your data with a rich set of data visualizations.
- A flexible UI to manage users, data sources, export data...

The user interface is based in the open-sourced [Apache Superset](http://superset.incubator.apache.org), so you can also refer to their documentation for advanced usage of the web interface.


## SQL Lab. Querying Code and Metadata

_If you prefer to work within the terminal via command line, you can open a SQL REPL running `sourced sql`_

Using the `SQL Lab` tab, from the web interface, you can analyze your dataset using SQL queries, and create charts from those queries with the `Explore` button.

You can find some sample queries in the [examples](../usage/examples.md).

If you want to know what the database schema looks like you can use either regular `SHOW` or `DESCRIBE` queries, or you can refer to the [diagram about gitbase entities and relations](https://docs.sourced.tech/gitbase/using-gitbase/schema#database-diagram).

```bash
$ sourced sql "SHOW tables;"
+--------------+
|    TABLE     |
+--------------+
| blobs        |
| commit_blobs |
| commit_files |
| commit_trees |
| commits      |
| files        |
| ref_commits  |
| refs         |
| remotes      |
| repositories |
| tree_entries |
+--------------+
```

```bash
$ sourced sql "DESCRIBE TABLE commits;"
+---------------------+-----------+
|        NAME         |   TYPE    |
+---------------------+-----------+
| repository_id       | TEXT      |
| commit_hash         | TEXT      |
| commit_author_name  | TEXT      |
| commit_author_email | TEXT      |
| commit_author_when  | TIMESTAMP |
| committer_name      | TEXT      |
| committer_email     | TEXT      |
| committer_when      | TIMESTAMP |
| commit_message      | TEXT      |
| tree_hash           | TEXT      |
| commit_parents      | JSON      |
+---------------------+-----------+
```


## UAST. Parsing code

_Please, refer to the [quick explanation about what Babelfish is](../usage/bblfsh.md) to know more about it._

You can get UASTs from the `UAST` tab (parsing files by direct input), or using the `UAST` gitbase function over blob contents on `SQL Lab` tab.


## Dashboards

_Please, refer to [Superset Tutorial, creating your first dashboard](http://superset.incubator.apache.org/tutorial.html) for more details._

The dashboards let you aggregate custom charts to show in the same place different metrics for your repositories.

You can create them:
- From the `Dashboard` tab, adding a new one, and then selecting new charts.
- From any chart view, the `Save` button will let you to add it into a new or existent one. 
