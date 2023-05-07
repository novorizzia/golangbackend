# script ini akan di run oleh !/bin/sh karna pada alpine bash shell tidak tersedia

#  memastikan script akan exit immidiately if command return non zero status
set -e

echo "run db migration"
# $DB_SOURCE source to get it's value
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"
# $@ = takes all parameters pass to the script and run it
exec "$@"