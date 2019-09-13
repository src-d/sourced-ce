# Babelfish UAST

_In the [Babelfish documentation](https://docs.sourced.tech/babelfish/), you will
find detailed information about Babelfish specifications, usage, examples, etc._

One of the most important components of **source{d} CE** is the UAST, which stands for:
[Universal Abstract Syntax Tree](https://docs.sourced.tech/babelfish/uast/uast-specification-v2).

UASTs are a normalized form of a programming language's AST, annotated with language-agnostic roles and transformed with language-agnostic concepts (e.g. Functions, Imports, etc.).

These enable an advanced static analysis of code and easy feature extraction for statistics or [Machine Learning on Code](https://github.com/src-d/awesome-machine-learning-on-source-code).


## UAST Usage

From the web interface, you can use the `UAST` tab, to parse files by direct input, or you can also get UASTs from the `SQL Lab` tab, using the `UAST(content)` [gitbase function](https://docs.sourced.tech/gitbase/using-gitbase/functions).

For the whole syntax about how to query the UASTs, you can refer to [How To Query UASTs With Babelfish](https://docs.sourced.tech/babelfish/using-babelfish/uast-querying)


## Supported Languages

To see which languages are available, check the table of [Babelfish supported languages](https://docs.sourced.tech/babelfish/languages).


## Clients and Connectors

The language parsing server (Babelfish) is available from the web interface, but you can also connect to the parsing server, deployed by **source{d} CE**, with several language clients, currently supported and maintained:

- [Babelfish Go Client](https://github.com/bblfsh/go-client)
- [Babelfish Python Client](https://github.com/bblfsh/client-python)
- [Babelfish Scala Client](https://github.com/bblfsh/client-scala)
