# dorp
## The K2CC Door Program

Dorp allows remote monitoring of K2CC station occupancy status by giving access to two pieces of information:
whether the door is open or shut, and whether the lights are on or off.

Dorp consists of two programs:
- A server, which displays current state on a web site
- A client, which pushes updates from the station to the server.

The client is intended to run on a Raspberry Pi, connected to a switch that monitors the door and
a photosensor to monitor the lights.

The server is intended to run on a computer with a static IP address.

Further documentation on how each of the programs work is in each of their respective directories.
Documeentation on the protocol can be found in protocol.md.
