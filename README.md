# Redash API Client #
[![Build Status][travis-image]][travis] [![Go Report Card][goreport-image]][goreport] [![License][license-image]][license]

## Overview ##

A Simple API client library for interacting with Redash written in Go. 

## Quick start ##

### Using go modules (aka. `go mod`) ###

In your go files, simply use:
```go
import "github.com/snowplow-devops/redash-client-go"
```

Then next `go mod tidy` or `go test` invocation will automatically
populate your `go.mod` with the last redash-client-go release, now
[![Version](https://img.shields.io/github/tag/snowplow-devops/redash-client-go.svg)](https://github.com/snowplow-devops/redash-client-go/releases).

Note you can use `go mod vendor` to vendor your dependencies.

## Usage ##

Functional examples can be found in
* https://github.com/snowplow-devops/redash-client-go/tree/master/examples 

## Development ##

Assuming git installed:

```bash
$ git clone https://github.com/snowplow-devops/redash-client-go
$ cd redash-client-go
$ make test
```

To remove all build files:

```bash
$ make clean
```

To format the golang code in the source directory:

```bash
$ make format
```

**Note:** Always run `format` before submitting any code.

**Note:** The `make test` command also generates a code coverage file which can be found at `build/coverage/coverage.html`.

### Copyright and license

The Redash Go Client is copyright 2019-2020 Snowplow Analytics Ltd.

Licensed under the **[Apache License, Version 2.0][license]** (the "License");
you may not use this software except in compliance with the License.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

[travis-image]: https://travis-ci.com/snowplow-devops/redash-client-go.png?branch=master
[travis]: https://travis-ci.com/snowplow-devops/redash-client-go

[release-image]: http://img.shields.io/badge/release-0.1.0-6ad7e5.svg?style=flat
[releases]: https://github.com/snowplow-devops/redash-client-go/releases

[license-image]: http://img.shields.io/badge/license-Apache--2-blue.svg?style=flat
[license]: http://www.apache.org/licenses/LICENSE-2.0

[goreport-image]: https://goreportcard.com/badge/github.com/snowplow-devops/redash-client-go
[goreport]: https://goreportcard.com/report/github.com/snowplow-devops/redash-client-go
