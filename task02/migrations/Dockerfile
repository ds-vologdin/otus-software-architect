FROM migrate/migrate:v4.13.0

ADD wait-for.sh /bin/
RUN chmod +x /bin/wait-for.sh
ADD *.sql /migrations/

ENTRYPOINT [ "/bin/wait-for.sh" ]
