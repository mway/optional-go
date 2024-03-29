#!/usr/bin/env just --justfile

main_pkg := ""
coverprofile := "cover.out"
dist_dir := env_var_or_default("DIST_DIR", "./dist")

default:
    @just --list | grep -v default

clean:
    rm -rf dist

###
# Testing

lint PKG="./...":
    golangci-lint run --new=false {{ PKG }}

test PKG="./..." *ARGS="":
    go test -race -failfast -count 1 -coverprofile {{ coverprofile }} {{ PKG }} {{ ARGS }}

vtest PKG="./..." *ARGS="": (test PKG ARGS "-v")

convert-coverage:
    go tool cover -html {{ coverprofile }}

cover: test && convert-coverage

mockgen VERSION="latest":
    command mockgen >/dev/null 2>&1 || go install go.uber.org/mock/mockgen@{{ VERSION }}

generate PKG="./...": mockgen
    go generate {{ PKG }}

###
# Benchmarks

alias benchmark := bench

bench PKG="./..." *ARGS="":
    go test -v -count 1 -run x -bench . {{ PKG }} {{ ARGS }} {{ ARGS }}

###
# Releasing

release VERSION: && _fetch_tags
    gh release create --title "{{ VERSION }}" --notes="" "{{ VERSION }}" && git fetch -t

release-assets VERSION: create-release-binaries && _fetch_tags
    gh release create --title "{{ VERSION }}" --notes="" "{{ VERSION }}" {{ dist_dir }}/*

###
# Binaries

run *ARGS="": _only_if_main_pkg
    go run {{ main_pkg }} {{ ARGS }}

install *ARGS="": _only_if_main_pkg
    go install {{ main_pkg }} {{ ARGS }}

create-release-binaries:
    #!/usr/bin/env bash
    set -eo pipefail

    # n.b. Windows is not currently supported.
    for os in darwin freebsd linux; do
        for arch in amd64 arm64; do
            GOARCH=$arch GOOS=$os bash -c "set -x; go build -o {{ dist_dir }}/git-flow-$arch-$os {{ main_pkg }}"
        done
    done

    for arm in 5 6 7; do
        GOARCH=arm GOOS=linux bash -c "set -x; go build -o {{ dist_dir }}/git-flow-armv$arm-linux {{ main_pkg }}"
    done

@_only_if_main_pkg:
    echo {{ if main_pkg == "" { error("not configured for executable programs") } else { "OK" } }} >/dev/null

@_fetch_tags:
    git fetch -pt
