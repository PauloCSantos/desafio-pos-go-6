FROM golang:1.20

WORKDIR /workspaces

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
