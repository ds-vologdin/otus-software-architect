FROM golang:1.15-alpine3.12 AS build

WORKDIR /app/
ENV GO111MODULE=on

ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . /app/
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/auth/

FROM scratch

COPY --from=build /app/app /bin/app
EXPOSE 8001
ENTRYPOINT ["/bin/app"]
