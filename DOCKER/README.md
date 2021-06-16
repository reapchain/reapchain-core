# Docker

## Supported tags and respective `Dockerfile` links

DockerHub tags for official releases are [here](https://hub.docker.com/r/reapchain/reapchain/tags/). The "latest" tag will always point to the highest version number.

Official releases can be found [here](https://github.com/reapchain/reapchain/releases).

The Dockerfile for reapchain is not expected to change in the near future. The master file used for all builds can be found [here](https://raw.githubusercontent.com/reapchain/reapchain/master/DOCKER/Dockerfile).

Respective versioned files can be found <https://raw.githubusercontent.com/reapchain/reapchain/vX.XX.XX/DOCKER/Dockerfile> (replace the Xs with the version number).

## Quick reference

- **Where to get help:** <https://reapchain.com/>
- **Where to file issues:** <https://github.com/reapchain/reapchain/issues>
- **Supported Docker versions:** [the latest release](https://github.com/moby/moby/releases) (down to 1.6 on a best-effort basis)

## Reapchain

Reapchain Core is Byzantine Fault Tolerant (BFT) middleware that takes a state transition machine, written in any programming language, and securely replicates it on many machines.

For more background, see the [the docs](https://docs.reapchain.com/master/introduction/#quick-start).

To get started developing applications, see the [application developers guide](https://docs.reapchain.com/master/introduction/quick-start.html).

## How to use this image

### Start one instance of the Reapchain core with the `kvstore` app

A quick example of a built-in app and Reapchain core in one container.

```sh
docker run -it --rm -v "/tmp:/reapchain" reapchain/reapchain init
docker run -it --rm -v "/tmp:/reapchain" reapchain/reapchain node --proxy_app=kvstore
```

## Local cluster

To run a 4-node network, see the `Makefile` in the root of [the repo](https://github.com/reapchain/reapchain/blob/master/Makefile) and run:

```sh
make build-linux
make build-docker-localnode
make localnet-start
```

Note that this will build and use a different image than the ones provided here.

## License

- Reapchain's license is [Apache 2.0](https://github.com/reapchain/reapchain/blob/master/LICENSE).

## Contributing

Contributions are most welcome! See the [contributing file](https://github.com/reapchain/reapchain/blob/master/CONTRIBUTING.md) for more information.
