FROM golang:1.20 # Build an image strating with golang version 1.20
WORKDIR /app # working directory /app 
RUN go install github.com/cosmtrek/air@latest # Air for hotstart 
RUN git config --global --add safe.directory /app # git config for air
COPY go.mod go.sum ./ # requirements 
RUN go mod download 

CMD ["air", "-c", ".air.toml", "test", "./..."]

COPY . .
