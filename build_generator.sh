#!/bin/bash

go mod download

go build -o generator ./cmd/bitvector_generator/main.go