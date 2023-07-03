#!/usr/bin/env bash
go run github.com/ogen-go/ogen/cmd/ogen \
	-clean \
	-skip-unimplemented \
	api.yaml
