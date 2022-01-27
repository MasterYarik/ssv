# SSV Specifications - Networking

**Status: WIP**

**TODO: add header**

This document contains the networking specification for `SSV.Network`.

## Overview

- [x] [Fundamentals](#fundamentals)
  - [x] [Stack](#stack)
  - [x] [Transport](#transport)
  - [x] [Messaging](#messaging)
  - [x] [Network Peers](#network-peers)
  - [x] [Identity](#identity)
  - [x] [Network Discovery](#network-discovery)
- [ ] [Wire](#wire)
  - [x] [Consensus](#consensus-protocol)
  - [ ] [Sync](#sync-protocol)
  - [ ] [Handshake](#handshake-protocol)
- [x] [Networking](#networking)
  - [x] [PubSub](#pubsub)
  - [x] [User Agent](#user-agnet)
  - [x] [Discovery](#discovery)
  - [x] [Netowrk ID](#network-id)
  - [x] [Subnets](#subnets)
  - [x] [Peers Connectivity](#peers-connectivity)
  - [x] [Forks](#forks)
  - [x] [High Availability](#high-availability)
  - [ ] [Security](#security)

## Fundamentals

### Stack

`SSV.Network` is a decentralized P2P network, consists of operator nodes grouped in multiple subnets.

The networking layer is built with [Libp2p](https://libp2p.io/), a modular framework for P2P networking that is used by multiple decetralized projects, inluding eth2.

### Transport

Network peers should support the following transports:
- `TCP` is used by libp2p for setting up communication channels between peers. default port: `12000`
- `UDP` is used for discovery purposes. default port: `13000`

[go-libp2p-noise](https://github.com/libp2p/go-libp2p-noise) is used to secure transport (based on [noise protocol](https://noiseprotocol.org/noise.html)).

Multiplexing of protocols over channels is achieved using [yamux](https://github.com/libp2p/go-libp2p-yamux) protocol.

### Messaging

Messages in the network are formatted with `protobuf` (NOTE: `v0` messages are encoded/decoded with JSON),
and being transported p2p with one of the following methods:

**Streams** 

Streams are used for direct messages between peers.

Libp2p allows to create a bidirectional stream between two peers and implement the corresponding wire messaging protocol. \
See more information in [IPFS specs > communication-model - streams](https://ipfs.io/ipfs/QmVqNrDfr2dxzQUo4VN3zhG4NV78uYFmRpgSktWDc2eeh2/specs/7-properties/#71-communication-model---streams).

**PubSub**

GossipSub ([v1.1](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md)) is the pubsub protocol used in `SSV.Network`

The main purpose is for broadcasting messages to a group (AKA subnet) of nodes. \
In addition, the machinary helps to determine liveliness and maintain peers scoring.


### Network Peers

There are several types of nodes in the network:

`Operator` is responsible for executing validators duties. \
It holds relevant registry data and the validators consensus data.

`Bootnode` is a public peer which is responsible for helping new peers to find other peers in the network.
It has a stable ENR that is provided with default configuration, so other peers could join the network easily.

`Exporter` is a public peer that is responsible for collecting and exporting information from the network. \
It collects registry data and consensus data (decided messages) of all the validators in the network. \
It has a stable ENR that is provided with default configuration, so it will have a stable connection with all peers and won't be affected by scoring, prunning, backoff etc.


### Identity

Identity in the network is based on two types of keys:

`Network Key` is used to create network/[libp2p identity](https://docs.libp2p.io/concepts/peer-id) (`peer.ID`), 
will be used by all network peers to setup a secured connection. \
Unless provided, the key will be generated and stored locally for future use, 
and can be revoked in case it was compromised. 

`Operator Key` is used for decryption of share's keys that are used for signing/verifying consensus messages and duties. \
Exporter and Bootnode does not hold this key.


### Network Discovery

[discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md) is used in `SSV.network` to complement discovery capalities that don't come with libp2p.

More information is available in [Networking / Discovery](#discovery)

------


## Wire

Network interaction includes several types of protocols, as detailed below.

## Consensus Protocol

`IBFT`/`QBFT` consensus protocol is used to govern `SSV` network.
`IBFT` ensures that consensus can be reached by a committee of `n` operator nodes while tolerating a certain amount of `f` faulty nodes as defined by `n ≥ 3f + 1`.

As part of the algorithm, nodes are exchanging messages with other nodes in the committee. \
Once the committee reaches consensus, the nodes will publish the decided message across the network.

More information regarding the protocol can be found in [iBFT annotated paper (By Blox)](/ibft/IBFT.md)

### Message Structure

`SignedMessage` is a wrapper for IBFT messages, it holds a message and its signature with a list of signer IDs:

```protobuf
syntax = "proto3";
import "gogo.proto";

// SignedMessage holds a message and it's corresponding signature
message SignedMessage {
  // message is the IBFT message
  Message message            = 1 [(gogoproto.nullable) = false];
  // signature is a signature of the IBFT message
  bytes signature            = 2 [(gogoproto.nullable) = false];
  // signer_ids is a list of the IDs of the signing operators
  repeated uint64 signer_ids = 3;
}

// Message represents an IBFT message
message Message {
  // type is the IBFT state / stage
  Stage type        = 1;
  // round is the current round where the message was sent
  uint64 round      = 2;
  // identifier is the message identifier
  bytes identifier      = 3;
  // sequence number is an incremental number for each instance, much like a block number would be in a blockchain
  uint64 seq_number = 4;
  // value holds the message data in bytes
  bytes value       = 5;
}
```

JSON example:
```json
{
  "message": {
    "type": 3,
    "round": 1,
    "identifier": "OTFiZGZjOWQxYzU4NzZkYTEwY...",
    "seq_number": 28276,
    "value": "mB0aAAAAAAA4AAAAAAAAADpTC1djq..."
  },
  "signature": "jrB0+Z9zyzzVaUpDMTlCt6Om9mj...",
  "signer_ids": [2, 3, 4]
}
```

**NOTE:** 
- all pubsub messages in the network are wrapped with libp2p's message structure
- `signer_ids` must be sorted, to allow hashing the entire message 

---

## Sync Protocols

There are several sync protocols, tha main purpose is to enable operator nodes to sync past decided message or to catch up with round changes.

In order to participat in some validator's consensus, a peer will first use sync protocols to align on past infromation.

Sync is done over streams as pubsub is not suitable in this case due to several reasons such as:
- API nature is request/response, unlike broadcasting in consensus messages
- Bandwidth - only one peer (usually) needs the data, it would be a waste to send redundant messages across the network.

### Message Structure

`SyncMessage` structure is used by all sync protocols, the type of message is specified in a dedicated field:

```protobuf
message SyncMessage {
  // MsgType is the type of sync message
  SyncMsgType Type                   = 1;
  // Identifier of the message (validator + role)
  bytes Identifier                      = 2;
  // Params holds the requests parameters
  repeated uint64 Params                = 3;
  // Messages holds the results (decided messages) of some request
  repeated proto.SignedMessage Messages = 4;
  // Error holds an error response if exist
  string Error                          = 5;
}

// SyncMsgType is an enum that represents the type of sync message 
enum SyncMsgType {
  // GetHighestType is a request from peers to return the highest decided/ prepared instance they know of
  GetHighestType       = 0;
  // GetInstanceRange is a request from peers to return instances and their decided/ prepared justifications
  GetInstanceRange     = 1;
  // GetCurrentInstance is a request from peers to return their current running instance details
  GetLatestChangeRound = 2;
}
```

Highest decided response:
```json
{
  "SignedMessages": [
    {
      "message": {
        "type": 3,
        "round": 1,
        "identifier": "...",
        "seq_number": 7943,
        "value": "Xmcg...sPM="
      },
      "signature": "g5y....7Dv",
      "signer_ids": [4,2,1]
    }
  ],
  "Type": 0,
  "Identifier": "..."
}
```

Error response:
```json
{
  "Identifier": "...",
  "Type": 2,
  "error": "EntryNotFoundError"
}
```

### Protocols

**TODO: add example request/response**

SSV nodes use the following stream protocols:


### 1. Highest Decided

This protocol is used by a node to find out what is the highest decided message for a specific validator.
In case there are no decided messages, it will return an empty array of messages.

`/ssv/sync/highest_decided/0.0.1`


### 2. Decided By Range

This protocol enables to sync decided messages in some specific range.

The request should specify the desired range, while the response will include all the found messages for that range.

`/ssv/sync/decided_by_range/0.0.1`


### 3. Last Change Round

This protocol enables a node to catch up with change round messages.

`/ssv/sync/last_change_round/0.0.1`

---

## Handshake protocol

The handshake protocol allows peers to identify, by exchanging signed information. \
It must be performed for every connection, and therefore forces nodes to 
authenticate / prove ownership of their operator key.

**TBD** Public, static nodes such as exporter requires registration 

`/ssv/auth/0.0.1`

The following information will be exchanged as part of the handshake:

```protobuf
syntax = "proto3";
import "gogo.proto";

// AuthMessage is a message that is used for authenticating nodes
message HandshakeMessage {
  // info is the node information to sign
  NodeInfo info = 1 [(gogoproto.nullable) = false];
  // signed is a signature of the message
  bytes signed     = 2 [(gogoproto.nullable) = false];
}

// NodeInfo contains node's information
message NodeInfo {
  // peer_id of the authenticating node
  bytes peer_id          = 1 [(gogoproto.nullable) = false];
  // operator_id of the authenticating node
  bytes operator_id      = 2 [(gogoproto.nullable) = true];
  // node_type is the type of the authenticating node
  uint64 node_type       = 3 [(gogoproto.nullable) = false];
  // execution_node is the eth1 node used by the node
  string execution_node  = 4; // TBD: actual value
  // consesnsus_node is the eth2 node used by the node
  string consesnsus_node = 5; // TBD: actual value
  // fork_version is the current fork used by the node
  uint32 fork_version    = 6;
  // fork_version is the current ssv-node version
  string node_version    = 7;

  // TODO: add cloud provider / region / ... 
}
```

---


## Networking

### Pubsub

The main purpose is for broadcasting messages to a group (AKA subnet) of nodes. \
In addition, the following are achieved as well:

- subscriptions metadata helps to get liveliness information of nodes
- pubsub scoring enables to prune bad/malicious peers based on network behavior and application-specific rules

The following parameters are used for initializing pubsub:

- `floodPublish` was turned on for better reliability, as peer's own messages will be propagated to a larger set of peers 
  (see [Flood Publishing](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#flood-publishing)) 
- `peerOutboundQueueSize` / `validationQueueSize` were increased to `600`, to avoid dropped messages on bad connectivity or slow processing
- `directPeers` includes the exporter peer ID, to ensure it gets all messages
- `subscriptionFilter` was injected to ensure a peer will connect to relevant topics, see [SubscriptionFilter](https://github.com/libp2p/go-libp2p-pubsub/blob/master/subscription_filter.go) interface
- (fork `v1`) `msgID` is a custom function that calculates a `msg-id` based on the message content hash. 
The default function uses the `sender` + `msg_seq` which we don't track, and enforces signature / verification for each message. 
As all the messages are being verified using the share key, it would be redundant to it also in the pubusb level.
Moreover, the default `msg-id` duplicates messages, causing it to be processed more than once, in case it was sent by multiple peers (e.g. decided message).
- (fork `v1`) `signPolicy` was set to `StrictNoSign` (required for custom `msg-id`) to avoid producing and verifying message signatures in the pubsub router
  - `signID` was set to empty (no author)

#### Pubsub Scoring

[Peer scoring](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#peer-scoring) was introduced in `gossipsub v1.1`,
the idea is that each individual peer maintains a score for other peers. 
The score is locally computed by each individual peer based on observed behaviour and is not shared.

`SSV.network` injects application specific scoring to apply connection and message scoring as part of pubsub scoring system. \
See [Connection Scoring](#connection-scoring) and [Message Scoring](#message-scoring) for more information.

Score thresholds are used by libp2p to determine whether a peer should be removed from topic's mesh, penalized or even ignored if the score drops too low. \
See [this section](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#score-thresholds) for more details regards the different thresholds. \
Thresholds values **TBD**, this section will be updated once that work is complete:

- `gossipThreshold`: -4000
- `publishThreshold`: -8000
- `graylistThreshold`: -16000
- `acceptPXThreshold`: 100
- `opportunisticGraftThreshold`: 5

[Score function](https://github.com/libp2p/specs/blob/master/pubsub/gossipsub/gossipsub-v1.1.md#the-score-function) is running pariodically and will determine the score of peers. During heartbeat, the score it checked and bad peers are handled accordingly.


### User Agent

Libp2p provides user agent mechanism with the [identify](https://github.com/libp2p/specs/tree/master/identify) protocol, 
which is used to exchange basic information with other peers in the network.

User Agent contains the node version and type, and in addition the operator id which might be reduced in future versions. \
See detailed format in [Forks / user agent](#fork-v0)


### Discovery

[discv5](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md) is a system for finding other participants in a peer-to-peer network, 
it is used in `SSV.network` to complement discovery.

DiscV5 works on top of UDP, it uses a DHT to store node records (`ENR`) of discovered peers. \
The discovery process walks randomaly on the nodes in the table that are not connected, filters them by `ENR` entries in order to connect with the most relevant nodes.

The communication is encrypted and authenticated using session keys, 
established in the [handshake process](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-theory.md#sessions).

**Bootnode** 

A peer that have a public, static ENR to enable new peers to join the network. For the sake of flexibility, 
bootnode/s ENR values are configurable and can be changed on demand by operators. \
Bootnode doesn't start a libp2p host for TCP communication, its role ends once a new peer finds existing peers in the network.

#### ENR

[Ethereum Node Records](https://github.com/ethereum/devp2p/blob/master/enr.md) is a format that holds peer information.
Records contain a signature, sequence (for republishing record) and arbitrary key/value pairs. 

`ENR` structure in `SSV.Network` consists of the following key/value pairs:

| Key         | Value                                                          | Status          |
|:------------|:---------------------------------------------------------------|:----------------|
| `id`        | name of identity scheme, e.g. "v4"                             | Done            |
| `secp256k1` | compressed secp256k1 public key, 33 bytes                      | Done            |
| `ip`        | IPv4 address, 4 bytes                                          | Done            |
| `tcp`       | TCP port, big endian integer                                   | Done            |
| `udp`       | UDP port, big endian integer                                   | Done            |
| `type`      | node type, integer; 1 (operator), 2 (exporter), 3 (bootnode)   | Done (`v0.1.9`) |
| `oid`       | operator id, 32 bytes                                          | Done (`v0.1.9`) |
| `version`   | fork version, integer                                          | -               |
| `subnets`   | bitlist, 0 for irrelevant and 1 for assigned subnet            | -               |

#### Discovery Alternatives

[libp2p's Kademlia DHT](https://github.com/libp2p/specs/tree/master/kad-dht) offers similar features, and even a more complete implemetation of Kademlia DHT.
Discv5 design is loosely inspired by the Kademlia DHT, but unlike most DHTs no arbitrary keys and values are stored. Instead, the DHT stores and relays node records.

Libp2p's Kad DHT will allow to advertise and find peers by multiple keys, e.g. topics/subnets.
As `ENR` has a size limit (`< 300` bytes), and therefore discv5 won't support multiple key/value pairs.

Notes:
- discv5 specifies [Topic Index](https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md#the-topic-index) that help to lookup relevant nodes in a smaller set of nodes, currently not fully implemented
- [discv5 specs](https://github.com/ethereum/devp2p/blob/master/discv5/discv5.md#comparison-with-other-discovery-mechanisms) details why discv5 was chosen over libp2p Kad DHT in Ethereum.

In `v0` discv5 is used, `v1` TBD, this section will be updated once that work is complete.

### Network ID

Network ID is a `32byte` key, that is used to distinguish between other networks (ssv and others).
Peers from other public/private libp2p networks (with different network ID) won't be able to read or write messages in the network, 
meaning that the key to be known and used by all network members.

It is done with [libp2p's private network](https://github.com/libp2p/specs/blob/master/pnet/Private-Networks-PSK-V1.md),
which encrypts/decrypts all traffic with the corresponding key,
regardless of the regular transport security protocol ([go-libp2p-noise](https://github.com/libp2p/go-libp2p-noise)).

**NOTE** discovery communication (UDP) won't use the network ID, but unknown nodes will be filtered anyway due to missing fields in their `ENR` entry as specified in [Discovery section](#discovery).


### Subnets

Consensus messages are being sent in the network over a subnet (pubsub topic), which the relevant peers should be subscribed to.

In addition to subnets, there is a global topic (AKA `main topic`) to publish all the decided messages in the network.

There are several options for how to setup topics in the network:

#### Subnets - fork v0

Each validator committee has a dedicated pubsub topic with all the relevant peers subscribed to it (committee + exporter).

It helps to reduce amount of messages in the network, but increases the number of topics which will grow up to the number of validators.

#### Subnet - fork v1

**TBD: main topic**

A subnet of operators is responsible for multiple committees,
reusing the same topic to communicate on behalf of multiple validators.

In comparison to `v0`, the number of topics will be reduced and the number of messages sent over the network should grow. \
As messages will be propagated to a larger set of nodes, we can expect better reliability (arrival of messages to all operators in the committee).

In addition, a larger group of operators provides:
- redundancy of decided messages across multiple nodes
- better security for subnets as more nodes will validate messages and can score bad/malicious nodes that will be pruned accordingly.

**Validators Mapping**

Validator's public key is mapped to a subnet using a hash function, 
which helps to distribute validators across subnets in a balanced way:

`hash(validatiorPubKey) % num_of_subnets`

Deteministic mapping is ensured as long as the number of subnets doesn't change, 
thefore its a fixed number (TBD 32 / 64 / 128).

**TBD** A dynamic number of subnets (e.g. `log(num_of_peers)`) which requires a different approach.


### Message Scoring

Message scorers track on operators behavior w.r.t incoming IBFT messages:

- Invalid message signature (`-100`)
- Message from operator w/o shared committees (`-1000`)

### Connection Scoring

Peer's connection score is determined after a successful handshake, and peers with low score will be pruned. 

Connection scores are based on the following properties:

- Shared subnets / committees (`25`)
- Static nodes such as `exporter` (`10000`)


### Peers Connectivity

In a fully-connected network, where each peer is connected to all other peers in the network,
running nodes will consume many resources to process all network related tasks e.g. parsing, peers management etc.

To lower resource consumption, the number of connected peers is limited, currently set to `250`. \
Once reached to peer limit, the node will connect only to relevant nodes with score above treshold, which is currently set to zero.


#### Connection Filters

Connection filters are executed upon new connection. \
Filters calculates the connection score of the new peer, and will terminate the connection if score is low.
In addition, it will mark the peer as pruned so following connections requests will be stopped at connection gater.

#### Connection Gating

Connection Gating allows to safeguard against bad/pruned peers that tries to connect multiple times. 
Inbound and outbound connections are intercepted and being checked before other components process the connection.

See libp2p's [ConnectionGater](https://github.com/libp2p/go-libp2p-core/blob/master/connmgr/gater.go) interface for more info.

### Forks

Future network forks will follow the general forks mechanism and design in SSV. \
The idea is to wrap procedures that have potential to be changed in future versions.

Currently, the following are covered:

- validator topic mapping
- message encoding/decoding
- user agent

#### Fork v0

**validator topic mapping**

Validator public key is used as the topic name:

`bloxstaking.ssv.<hex(validator-public-key)>`

**message encoding/decoding**

JSON is used for encoding/decoding of messages.

**user agent**

User Agent contains the node version and type, and in addition the operator id.

`SSV-Node:v0.x.x:<node-type>:<?operator-id>`

#### Fork v1 (TBD)

**validator topic mapping**

Validator public key hash is used to determine the validator's subnet which is the topic name:

`bloxstaking.ssv.<hash(validatiorPubKey) % num_of_subnets>`


### NAT port map

libp2p enables to configure a `NATManager` that will attempt to open a port in the network's firewall using `UPnP`.


### Security

The following measures are used to protect against malicious peers and denial of service attacks:
- The number of connected peer is limited to `250`. Once reaching limit, the node should connect only with peers that have high connection score.
- Connection score that is determined during discovery process, ensures that the node will try to connect only with relevant peers.
- Connection filters determine the score for inbound connections, 
- Connection gater protects against peers which were pruned in the past but tries to connect again before backoff timeout (5 min). 
it kicks in in an early stage, before the other components processes the request to avoid resources consumption.

DiscV5 specs specifies potential vulnerabilities in the discovery system, 
see (discv5-rationale/security-goals)[https://github.com/ethereum/devp2p/blob/master/discv5/discv5-rationale.md#security-goals].
Not all vulnerabilities applies to `SSV.Network`, and some were mitigated.

**TODO: complete**

### High Availability

As it participats in a decentralized p2p network, HA for an ssv node is not trivial. \ 
The reason is that running multiple instances might lead to slashing and even disturb the consensus as it creates ambiguity and conflicts.

Ideas TBD:

#### Hub Node

`Hub Node` is a node that handles the network layer of ssv node, 
it could be connected to multiple `Worker` nodes and stream the entire network layer messages to/from them.

A possible implmentation could be an SSV node with a proxied network layer that uses websocket to communicate with worker nodes.

#### Subnet Partitions

Subnet partitions separates a given set of subnets into `n` indenpendant subsets of subnets, 
which could be assigned to `n` running instances of the same operator, each working on different subsets.

That will help decrease the damage in case some node fails, as only a portion of the assigned validators will be affected, while the other healthy instances keeps doing tasks in their subnets.