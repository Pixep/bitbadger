# BitBadger
![goreportcard](https://goreportcard.com/badge/github.com/Pixep/bitbadger)

BitBadger serves repository badges, and currently supports
* Open PR count (BitBucket Cloud)
* Open PR average duration (BitBucket Cloud)

Badges are generated from 

* Retrieve Open PR count and Open PR average duration from BitBucket Cloud
* Returns a custom badge from shields.io with the desired metric
* Badges are colored based on metric value

## Getting Started

Download and run `make` to build the binary. Go 1.12 is used, and can be installed from [https://golang.org/doc/install](https://golang.org/doc/install).

### Basic usage

Start the server, providing repository credentials. Only BitBucket Cloud is currently supported.
```
bitbadger [--port <port>] <username> <password>
```

Navigate or link the following URL for badges

`http://<server>:<port>/<username-or-group>/<repo-slug>/<metric-type>`, where metric can be `open-pr-count` or `avg-pr-time`
