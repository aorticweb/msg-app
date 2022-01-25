#!/bin/bash
# Perform postgres migrations
echo -e "\n=== Trying to connect to pg database ==="
dbmate --url $POSTGRES_URL wait
echo -e "\n=== DB is responsive, beginning migrations ==="
dbmate --url $POSTGRES_URL --migrations-dir /db/migrations up
wait
echo -e "\n=== PG migrations complete ==="