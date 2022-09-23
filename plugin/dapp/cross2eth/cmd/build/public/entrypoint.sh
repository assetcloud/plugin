#!/usr/bin/env bash
/root/chain -f /root/chain.toml &
# to wait nginx start
sleep 15
/root/chain -f "$PARAFILE"
