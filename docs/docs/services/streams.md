# Bidirectional Streaming

The **Stream service** defines abstract interfaces for bidirectional streaming between a client and a server. It exposes simple and efficient API for opening bidirectional streams, and recycling streams when opening multiple connections.

Streams are an efficient way to send multiple messages back and forth between two nodes. This is extremely useful in decentralized network where communication overhead can massively reduce performance, and communication happens at scale.

## Streamers

A common issue that arises when using bidirectional streams is figuring out which node should take on the role of the client, and which node should take on the role of the server. The **Streamer** interface is designed to hide these details from the user.

Below is an example of using the default implementation available in the `stream` package.

```go
streamer := NewStreamer(addr, client, server)
stream, err := streamer.Open(remoteMultiAddr)
if err != nil {
    log.Fatalf("cannot open stream with %v: %v", remoteMultiAddr, err)
}
defer stream.Close()
```

Opening a stream requires no knowledge of whether we should actively attempt a connection against a remote server, or passively wait for a connection from a remote client. This simplifies the control flow for the user and prevents nodes from accidentally creating two connections with each other, reducing resource usage.

## Recycling

It is often convenient to create multiple streams to a single node; simplifying control flows and ownership. However, done naively this can consume a lot of resources since maintaining a large number of network connections can be heavy.

Recycling streams allows us to reduce resource usage by reusing a single underlying connection for multiple streams.

```go
streamer := NewStreamer(addr, client, server)
streamer = NewStreamRecycler(streamer)
stream, err := streamer.Open(remoteMultiAddr)
if err != nil {
    log.Fatalf("cannot open stream with %v: %v", remoteMultiAddr, err)
}
defer stream.Close()
```

Using the **Streamer** interface, we are able to use the recycler without making any changes to our existing code. Once all streams to a node are closed, the underlying connection will be automatically freed.