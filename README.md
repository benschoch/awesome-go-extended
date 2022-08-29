# Extended HTML for [Awesome Go](https://github.com/avelino/awesome-go)

This repo provides a slightly enhanced HTML version of the awesome [Awesome Go](https://github.com/avelino/awesome-go)
repository.

All credits for the content go the project and its countless contributors. Thanks for providing such a great overview of awesome packages for Go – [Share some love if you like it](https://github.com/sponsors/avelino) 🦄

## TL;DR

You can download and open the [index.html](./index.html) file or head your browser to the github.io page.

## Build HTML

1. Get a [GitHub access token](https://github.com/settings/tokens)
2. Build the HTML by running this command with your access token:
    ```shell
    GITHUB_ACCESS_TOKEN=... go run main.go
    ```
   _You can run it without a token, but you will probably reach the rate limit._

3. Open [index.html](./index.html) file in your favorite browser

## Technical details

The generation of the HTML works as follows:

1. The `README.md` of [Awesome Go](https://github.com/avelino/awesome-go) is downloaded and parsed into categories and
   packages.
2. Details about GitHub repositories are being loaded and added.
3. The data is rendered as HTML, with the help of [html/template](https://pkg.go.dev/html/template).
   The template itself mainly uses [Bootstrap](https://getbootstrap.com/), [jQuery](https://jquery.com/)
   and [DataTables](https://datatables.net/)

To limit repetitive requests to GitHub, the responses are cache in a local file cached
with [onecache](https://github.com/adelowo/onecache).
The cache can be flushed by passing the `-flush=true` flag to the run command:

```shell
GITHUB_ACCESS_TOKEN=... go run main.go -flush=true
```

## Disclaimer

This tool highly relies on the current format of the [Awesome Go](https://github.com/avelino/awesome-go) `README.md`.
As long as that doesn't change too much, the parsing should be able to handle it.