# ObfsVPN

The `obfsvpn` module contains a Go package that provides server and client components to
use variants of the obfs4 obfuscation protocol. It is intended to be used as a
drop-in Pluggable Transport for OpenVPN connections (although it can be used
for other, more generic purposes).

A docker container will be provided to facilitate startng an OpenVPN service that
is accessible via the obfuscated proxy too.

You can read more online about how obfsvpn is used to provide [circumvention](https://docs.leap.se/circumvention/)
tactics to the LEAP VPN Clients, and in particular about the design of the
[Hopping Pluggable Transport](https://docs.leap.se/circumvention/hopping/).

## Protocol stack

```
--------------------
 application data
--------------------
      OpenVPN
--------------------
   obfsvpn proxy
--------------------
       obfs4
--------------------
   wire transport
--------------------
```

- Application data is written to the specified interface (typically a `tun`
  device started by `OpenVPN`).
- `OpenVPN` provides end-to-end encryption and a reliability layer. We'll be
  testing with the `2.5.x` branch of the reference OpenVPN implementation.
- `obfs4` is used for an extra layer of encryption and obfuscation. It is a
  look-like-nothing protocol that also hides the key exchange to the eyes of
  the censor.
- Wire transport is, by default, TCP. Other transports will be explored to
  facilitate evasion: `KCP`, `QUIC`?

## Testing

There is an entirely automated docker-compose based network sandbox that can be used for testing.

To bring up all of the services in the traditional/non-hopping mode:

```sh
$ make integration
```

This will build and launch the services in their correct configurations. You can then see them running:

```sh
$ docker-compose ps
          Name                        Command               State         Ports
--------------------------------------------------------------------------------------
obfsvpn_client_1           dumb-init /usr/bin/start.sh      Up
obfsvpn_obfsvpn-1_1        dumb-init /opt/obfsvpn/sta ...   Up
obfsvpn_obfsvpn-2_1        dumb-init /opt/obfsvpn/sta ...   Up
obfsvpn_openvpn-server_1   dumb-init /opt/openvpn-ser ...   Up      5540/tcp, 5540/udp
```

You can get logs from one, more, or all of the services:

```sh
$ docker-compose logs client
$ docker-compose logs client openvpn-server
# to tail all logs:
$ docker-compose logs -f
```

You can then run arbitrary commands on any of the services to debug, test performance, etc:

```sh
$ docker-compose exec client ip route
0.0.0.0/1 via 10.8.0.5 dev tun0
default via 192.168.80.1 dev eth0
10.8.0.1 via 10.8.0.5 dev tun0
10.8.0.5 dev tun0 proto kernel scope link src 10.8.0.6
127.0.0.1 via 192.168.80.1 dev eth0
128.0.0.0/1 via 10.8.0.5 dev tun0
192.168.80.0/20 dev eth0 proto kernel scope link src 192.168.80.3

$ docker-compose exec client ping 8.8.8.8
PING 8.8.8.8 (8.8.8.8): 56 data bytes
64 bytes from 8.8.8.8: seq=0 ttl=109 time=12.495 ms
64 bytes from 8.8.8.8: seq=1 ttl=109 time=13.614 ms
64 bytes from 8.8.8.8: seq=2 ttl=109 time=13.900 ms

$ docker-compose exec openvpn-server iperf3 -s --bind-dev tun0

❯ docker-compose exec client iperf3 -c 10.8.0.1 --bind-dev tun0
Connecting to host 10.8.0.1, port 5201
[  5] local 10.8.0.6 port 51390 connected to 10.8.0.1 port 5201
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec  40.9 MBytes   343 Mbits/sec  206    801 KBytes
[  5]   1.00-2.00   sec  36.2 MBytes   304 Mbits/sec  106    601 KBytes
[  5]   2.00-3.00   sec  41.2 MBytes   346 Mbits/sec    6    456 KBytes
[  5]   3.00-4.00   sec  40.0 MBytes   336 Mbits/sec    0    512 KBytes
[  5]   4.00-5.00   sec  42.5 MBytes   357 Mbits/sec    0    565 KBytes
[  5]   5.00-6.00   sec  43.8 MBytes   367 Mbits/sec    0    615 KBytes
[  5]   6.00-7.00   sec  36.2 MBytes   304 Mbits/sec   22    457 KBytes
[  5]   7.00-8.00   sec  41.2 MBytes   346 Mbits/sec    0    525 KBytes
[  5]   8.00-9.00   sec  41.2 MBytes   346 Mbits/sec    0    575 KBytes
[  5]   9.00-10.00  sec  40.0 MBytes   336 Mbits/sec    0    610 KBytes
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate         Retr
[  5]   0.00-10.00  sec   403 MBytes   338 Mbits/sec  340             sender
[  5]   0.00-10.02  sec   401 MBytes   336 Mbits/sec                  receiver

iperf Done.
```

### Testing PT3 Hopping

The PT3 Hopping architecture can be brought up in an almost identical way, except that calls to docker-compose require an `--env-file ./.env.hopping` parameter to distinguish between the two strategies.

The integration test will bring everything up for you:

```sh
$ make integration-hopping
```

Then when you want to run commands, add the `--env-file` argument:

```sh
❯ docker-compose --env-file ./.env.hopping exec client ping -c 3 -I tun0 8.8.8.8
PING 8.8.8.8 (8.8.8.8): 56 data bytes
64 bytes from 8.8.8.8: seq=0 ttl=113 time=12.829 ms
64 bytes from 8.8.8.8: seq=1 ttl=113 time=19.346 ms
64 bytes from 8.8.8.8: seq=2 ttl=113 time=19.013 ms

--- 8.8.8.8 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 12.829/17.062/19.346 ms
```

...

## Android

Assuming you have the `android ndk` in place, you can build the bindings for android using `gomobile`:

```
make build-android
```
