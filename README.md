# chronos [![build status](https://img.shields.io/travis/kataras/chronos/master.svg?style=flat-square)](https://travis-ci.org/kataras/chronos)[![report card](https://img.shields.io/badge/report%20card-a%2B-ff3333.svg?style=flat-square)](http://goreportcard.com/report/kataras/chronos)[![github issues](https://img.shields.io/github/issues/kataras/chronos.svg?style=flat-square)](https://github.com/kataras/chronos/issues?q=is%3Aopen+is%3Aissue)[![github closed issues](https://img.shields.io/github/issues-closed-raw/kataras/chronos.svg?style=flat-square)](https://github.com/kataras/chronos/issues?q=is%3Aissue+is%3Aclosed)[![release](https://img.shields.io/github/release/kataras/chronos.svg?style=flat-square)](https://github.com/kataras/chronos/releases)[![learn by examples](https://img.shields.io/badge/learn%20by-examples-0077b3.svg?style=flat-square)](https://github.com/kataras/chronos/tree/master/_examples)[![docs](https://img.shields.io/badge/chronos%20-godocs-0366d6.svg?style=flat-square)](https://godoc.org/github.com/kataras/chronos)

Chronos provides an easy way to limit X operations per Y time in accuracy of nanoseconds.
Chronos is a crucial tool when calling external, remote or even internal APIs with rate limits.

### 📑 Table Of Content

* [Installation](#-installation)
* [Latest changes](HISTORY.md)
* [Examples](#-learn)
    * [Simple native](_examples/native-acquire)
    * [Simple http client](_examples/http)
* [Community & Support](#-community)
* [Versioning](#-version)
* [People](#-people)

### 🚀 Installation

The only requirement is the [Go Programming Language](https://golang.org/dl/).

```bash
$ go get -u github.com/kataras/chronos
```

> `chronos` has zero third-party dependencies rathen than the Go's standard library.

`Chronos` is very simple to use tool, you just declare the maximum operations(max) and in every Y time
that those operations are allowed to be executed(per). `Chronos` has the accuracy of nanoseconds.

You can use `chronos#Acquire` anywhere you want, but be careful because it waits/blocks for
availability to be executed.

```go
import (
    "time"

    "github.com/kataras/chronos"
)

func main() {
    // Allow maximum of 5 MyActions calls
    // per 3 seconds.
    c := chronos.New(5, 3 * time.Second)

    // The below actions will execute immediately.
    MyAction(c)
    MyAction(c)
    MyAction(c)
    MyAction(c)
    MyAction(c)
    // The below action will execute after 3 seconds from the last call
    // or when 3 seconds passed without any c.Acquire call.
    MyAction(c)
}

func MyAction(c *chronos.C) {
    <- c.Acquire()
    // Do something here...
}
```

The `<- c.Acquire()` is blocking when needed, no any further actions neeeded by the end-developer.

### Helpers

The [ext](ext) package provides two subpackages; `function` and `http`. Navigate there to learn how they can help you.

### 📖 Learn

The awesome `chronos` community is always adding new content to learn from, [_examples](_examples/) is a great place to get started!

Read the [godocs](https://godoc.org/github.com/kataras/chronos) for a better understanding.

### 👥 Community

Join the welcoming community of fellow `chronos` developers in [rocket.chat](https://kataras.rocket.chat/channel/chronos)

- [Post](https://github.com/kataras/chronos/issues/new) a feature request or report a bug
- :star: and watch the public [repository](https://github.com/kataras/chronos/stargazers), will keep you up to date
- :earth_americas: publish [an article](https://medium.com/search?q=chronos) or share a [tweet](https://twitter.com/hashtag/golang) about your personal experience with `chronos`.

### 📌 Version

Current: [0.0.1](VERSION)

Read more about Semantic Versioning 2.0.0

 - http://semver.org/
 - https://en.wikipedia.org/wiki/Software_versioning
 - https://wiki.debian.org/UpstreamGuide#Releases_and_Versions

### 🥇 People

The original author of `chronos` is [@kataras](https://github.com/kataras), you can reach him via;

- [Medium](https://medium.com/@kataras)
- [Twitter](https://twitter.com/makismaropoulos)
- [Dev.to](https://dev.to/@kataras)
- [Facebook](https://facebook.com/kataras.gopher)
- [Mail](mailto:kataras2006@hotmail.com?subject=Chronos%20I%20need%20some%20help%20please)

[List of all Authors](AUTHORS)

[List of all Contributors](https://github.com/kataras/chronos/graphs/contributors)

## License

This software is licensed under the open source 3-Clause BSD License.

You can find the license file [here](LICENSE), for any questions regarding the license please [contact with me](mailto:mailto:kataras2006@hotmail.com?subject=Chronos%20License).