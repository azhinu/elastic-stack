#!/bin/bash
set -eo pipefail
# Set your [remote_monitoring_user] password below:
# Change [http] to [https] is you setted kibana to use TLS.
health=$(curl -fsSL --insecure 'http://remote_monitoring_user:<passwd>@localhost:5601/api/status' | /usr/libexec/platform-python -c "import sys, json; print(json.load(sys.stdin)['status']['overall']['level'])")
if [ "$health" = 'green' ] || [ "$health" = "yellow" ]; then
        exit 0
fi
echo >&2 "unexpected health status: $health"
exit 1
