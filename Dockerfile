FROM golang:1.24-alpine 

WORKDIR /app/server

COPY server/go.mod server/go.sum ./
RUN go mod download && go mod verify

COPY server .

RUN go build -o ./bin/hackernoon ./cmd

ENV GO_ENV=production

EXPOSE 3000

ENTRYPOINT ["/app/bin/hackernoon"]