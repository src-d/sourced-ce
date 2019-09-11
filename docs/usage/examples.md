# SQL Examples to Analyze Your Data

_If you want to know what the database schema looks like, you can refer to the [diagram about gitbase entities and relations](https://docs.sourced.tech/gitbase/using-gitbase/schema#database-diagram), or just use regular `SHOW` or `DESCRIBE` queries._

_In gitbase repository, you will find more [SQL examples of queries](https://docs.sourced.tech/gitbase/using-gitbase/examples)._


## Index

* [Queries For Repositories](#queries-for-repositories)
* [Queries With Files](#queries-with-files)
* [Queries With UASTs](#queries-with-uasts)
* [Queries About Comitters](#queries-about-comitters)


## Queries For Repositories

**Show me the repositories I am analyzing:**

```sql
SELECT * FROM repositories;
```

**Last commit messages for HEAD for every repository**

```sql
SELECT commit_message
FROM refs
NATURAL JOIN commits
WHERE ref_name = 'HEAD';
```

**Top 10 repositories by commit count from [HEAD](https://git-scm.com/book/en/v2/Git-Internals-Git-References#ref_the_ref):**

```sql
SELECT repository_id,commit_count
FROM (
    SELECT
        repository_id,
        COUNT(*) AS commit_count
    FROM ref_commits
    WHERE ref_name = 'HEAD'
    GROUP BY repository_id
) AS q
ORDER BY commit_count DESC
LIMIT 10;
```

**10 top repos by file count at HEAD**

```sql
SELECT repository_id, num_files FROM (
    SELECT COUNT(*) num_files, repository_id
    FROM refs
    NATURAL JOIN commit_files
    WHERE ref_name = 'HEAD'
    GROUP BY repository_id
) AS t
ORDER BY num_files DESC
LIMIT 10;
```


## Queries With Files

**Query for all LICENSE & README files across history:**

```sql
SELECT file_path, repository_id, blob_size
FROM files
WHERE
    file_path = 'LICENSE'
    OR file_path = 'README.md';
```

**Query all files at [HEAD](https://git-scm.com/book/en/v2/Git-Internals-Git-References#ref_the_ref):**

```sql
SELECT cf.file_path, f.blob_size
FROM ref_commits rc
NATURAL JOIN commit_files cf
NATURAL JOIN files f
WHERE
    rc.ref_name = 'HEAD'
    AND rc.history_index = 0;
```


## Queries With UASTs

_**Note**: UAST values are returned as binary blobs; they're best visualized in the web UI interface rather than the CLI where are seen as binary data._

**Retrieve the UAST for all files at [HEAD](https://git-scm.com/book/en/v2/Git-Internals-Git-References#ref_the_ref):**

```sql
SELECT * FROM (
    SELECT cf.file_path,
        UAST(f.blob_content, LANGUAGE(f.file_path,  f.blob_content)) as uast
    FROM ref_commits r
    NATURAL JOIN commit_files cf
    NATURAL JOIN files f
    WHERE
        r.ref_name = 'HEAD'
        AND r.history_index = 0
) t WHERE uast != '';
```


## Queries About Comitters

**Top committers per repository**

```sql
SELECT * FROM (
    SELECT
        commit_author_email as author,
        repository_id as id,
        count(*) as num_commits
    FROM commits
    GROUP BY commit_author_email, repository_id
) AS t
ORDER BY num_commits DESC;
```
