#!/bin/bash
set -eo pipefail

health=$(curl -fsSL --insecure 'http://localhost:9600' | python3 -c "import sys, json; print(json.load(sys.stdin)['status'])")
if [ "$health" = 'green' ] || [ "$health" = "yellow" ]; then
        exit 0
fi
echo >&2 "unexpected health status: $health"
exit 1
