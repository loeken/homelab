FROM golang:alpine

RUN apk update
RUN apk add git terraform github-cli curl

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

RUN curl -L https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64 -o /usr/bin/cloudflared
RUN mkdir /app
RUN chmod +x /usr/bin/cloudflared

WORKDIR /app
COPY . /app
RUN go build -ldflags="-s -w" -o setup .
CMD ["go", "run", "/app/setup.go"]