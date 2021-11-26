#!/bin/bash

container_name=test-postgres
docker_stat="$(docker container inspect -f '{{.State.Status}}' $container_name)"

# main directory
TEST_DIR_ROOT=/tmp/test_docker
# script creating db with data
INIT_DB_SCRIPT=init-db-data.sh

make_test_dir() {
   echo "Creating temp directory for docker container"
   rm -rf ${TEST_DIR_ROOT}
   mkdir ${TEST_DIR_ROOT}
   mkdir ${TEST_DIR_ROOT}/data
   mkdir ${TEST_DIR_ROOT}/sql
   cp ./db/*.sql ${TEST_DIR_ROOT}/sql/
   cp ./tools/docker-entrypoint-initdb.d/${INIT_DB_SCRIPT} ${TEST_DIR_ROOT}/${INIT_DB_SCRIPT}
}

docker container rm --force $container_name
make_test_dir
docker run -d \
      --name $container_name \
      -p 5433:5432 \
      -e POSTGRES_PASSWORD=test \
      -e POSTGRES_USER=user \
      -e POSTGRES_DB=test \
      -v ${TEST_DIR_ROOT}/data:/var/lib/postgresql/data \
      -v ${TEST_DIR_ROOT}/sql:/tmp/data \
      -v ${TEST_DIR_ROOT}/${INIT_DB_SCRIPT}:/docker-entrypoint-initdb.d/${INIT_DB_SCRIPT} \
      postgres   