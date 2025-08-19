#!/usr/bin/env bash
set -euo pipefail

host=${1:-localhost}
port=${2:-6379}

for i in $(seq 1 30); do
  if redis-cli -h "$host" -p "$port" ping &>/dev/null; then
    echo "redis is up"
    exit 0
  fi
  sleep 1
done

echo "redis did not start in time" >&2
exit 1
