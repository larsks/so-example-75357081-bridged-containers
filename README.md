# Bridged containers

I created this repository while investigating <https://stackoverflow.com/q/75357081/147356>.

I haven't been able to reproduce the behavior described in the question:

- It works with the original call to `netlink.LinkSetHardwareAddr`
- It works if we remove that call completely
- It works if we set an explicit hardware address

## Before you begin

In order to ping or otherwise access external addresses, we need to make the following changes:

- Ensure that forwarding isn't blocked by a rule or default policy in the filter `FORWARD` chain
- Ensure that we have an appropriate `MASQUERADE` rule in the nat `POSTROUTING` chain

If you're running on an Ubuntu system, the `rules.v4` file included here will set that up for you:

1. Install `iptables-persistent`:

    ```
    apt -y install iptables-persistent
    ```

2. Copy `rules.v4` to `/etc/iptables`:

    ```
    cp rules.v4 /etc/iptables/rules.v4
    ```

3. Reboot.

## Running the tests

The `run_test.sh` script will set up a bridge, and then start attaching containers. After attaching each container, it will attempt to ping a remote target to verify that connectivity is working as expected.

For each container we add to the bridge, we ping the remote target both from the new container and from the first container.
