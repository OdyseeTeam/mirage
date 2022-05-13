#!/usr/bin/env bash
# example: ./scripts/release.sh patch
# available: major.minor.patch
go install github.com/caarlos0/svu@latest
git tag `svu "$1"`
git push --tags
