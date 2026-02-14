FROM golang:1.26
WORKDIR /src
COPY . .
RUN go mod download

WORKDIR /src/cmd/action
RUN CGO_ENABLED=0 go build -o /usr/bin/action .

ENTRYPOINT ["/usr/bin/action"]