#!/bin/sh
rm -r /tmp/versiondb || true
cd old
nix develop -c go run ./main.go /tmp/versiondb | wc -l
cd ../new
nix develop -c go run ./query/main.go /tmp/versiondb | wc -l
nix develop -c go run ./fix/main.go /tmp/versiondb | wc -l
nix develop -c go run ./query/main.go /tmp/versiondb | wc -l
