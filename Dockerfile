FROM golang:alpine AS build

COPY . /server/

WORKDIR /server/

RUN go build app/cmd/main.go

FROM alpine:latest

COPY --from=build /server/main .

EXPOSE 8080

ENV POSTGRES_HOST=localhost
ENV POSTGRES_PORT=5432
ENV POSTGRES_USERNAME=db_pg
ENV POSTGRES_PASSWORD=db_postgres
ENV POSTGRES_DATABASE=db_tenders
ENV SERVER_ADDRESS=":8080"

ENTRYPOINT ["./main"]

