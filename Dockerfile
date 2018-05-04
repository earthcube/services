# Start from scratch image and add in a precompiled binary
# Cross compile if needed...
# CGO_ENABLED=0  go build .
# CGO_ENABLED=0 go build .
# docker build -t earthcube/p418services:latest -t earthcube/p418services:0.0.1 .
# docker run -d -p 6789:6789  earthcube/p418services:0.0.1
FROM scratch

# Add in the static elements (could also mount these from local filesystem)
# later as the indexes grow
ADD services /
ADD log /log
ADD indexcatalog.json /
# ADD logs /logs  or mount in the volume in the compose or docker command

# NOTE make sure to mount the bleve indexes to /index

# Add our binary
CMD ["/services"]

# Document that the service listens on this port
EXPOSE 6789
