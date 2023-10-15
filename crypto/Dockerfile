FROM golang:1.20
WORKDIR /app
ENV CGO_ENABLED 1
RUN go install github.com/cosmtrek/air@latest
RUN git config --global --add safe.directory /app
COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml", "test", "./..."]

COPY . .
