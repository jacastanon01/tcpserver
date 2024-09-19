# TCP Server in Go

This application will demonstrate how to set up a TCP connection using Go. 

## What is TCP?

Transmission Control Protocol (TCP) is built on top of the Internet Protocol Suite (TCP/IP) and establishes a connection between two computers. A request is sent from the client to the server and the response is sent back. A TCP connection guarantees delivery because data packets are transmitted in segments, and the protocol ensures that packets are received in the correct order. If any packet is lost or out of order, TCP handles retransmission, maintaining reliable data transfer. This connection is handled through a three-step "handshake" process: SYN, SYN-ACK, and ACK, which synchronizes both ends. Once the connection is established, data is transmitted in segments, and TCP ensures that all segments arrive intact and in order.
