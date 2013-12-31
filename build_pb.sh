#!/bin/bash
protoc --plugin=${GOPATH}/bin/protoc-gen-go message.proto  --go_out=.
