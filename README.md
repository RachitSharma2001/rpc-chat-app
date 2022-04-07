# Concurrent Chat Application following RPC protocol
This is an Go application that makes use of gRPC, protocol buffers, and concurrency in order to make a concurrent chat application. What this means is that multiple clients can send messages to each other, and they will recieve messages concurrently. 

# How to run 
On one terminal, cd into the root directory of this project, and then do <b> "go run main.go" </b>. Then, on another terminal, cd into the client directory, and do <b> "go run client.go" </b>. You can set up as many client terminals as you want. 