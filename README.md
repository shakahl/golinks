# golinks

[![Build Status](https://cloud.drone.io/api/badges/prologic/msgbus/status.svg)](https://cloud.drone.io/prologic/msgbus)
[![CodeCov](https://codecov.io/gh/prologic/msgbus/branch/master/graph/badge.svg)](https://codecov.io/gh/prologic/msgbus)
[![Go Report Card](https://goreportcard.com/badge/prologic/msgbus)](https://goreportcard.com/report/prologic/msgbus)
[![GoDoc](https://godoc.org/github.com/prologic/msgbus?status.svg)](https://godoc.org/github.com/prologic/msgbus) 
[![Sourcegraph](https://sourcegraph.com/github.com/prologic/msgbus/-/badge.svg)](https://sourcegraph.com/github.com/prologic/msgbus?badge)

golinks is a web app that allows you to create smart bookmarks, commands and aliases by pointing your web browser's default search engine at a running instance. Similar to bunny1 or yubnub.

## Installation

### Source

```#!bash
$ go get github.com/prologic/golinks
```

### OS X Homebrew

There is a formula provided that you can tap and install from
[prologic/homebrew-golinks](https://github.com/prologic/homebrew-golinks):

```#!bash
$ brew tap prologic/golinks
$ brew install golinks
```

**NB:** This installs the latest released binary; so if you want a more
recent unreleased version from master you'll have to clone the repository
and build yourself.

## Usage

Run golinks:

```#!bash
$ golinks -bind 127.0.0.1:8000
```

Set your browser's default golinks engine to http://localhost:8000/?q=%s

Then type `help` to view the main help page, `g foo bar` to perform a [Google](https://google.com) search for "foo bar" or `list` to list all available commands.

## Contributing

golinks is considered "production" software and is used daily. If you find this interresting or useful please fork and contribute back via pull-requests! If you find bugs or have ideas for new features, please file an issue!

## License

MIT
