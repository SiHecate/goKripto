FROM golang:1.21.0
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
