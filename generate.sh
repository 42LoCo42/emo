#!/usr/bin/env bash
dir="api"
mkdir -p "$dir"
oapi-codegen \
	-package "$dir" \
	-generate "types,server,client" \
	<(yq < "api.yaml" | ref-merge) \
	> "$dir/api.go"
