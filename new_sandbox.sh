#!/usr/bin/env bash
set -euo pipefail


container_name=$1
container_ip=$2
br_name=$3
gw_ip=$4


echo "creating container $container_name connected to br $br_name with ip $container_ip, gw $gw_ip"

docker rm -f ${container_name} > /dev/null 2>&1
docker run -d --network none --name ${container_name} docker.io/alpine:latest sleep inf
pid=$(docker inspect -f '{{.State.Pid}}' ${container_name})
ln -sfT /proc/$pid/ns/net /var/run/netns/${container_name}

host_veth="veth_"$pid"_ext"
container_veth="veth_"$pid"_int"

ip link add dev $host_veth type veth peer name $container_veth
ip link set dev $host_veth master $br_name
ip link set dev $host_veth up
ip link set dev $container_veth netns ${container_name}
ip netns exec ${container_name} ip link set dev $container_veth up
ip netns exec ${container_name} ip address add $container_ip dev $container_veth
ip netns exec ${container_name} ip route add default via $gw_ip
