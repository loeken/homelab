docker run -u "$(id -u):$(id -g)" -v "$(pwd)":/keys authelia/authelia:latest authelia crypto pair rsa generate --bits 4096 --directory /keys
