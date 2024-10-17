#!/bin/bash

if [ -z "$1" ]; then
  echo "use: ./bench.sh write_filename"
  exit 1
fi

# Run the Go benchmarks and write to the specified or default file
go test -bench . -benchmem -count=10 ./... > "$1"