# Mirage
[![Mirage](https://github.com/OdyseeTeam/mirage/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/OdyseeTeam/mirage/actions/workflows/go.yml)
[![Latest release](https://badgen.net/github/release/OdyseeTeam/mirage?cache=600)](https://github.com/OdyseeTeam/mirage/releases)
![Docker Image Version (latest semver)](https://img.shields.io/docker/v/odyseeteam/mirage)

Mirage is a software that helps Odysee distribute static content on its platform. The following features are covered:
1) images retrieval from any origin specified
2) images optimization needed for web distribution
3) images caching via [OdyseeTeam/gody-cdn](https://github.com/OdyseeTeam/gody-cdn)

## Installation
to be updated

## Usage
To be updated

## Building from Source
This project requires [Go v1.17](https://golang.org/doc/install).

On Ubuntu you can install it with `sudo snap install go --classic`.

if you want a specific version of go, you can always do `sudo snap refresh go --channel=1.17/stable` for example.

```
git clone git@github.com:OdyseeTeam/mirage.git
cd mirage
make
```

You may choose different targets:
- make test: run go tests
- make lint: run linters
- make linux: build linux binary
- make macos: build mac os binary
- make image: build docker image
- make publish_image: push docker image to docker hub
- make retag: move previous tag to current commit

## Contributing
Feel free to open a pull request or an issue anytime you like!

## License
This project is MIT licensed.

## Security
We take security seriously. Please contact security@odysee.com regarding any security issues.

## Contact
The primary contact for this project is [@Nikooo777](https://github.com/Nikooo777) (niko-at-odysee.com)