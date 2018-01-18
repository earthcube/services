# Start from scratch image and add in a precompiled binary
# CGO_ENABLED=0 env GOOS=linux go build .
# docker build -t earthcube/p418services:latest -t earthcube/p418services:0.0.1 .
# docker run -d -p 6789:6789  earthcube/p418services:0.0.1
FROM scratch

# Add in the static elements (could also mount these from local filesystem)
# later as the indexes grow
ADD services /

# Add our binary
CMD ["/services"]

# Document that the service listens on this port
EXPOSE 6789