#!/usr/bin/env bash

bridge_name=br1

for node in {0..2}; do
	name="node$node"
	docker rm -f $name > /dev/null 2>&1
	rm -f /run/netns/$name
done

ip link del $bridge_name
