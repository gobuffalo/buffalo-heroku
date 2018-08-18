<p align="center"><img src="https://github.com/gobuffalo/buffalo/blob/master/logo.svg" width="360"></p>

<p align="center">
<a href="https://godoc.org/github.com/gobuffalo/buffalo-heroku"><img src="https://godoc.org/github.com/gobuffalo/buffalo-heroku?status.svg" alt="GoDoc" /></a>
<a href="https://goreportcard.com/report/github.com/gobuffalo/buffalo-heroku"><img src="https://goreportcard.com/badge/github.com/gobuffalo/buffalo-heroku" alt="Go Report Card" /></a>
</p>

# buffalo-heroku

This is a plugin for [https://gobuffalo.io](https://gobuffalo.io) that makes working with Heroku easier.

It assumes you are using Docker to deploy to Heroku. It is recommended you read [https://devcenter.heroku.com/articles/container-registry-and-runtime](https://devcenter.heroku.com/articles/container-registry-and-runtime) first.

## Installation

```bash
$ go get -u -v github.com/gobuffalo/buffalo-heroku
```

## Help

```bash
$ buffalo heroku --help
```

## Pre-Requisites

* You should absolutely have read [https://devcenter.heroku.com/articles/container-registry-and-runtime](https://devcenter.heroku.com/articles/container-registry-and-runtime) first.
* You should have the Heroku CLI installed [https://devcenter.heroku.com/articles/heroku-cli](https://devcenter.heroku.com/articles/heroku-cli).

## New

The `buffalo heroku new` command will setup and create a new Heroku app for you, with a bunch of defaults that **I** find nice.

### Flags/Options

There are a lot of flags and options you can use to tweak the Heroku app you create. Use the `--help` flag to see a list of them all.

```bash
$ buffalo heroku new -h
```

## Deploying

The initial `setup` command will do a deploy at the end, but after that you'll want to use the `buffalo heroku deploy` command to push a new version of your application, it'll even try to run your migrations for you.

```bash
$ buffalo heroku deploy
```

