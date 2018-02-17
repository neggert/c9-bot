#! /bin/sh

source $1
shift

eval $(docker-machine env $DOCKER_MACHINE)

docker-compose $@