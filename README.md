# tcp-extension-cable [![Build Status](https://travis-ci.org/lxjhk/tcp-extension-cable.svg?branch=master)](https://travis-ci.org/lxjhk/tcp-extension-cable)
Create a tunnel between a tcp port on A and another tcp port on B with simple masking of the transmitted bytes.

# Motivation
This program works like a reverse proxy in that you visit a machine that is not the real server but your traffic is forwarded to the real server. However, unlike a typical proxy service, this program masks the bytes that are being forwarded which is useful if somewhere along the original communication channel, packets containing a certain pattern is dropped.

Topology:

```
Before: client->server
After: client->springboard server->local-client server->server
```

The transmission in `springboard server->local-client server` is masked on the byte level.


# How to use
There are two modes: local-client mode and springboard mode. In a typical setup, one machine will be in the local-client mode which listens for connection from the springboard server and then establishes connections with the destination server. Another machine will be in the springboard mode which listens for incoming client connections and then establishes connection with the local-client server. 

When a client connects to the springboard server, the springboard server starts a coneection with the local-client server which then starts a connection with the actual server.

## local-client

Command to start in the local-client mode:

```
tcpec -mode=local-client -da=x.x.x.x:5000 -lp=8844
```

`da` represents the destination address. It represents the endpoint that the connection ultimately points to.

`lp` is the port number that springboard will use to connect to the local-client.


## springboard

Command to start in the springboard mode:

```
tcpec -mode=springboard -da=x.x.x.x:8844 -lp=8855
```

`da` should be set to the local-client address with local-client's `lp`

`lp` is the port number that clients will use to connect to the springboard.


