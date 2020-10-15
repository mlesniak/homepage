#!/usr/bin/env sh

# trap ctrl-c
trap ctrl_c INT

function ctrl_c() {
  # Jup, every python process...
  pid=$(ps ax|grep SimpleHTTPServer|sort -n|cut -d' ' -f1)
  kill -9 $pid >/dev/null 2>&1
}

go build
cd docs/ && python -m SimpleHTTPServer >/dev/null 2>&1 &
watch -n 1 aperol