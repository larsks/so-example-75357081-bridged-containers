#!/usr/bin/env bash

set -e

bridge_name=br1
bridge_ip=10.10.10.1/24
ping_target=8.8.8.8

# ensure /run/netns exists
mkdir -p /run/netns

echo "create the bridge"
go run br_init.go -name $bridge_name -ip $bridge_ip

for node in {0..2}; do
	name="node$node"
	echo "create container $name"
	bash new_sandbox.sh $name 10.10.10.$((100 + $node))/24 br1 10.10.10.1

	echo "try to ping $ping_target from node0"
	if ! ip netns exec node0 ping -c2 $ping_target; then
		echo "ping from node0 failed!"
		exit 1
	fi

	echo "try to ping $ping_target from $name"
	if ! ip netns exec $name ping -c2 $ping_target; then
		echo "ping from $name failed!"
		exit 1
	fi

	cat <<-EOF
	
	---------------------------------------------------------------------------

	EOF
done

echo "all tests successful"
