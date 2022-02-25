#!/bin/sh
set -e

scriptdir=$(cd $(dirname $0); pwd)
cd $scriptdir/../../../

case $1 in
  build)
    docker build --build-arg USER_ID=$(id -u) --build-arg GROUP_ID=$(id -g) --rm -f $scriptdir/Dockerfile -t telemetry-all-in-one:latest .
    ;;
  clean)
    docker image prune -f
    ;;
  *)
    echo "specify 'build' or 'clean' as the first arg"
    exit 1
    ;;
esac
