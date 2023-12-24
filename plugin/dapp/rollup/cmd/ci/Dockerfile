FROM ubuntu:20.04

WORKDIR /root
COPY chain chain
COPY chain-cli chain-cli
COPY *.toml ./

CMD /root/chain -f /root/$ChainConfig
