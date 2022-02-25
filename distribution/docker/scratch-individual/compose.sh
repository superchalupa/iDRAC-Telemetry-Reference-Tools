#!/bin/sh

scriptdir=$(cd $(dirname $0); pwd)
topdir=$(cd $scriptdir/../../..; pwd)
cd $topdir

set -e

# relies on already-built docker images. use 'build-docker.sh' to create if needed.
echo ARGS: $@
echo 1: $1

echo "USER_ID=$(id -u)" > $topdir/.env
echo "GROUP_ID=$(id -g)" >> $topdir/.env

case $1 in
  config)
    docker-compose --project-dir $topdir -f $scriptdir/docker-compose.yml config
    ;;

  stop)
    docker-compose --project-dir $topdir -f $scriptdir/docker-compose.yml rm -f
    ;;

  build-and-start)
    docker-compose --project-dir $topdir -f $scriptdir/docker-compose.yml up --build --remove-orphans
    ;;

  start)
    docker-compose --project-dir $topdir -f $scriptdir/docker-compose.yml up
    ;;

  *)
    echo "Specify 'start' or 'stop' as the first argument"
    exit 1
    ;;
esac
