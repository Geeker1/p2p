

# P2P Implementation

# Server
Server has a location to big file, splits it into chunks and sends to the various client servers.


# How do clients know each other

# Tracker
A seperate server that keeps track of servers and the chunks they have, so clients can query tracker for specific server having the chunk they want.


# Gossip Protocol
Clients are able to talk to themselves and figure out where the chunks they want are.

## Server Process
- Server starts, splits file into chunks and stores them, waiting for available clients.
- When clients become available, it sends chunks to clients, sharing load in best way possible.

## Client Process
- Client starts, queries server for chunk-list to keep track of.
- For every chunk client receives, it tells tracker it now has that chunk.
- Client checks list of chunks left to fetch, picks one and request tracker to send peer having that chunk.
    - Note that client can also receive chunk from server and if chunk is in list of chunks to fetch, client accepts chunk
- once chunk list is empty, we know we've fetched all chunks. We rearrange chunks to file and save in client directory.

## Tracker Process
- Stores information only about chunk-peer locations and updates chunk-store with new peer when client sends chunk-finished request.

## Noticing Caveats already
- There would probabaly be a state issue on client side if we are faced with issue where server sends chunk and client also receiving chunk from another peer.


