FROM alpine:3.12

RUN apk add --no-cache postgresql-client
ADD wait-for.sh /usr/bin/
RUN chmod a+x /usr/bin/wait-for.sh
ADD init_db.sql /sql/

ENTRYPOINT ["/usr/bin/wait-for.sh"]
