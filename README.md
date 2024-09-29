# PingAnalyzer
A tool to ping destinations with Goroutines. Different services have different methods to interact with each other. You can find more in the architecture section.

## Architecture
The "ping-service" is the core module of ping functionality. It is responsible for "start" and "stop" pings with desired options. Upon receiving a ping request, it stores the request information and generates a new separate goroutine to do the ping. Using a specific boolean channel, the ping-service can stop the allocated goroutine. The "server" is the central unit that interacts with different parts like CLI or Web GUI (to be developed). A server-side streaming gRPC is implemented between the ping-service and the server in order for ping replies to be transferred to the central unit continuously.


<p align="center">
<img src="./assets/Architecture%20diagram.png" width="500" />
</p>


## How to use?
Clone and build the source code. Each service gives you an executable (server and ping service for now). Run both of them. Temporarily, a simple CLI is implemented in the server service. You can run these commands there:

```cmd
ping <destination IP> <count>
stop <ping id number>
```


<p align="center">
<img src="./assets/example.png" width="400" />
</p>
