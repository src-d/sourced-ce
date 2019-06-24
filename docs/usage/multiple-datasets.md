# Working With Multiple Data Sets

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
