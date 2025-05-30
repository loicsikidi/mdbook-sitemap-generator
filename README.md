# mdbook-sitemap-generator

> [!NOTE]
> This project is a fork of the original [mdbook-sitemap-generator](https://github.com/rxdn/mdbook-sitemap-generator) which is no longer maintained.
> Contrary to the original, this version is written in Golang.

## What is this?

mdbook-sitemap-generator is a simple utility to generate sitemap.xml files for mdbook projects.

## Installation

Binaries are distributed on the [Github Releases Page](https://github.com/loicsikidi/mdbook-sitemap-generator/releases).

It is also possible to install this utility via go, using `go install github.com/loicsikidi/mdbook-sitemap-generator@latest`.

### Verify binaries integrity using cosign

To verify the integrity of the downloaded binaries, you can use the provided checksums in [Github Releases Page](https://github.com/loicsikidi/mdbook-sitemap-generator/releases).

In order to verify the integrity and trustworthiness of the later, you can use the `cosign` tool, which is a part of the [sigstore project](https://sigstore.dev/).

Run the following commands to do so:

```bash
version=1.0.1 # replace with the version you want to verify

curl -sSL "https://github.com/loicsikidi/mdbook-sitemap-generator/releases/download/v${version}/mdbook-sitemap-generator_${version}_checksums.txt-keyless.pem" -o keyless.pem
curl -sSL "https://github.com/loicsikidi/mdbook-sitemap-generator/releases/download/v${version}/mdbook-sitemap-generator_${version}_checksums.txt-keyless.sig" -o keyless.sig
curl -sSL "https://github.com/loicsikidi/mdbook-sitemap-generator/releases/download/v${version}/mdbook-sitemap-generator_${version}_checksums.txt" -o checksums.txt

cosign verify-blob --certificate keyless.pem --signature keyless.sig checksums.txt \
--certificate-identity "https://github.com/loicsikidi/mdbook-sitemap-generator/.github/workflows/release.yaml@refs/tags/v${version}" \
--certificate-oidc-issuer https://token.actions.githubusercontent.com
# output: Verified OK
```

Once the verification is successful, you can use the `checksums.txt` file to verify the integrity of each artefact using `sha256sum` command.

## Usage

This project can be used as a mdbook plugin, or as a standalone utility.

### Used as a mdbook plugin

This utility can be used as a plugin for mdbook, which will automatically generate a sitemap.xml file during the build process (i.e., `mdbook build`).

#### Configuration

To use this utility as a plugin, you need to add the following configuration to your `book.toml` file:

```toml
[output.sitemap-generator]
# Domain of the site, used to generate absolute URLs
# This field is REQUIRED
domain = "docs.example.com"
```

> [!NOTE]
> `sitemap.xml` file which will be stored in `./book/sitemap-generator` directory.

### Used as a standalone utility

The utility should be run from the root of the project.

```
USAGE:
    mdbook-sitemap-generator --domain <DOMAIN> [OPTIONS]

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
$ mdboo@k-sitemap-generator --domain docs.example.com --output book/sitemap.xml
```

> [!TIP]
> The utility will automatically detect the book's root directory by parsing `book.toml` and fallback to `src` if it finds nothing.
