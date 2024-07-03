#!/bin/bash

go mod download

go build -o bitvector ./cmd/bitvector/main.go