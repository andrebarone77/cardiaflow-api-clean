#!/bin/bash

go build -ldflags="-s -w" -o build/cardiaflow ./cmd/api