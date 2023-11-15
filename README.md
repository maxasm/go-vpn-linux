# How to create a VPN on Linux using a TUN device

1. Ensure that the `TUN` module is loaded.

```bash
lsmod | grep tun
```

If the above command produces no output then the TUN module is not loaded.

You can load the `TUN` module by running the command below.

```bash
sudo modprobe tun
```

If the above command produces and error this means the module is not installed.
Run the following command to install the TUN module (on Arch Linux).

```bash
sudo pacman -S linux-headers
```

2. Create the `TUN` interface 

You can use the `ip` command to create a `TUN` interface.

```bash
sudo ip tuntap add mode tun name <tun-interface-name>
```

3. Check if the `TUN` device has been created

```bash
ip link show dev <tun-interface-name>
```

4. Bring the device `UP`

You need to bring-up the interface.
When an interface is `up` it means that the operating system has configured the interface and it is ready to send and receive network traffic.
To bring up the interface, run the command below.

```bash
sudo ip link set <tun-interface-name> up
```

5. Assign an IP and Subnet Mask to the interface

Once the interface is up, you need to assign an IP and a Subnet Mask to it.
When assigning an IP and Subnet Mask, make sure that the network does not conflict with any other network.
You are actually acting as a router in the sence that a router assigns IP addresses.

```bash
sudo ip addr add 192.168.1.1/24 dev <tun-interface-name>
```

6. Delete the `TUN` interface

You can delete the TUN interface using the command below

```bash
sudo ip tuntap del dev <tun-interface-name> mode tun
```
