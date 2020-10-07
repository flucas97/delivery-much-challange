#!/usr/bin/env bash
set -e
if [ $# -eq 0 ]; then /app/api; fi;
exec $@

