FROM golang:1.20
WORKDIR /usr/src/app
RUN go install github.com/cosmtrek/air@latest
COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

COPY . .
