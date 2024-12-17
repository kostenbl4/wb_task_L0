#!/bin/sh

DBSTRING="host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=disable"

/goose -dir /migrations postgres "$DBSTRING" up