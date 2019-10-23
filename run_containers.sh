#!/bin/bash


docker-compose down -v
docker-compose up -d

do_more=true
while [ $do_more == true ]; do
  for service in $(docker-compose ps); do
    do_more=false
    if [ $service == "Exit" ]; then
      docker-compose up -d
      do_more=true
      break
    fi
  done
done

# Just to make sure
docker-compose up -d

echo "All containers are running"


