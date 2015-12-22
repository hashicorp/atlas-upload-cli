Atlas Upload CLI
================
[![Latest Version](http://img.shields.io/github/release/hashicorp/atlas-upload-cli.svg?style=flat-square)][release]
[![Build Status](http://img.shields.io/travis/hashicorp/atlas-upload-cli.svg?style=flat-square)][travis]
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[release]: https://github.com/hashicorp/atlas-upload-cli/releases
[travis]: http://travis-ci.org/hashicorp/atlas-upload-cli
[godocs]: http://godoc.org/github.com/hashicorp/atlas-upload-cli

The Atlas Upload CLI is a lightweight command line interface for uploading
application code to [Atlas][] to kick off deployment processes. This is the CLI
used to power the `vagrant push` command and other parts of [Atlas Go][] with
the Atlas strategy.

It can also be downloaded and used externally with other systems (such as a CI
service like Jenkins or Travis CI) to initiate Atlas-based deploys.

Usage
-----

```
atlas-upload [options] slug path

  Upload application code or artifacts to Atlas for initiating deployments.

  "slug" is the name of the <username>/<application_name> to upload to within Atlas.

  If path is a directory, it will be compressed (gzip tar) and uploaded
  in its entirety. The root of the archive will be the path. For clarity:
  if you upload the "foo/" directory, then the file "foo/version" will be
  "version" in the archive since "foo/" is the root.

  A path must be specified. Due to the nature of this application, it does
  not default to using the current working directory automatically.

Options:

  -exclude=<path>     Glob pattern of files or directories to exclude (this may
                      be specified multiple times)
  -include=<path>     Glob pattern of files/directories to include (this may be
                      specified multiple times, any excludes will override
                      conflicting includes)
  -address=<url>      The address of the Atlas server
  -token=<token>      The Atlas API token
  -vcs                Get lists of files to exclude and include from version
                      control system (Git, Mercurial or Subversion)
  -metadata<k=v>      Arbitrary key-value (string) metadata to be sent with the
                      upload; may be specified multiple times

  -debug              Turn on debug output
  -version            Print the version of this application
```

FAQ
---
**Q: Can I specify my Atlas access token via an environment variable?**<br>
A: All of HashiCorp's products support the `ATLAS_TOKEN` environment variable.
You can set this value in your shell profile or securely in your environment and
it will be used.


Contributing
------------
To hack on the Atlas Upload CLI, you will need a modern [Go][] environment. To
compile the `atlas-upload` binary and run the test suite, simply execute:

```shell
$ make
```

This will compile the `atlas-upload` binary into `bin/atlas-upload` and
run the test suite.

If you just want to run the tests:

```shell
$ make test
```

Or to run a specific test in the suite:

```shell
go test ./... -run SomeTestFunction_name
```

Submit Pull Requests and Issues to the [Atlas Upload CLI project on GitHub][Atlas Upload CLI].


[Atlas]: https://atlas.hashicorp.com "HashiCorp's Atlas"
[Atlas Go]: https://github.com/hashicorp/atlas-go "Atlas Go on GitHub"
[Atlas Upload CLI]: https://github.com/hashicorp/atlas-upload-cli "Atlas Upload CLI on GitHub"
[Go]: http://golang.org "Go the language"
