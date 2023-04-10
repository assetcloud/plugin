FROM ubuntu:20.04

WORKDIR /root
COPY chain chain
COPY chain-cli chain-cli
COPY chain.toml ./

CMD ["/root/chain", "-f", "/root/chain.toml"]
