FROM ubuntu:16.04

WORKDIR /root
COPY chain chain
COPY chain-cli chain-cli
COPY entrypoint.sh entrypoint.sh
COPY chain.toml chain*.toml ./


CMD ["/root/chain", "-f", "/root/chain.toml"]
