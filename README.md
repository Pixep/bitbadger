# BitBadger
[![Build Status](https://travis-ci.com/Pixep/bitbadger.svg?branch=master)](https://travis-ci.com/Pixep/bitbadger)
[![Go Report Card](https://goreportcard.com/badge/github.com/Pixep/bitbadger)](https://goreportcard.com/report/github.com/Pixep/bitbadger)

## Description

BitBadger is an HTTP/S server that creates linkable badges for BitBucket Cloud repositories. The core features are:
* Generates badges to be used in README.md Markdown documentations or anywhere
* Provide repository health metrics, as badges, to quickly identify areas to improve
    * See below for a list of supported badges/metrics
* Supports BitBucket Cloud repositories
* Runs as an HTTP or HTTPS server

## Supported badges

The following metrics can be used to monitor your repository:
* Average PR merge time (BitBucket Cloud)
    * `avg-pr-merge-time`
    * ![avg-pr-merge-time](doc/avg-pr-merge-time.svg)
* Open PR count (BitBucket Cloud)
    * `open-pr-count`
* Open PR average age (BitBucket Cloud)
    * `open-pr-avg-age`
* Oldest Open PR age (BitBucket Cloud)
    * `oldest-open-pr-age`

Badges are generated using:
* BitBucket Cloud API v2
    * See [BitBucket Cloud API v2 reference](https://developer.atlassian.com/bitbucket/api/2/reference/resource/repositories/%7Busername%7D/%7Brepo_slug%7D/pullrequests)
* shields.io badge generation service

## Getting Started

* Clone this repository
* Run `make` to build BitBadger. It uses Go 1.12, which can be easily installed from [https://golang.org/doc/install](https://golang.org/doc/install).
* Run it `./bitbadger ...`

### Basic usage

Due to BitBucket Cloud API, you will need to provide credentials to run the server. Only BitBucket Cloud is currently supported at the time.

Run the server using HTTP using `--insecure` flag.
```
bitbadger [--port <port>] --insecure <username> <password>
```

To run using HTTPS, you will need run it providing the path to your private key and certificate.
```
bitbadger [--port <port>] --cert <certificate-file> --key <private-key-file> <username> <password>
```

### Link badges

To link the badges, use the following URL:

`http://<server>:<port>/<username-or-group>/<repo-slug>/<badge-type>`

* `<username-or-group>`: Owning user or group, as visible in your repository URL
* `<repository-slug>`: Repository slug, as visible in your repository URL
* `<badge-type>`: One of
    * `open-pr-count`, `open-pr-avg-age`, `oldest-open-pr-age`, or `avg-pr-merge-time`

Markdown example:

`![avg-pr-merge-time](https://yourserver:34000/myuser/myrepository/avg-pr-merge-time.svg)`

![avg-pr-merge-time](doc/avg-pr-merge-time.svg)

## Advanced usage

### Caching

BitBadger supports caching requests to minimize traffice and latency. Note that caching is disabled by default. You can enable and adjust the caching behavior using the following options:

* `--cachevalidity`: Validity duration of the cache, in minutes. Defaults to `0`, which disables caching.
* `--maxcached`: Maximum number of cached requests. Defaults to `100`

### Command line options

```
GLOBAL OPTIONS:
   --debug, -d             Enable debug mode
   --insecure, -i          Enable insecure HTTP, without TLS
   --cert value, -c value  Path to TLS certificate
   --key value, -k value   Path to TLS private key
   --port value, -p value  Set the port that the server listens on (default: 34000)
   --cachevalidity value   Set for how long the requests should be cached in minutes (default: 0)
   --maxcached value       Set the maximum number of cached requests (default: 100)
   --help, -h              show help
   --version, -v           print the version
```

## Aknowledgment

Thanks to https://shields.io for their badge generation service.
