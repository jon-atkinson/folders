#!/bin/bash

if [ -z "$1" ]; then
  echo "use: ./bench.sh write_filename"
  exit 1
fi

go test -bench . -benchmem -count=10 ./... > "./benchmark/$1"