# mdbook-sitemap-generator

> INFO: this is a fork of the original [mdbook-sitemap-generator](https://github.com/rxdn/mdbook-sitemap-generator) which is no longer maintained.
> Contrary to the original, this version is written in Golang.

## What is this?

mdbook-sitemap-generator is a simple utility to generate sitemap.xml files for mdbook projects.

## Installation

Binaries are distributed on the [Github Releases Page](https://github.com/loicsikidi/mdbook-sitemap-generator/releases).

It is also possible to install this utility via go, using `go install github.com/loicsikidi/mdbook-sitemap-generator@latest`.

## Usage

The utility should be run from the root of the project.

```
USAGE:
    mdbook-sitemap-generator [OPTIONS] --domain <DOMAIN>

OPTIONS:
    --domain <DOMAIN>
    --help               Print help information
    --version            Print version information
    --output <OUTPUT>
```

When running the utility, you must pass the site's domain on URL via the `--domain` flag, for example, `--domain docs.example.com`.

If the `--output` flag is not passed, the sitemap will be written to stdout.

For example:
```
$ ls
book  book.toml  src
$ mdbook-sitemap-generator --domain docs.example.com --output book/sitemap.xml
```

> TIP: the utility will automatically detect the book's root directory by parsing `book.toml` and fallback to `src` if it find nothing.
