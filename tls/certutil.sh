#!/usr/bin/env sh

docker run --rm -it \
  -v "$(dirname "$(readlink -f "$0")")":/usr/share/elasticsearch/tls \
  docker.elastic.co/elasticsearch/elasticsearch:8.1.0 \
  bin/elasticsearch-certutil "$@"
