#!/bin/sh

set -e

if [ "$1" = 'services' ]; then
    /services -host ${GEODEX_HOST}
else
    exec "$@"
fi
