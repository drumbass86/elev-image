#!/bin/bash

container_name=dev-postgres
docker_stat="$(docker container inspect -f '{{.State.Status}}' $container_name)"

if [[ $docker_stat != "running" ]]; then
   if [[ $docker_stat != "paused" || $docker_stat == "" ]]; then   
      echo "container $container_name donot paused or running. Kill and started new container!"
      docker container rm $container_name
      docker run -d \
      --name $container_name \
      -p 5432:5432 \
      -e POSTGRES_PASSWORD=123qwe \
      -e POSTGRES_USER=user \
      -v ${HOME}/postgres/postgresql-data:/var/lib/postgresql/data \
      postgres
   else
      docker container start $container_name
   fi
else
   echo "container $container_name already running!"
fi