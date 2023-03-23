# We need to build in a Linux environment to support C libraries, e.g. RocksDB.
# We use Debian instead of Alpine, so that we can use binary database packages
# instead of spending time compiling them.
FROM golang:1.15

RUN apt-get -qq update -y && apt-get -qq upgrade -y >/dev/null
RUN apt-get -qq install -y libleveldb-dev librocksdb-dev >/dev/null

# Set up build directory /src/reapchain-core
ENV REAPCHAIN_BUILD_OPTIONS badgerdb,boltdb,cleveldb,rocksdb
WORKDIR /src/reapchain-core

# Fetch dependencies separately (for layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Build ReapchainCore and install into /usr/bin/reapchain-core
COPY . .
RUN make build && cp build/reapchain-core /usr/bin/reapchain-core
COPY test/e2e/docker/entrypoint* /usr/bin/
RUN cd test/e2e && make maverick && cp build/maverick /usr/bin/maverick
RUN cd test/e2e && make node && cp build/node /usr/bin/app

# Set up runtime directory. We don't use a separate runtime image since we need
# e.g. leveldb and rocksdb which are already installed in the build image.
WORKDIR /reapchain-core
VOLUME /reapchain-core
ENV TMHOME=/reapchain-core

EXPOSE 26656 26657 26660 6060
ENTRYPOINT ["/usr/bin/entrypoint"]
CMD ["node"]
STOPSIGNAL SIGTERM
