#!/usr/bin/env bash
docker run --rm -it \
  -v "$(dirname "$(readlink -f "$0")")":/usr/share/elasticsearch/tls \
  docker.elastic.co/elasticsearch/elasticsearch:7.16.1 \
  bin/elasticsearch-certutil "$@"
