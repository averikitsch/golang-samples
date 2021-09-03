FROM golang
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get install -y \
  apt-utils \
  dumb-init \
  nfs-common \
  && apt-get clean

# Create all mount points
RUN mkdir -p /mnt/nfs/filestore

COPY run.sh /run.sh
COPY main.go /go/src/app/main.go
COPY go.mod /go/src/app/
WORKDIR /go/src/app
RUN go build -v -o server
WORKDIR /
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/run.sh"]