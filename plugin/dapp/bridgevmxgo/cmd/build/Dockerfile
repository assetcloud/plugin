FROM ubuntu:16.04

WORKDIR /root
COPY chain chain
COPY chain-cli chain-cli
COPY chain.toml chain*.toml ./
COPY entrypoint.sh entrypoint.sh

CMD ["/root/chain", "-f", "/root/chain.toml"]
