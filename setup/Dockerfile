FROM golang:alpine

RUN apk update
RUN mkdir /app
WORKDIR /app
COPY . /app
RUN go build -ldflags="-s -w" -o setup .
CMD ["go", "run", "/app/setup.go"]