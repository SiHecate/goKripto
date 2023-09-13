FROM golang:1.20
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go install github.com/cosmtrek/air@latest
RUN go mod download

CMD ["air", "-c", ".air.toml"]

COPY . .
