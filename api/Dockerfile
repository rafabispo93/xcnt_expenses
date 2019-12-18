FROM golang:1.13 AS build

RUN curl -sL https://taskfile.dev/install.sh | sh
RUN go get github.com/99designs/gqlgen

COPY . /app
WORKDIR /app

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN task build

FROM scratch AS deploy
COPY --from=build /app/build/api .
EXPOSE 8080
ENTRYPOINT ["./api"]
