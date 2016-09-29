FROM alpine:3.3

RUN apk --update add sqlite

CMD echo "CREATE TABLE IF NOT EXISTS foo(bar int);" | sqlite3 /data/database.db
