#!/bin/bash
set -eo pipefail
# Set your [remote_monitoring_user] password below:
health=$(curl -fsSL --insecure 'https://remote_monitoring_user:<passwd>@localhost:9200/_cat/health?h=status')
if [ "$health" = "green" ] || [ "$health" = "yellow" ]; then
        exit 0
fi
echo >&2 "unexpected health status: $health"
exit 1
