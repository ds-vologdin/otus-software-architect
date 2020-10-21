FROM golang:1.15-alpine3.12 AS build

ADD . /app/
WORKDIR /app/
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM scratch

COPY --from=build /app/app /bin/app
EXPOSE 8000
ENTRYPOINT ["/bin/app"]
