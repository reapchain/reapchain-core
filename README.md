# Reapchain Core

![banner](docs/reapchain_logo.png)

[![version](https://img.shields.io/github/tag/reapchain/reapchain-core.svg)](https://github.com/reapchain/reapchain-core/releases/latest)
[![Go version](https://img.shields.io/badge/go-1.18+-blue.svg)](https://github.com/moovweb/gvm)
[![Discord chat](https://img.shields.io/discord/669268347736686612.svg)](https://discord.gg/AzefAFd)

[Reapchain Core](https://reapchain.com/file/kr/ReapChain_WhitePaper_0.9_kr.pdf) is forked from Tendermint Core [v0.34.9](https://github.com/tendermint/tendermint/tree/v0.34.9) on 2021-07-16.
And we synced up with Tendermint-[v0.34.20](https://github.com/tendermint/tendermint/tree/v0.34.20) on 2023-04-07.

## Releases

Please do not depend on `main` as your production branch. Use [releases](https://github.com/reapchain/reapchain-core/releases) instead.

## Minimum requirements

| Requirement | Notes            |
| ----------- |------------------|
| Go version  | Go1.18 or higher |


# Quick Start
## git clone
```
git clone https://github.com/reapchain/reapchain-core.git
```
## Build
```
make build
make install
```
## Check version
```
podc version
```

## Local Standalone
```
cd localnet
# ./launch_local_testnet.sh [NUMBER_OF_STANDING] [NUMBER_OF_STEERING_CANDIDATE] [DIRECTORY]
./launch_local_testnet.sh 2 1 local-testnet
```
### Check Process
```
ps -ef | grep podc
```
### Check Log
```
tail -f ./local-testnet/logs/standing0_log.log
```

## Tools

Benchmarking is provided by [`tm-load-test`](https://github.com/informalsystems/tm-load-test).
Additional tooling can be found in [/docs/tools](/docs/tools).

## Applications

- [Cosmos SDK](http://github.com/reapchain/cosmos-sdk); a cryptocurrency application framework
- [Ethermint](https://github.com/reapchain/ethermint); Ethereum on Reapchain Core

## Research

- [The latest gossip on BFT consensus](https://arxiv.org/abs/1807.04938)
- [Original Reapchain Whitepaper](https://reapchain.com/file/kr/ReapChain_WhitePaper_0.9_kr.pdf)
