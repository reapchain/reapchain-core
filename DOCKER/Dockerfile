# stage 1 Generate ReapchainCore Binary
FROM golang:1.15-alpine as builder
RUN apk update && \
    apk upgrade && \
    apk --no-cache add make
COPY / /reapchain-core
WORKDIR /reapchain-core
RUN make build-linux

# stage 2
FROM golang:1.15-alpine
LABEL maintainer="hello@reapchain-core.com"

# ReapchainCore will be looking for the genesis file in /reapchain-core/config/genesis.json
# (unless you change `genesis_file` in config.toml). You can put your config.toml and
# private validator file into /reapchain-core/config.
#
# The /reapchain-core/data dir is used by reapchain-core to store state.
ENV TMHOME /reapchain-core

# OS environment setup
# Set user right away for determinism, create directory for persistence and give our user ownership
# jq and curl used for extracting `pub_key` from private validator while
# deploying reapchain-core with Kubernetes. It is nice to have bash so the users
# could execute bash commands.
RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq bash && \
    addgroup tmuser && \
    adduser -S -G tmuser tmuser -h "$TMHOME"

# Run the container with tmuser by default. (UID=100, GID=1000)
USER tmuser

WORKDIR $TMHOME

# p2p, rpc and prometheus port
EXPOSE 26656 26657 26660

STOPSIGNAL SIGTERM

COPY --from=builder /reapchain-core/build/reapchain-core /usr/bin/reapchain-core

# You can overwrite these before the first run to influence
# config.json and genesis.json. Additionally, you can override
# CMD to add parameters to `podc node`.
ENV PROXY_APP=kvstore MONIKER=dockernode CHAIN_ID=dockerchain

COPY ./DOCKER/docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["node"]

# Expose the data directory as a volume since there's mutable state in there
VOLUME [ "$TMHOME" ]

