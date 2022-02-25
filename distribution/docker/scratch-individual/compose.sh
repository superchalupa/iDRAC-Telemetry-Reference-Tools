#!/bin/sh

scriptdir=$(cd $(dirname $0); pwd)
topdir=$(cd $scriptdir/../../..; pwd)
cd $topdir

set -e

# Can automatically rebuild docker images if needed
echo "USER_ID=$(id -u)" > $topdir/.env
echo "GROUP_ID=$(id -g)" >> $topdir/.env
PROFILE_ARG=

opts=$(getopt \
  -n $0 \
  --longoptions "build" \
  --longoptions "influx-pump" \
  --longoptions "prometheus-pump" \
  --longoptions "splunk-pump" \
  --longoptions "influx-test-db" \
  --longoptions "prometheus-test-db" \
  --longoptions "grafana" \
  -- "$@")
if [ $? -ne 0 ]; then
  echo "options not recognized"
  exit 1
fi
eval set -- "$opts"
while true; do
  echo "processing arg: $1"
  case "$1" in
    --build)
      BUILD_ARG="--build --remove-orphans"
      ;;
    --influx-pump)
      PROFILE_ARG="$PROFILE_ARG --profile influx-pump"
      ;;
    --prometheus-pump)
      PROFILE_ARG="$PROFILE_ARG --profile prometheus-pump"
      ;;
    --splunk-pump)
      PROFILE_ARG="$PROFILE_ARG --profile splunk-pump"
      ;;
    --influx-test-db)
      PROFILE_ARG="$PROFILE_ARG --profile influx-test-db"
      ;;
    --prometheus-test-db)
      PROFILE_ARG="$PROFILE_ARG --profile prometheus-test-db"
      ;;
    --grafana)
      PROFILE_ARG="$PROFILE_ARG --profile grafana"
      ;;
    --)
      shift
      break
      ;;
  esac
  shift
done


case $1 in
  stop)
    docker-compose --project-dir $topdir -f $scriptdir/docker-compose.yml rm -f
    ;;

  start)
    docker-compose --project-dir $topdir -f $scriptdir/docker-compose.yml ${PROFILE_ARG} up ${BUILD_ARG}
    ;;

  *)
    echo "Specify 'start' or 'stop'"
    exit 1
    ;;
esac
