# Go WebSockets Chat Application.\n
## How do i get it running?\n

```
> go build -o chat
> ./chat --addr=":Serve at what port?"
> ./chat --addr=":8000"
```
###### Sample output if `r.tracer = trace.New(os.Stdout)`\n

* 2018/04/19 13:37:50 Starting to serve at: :8000\n
New client joined
New client joined
Message received: Hi room?

--Sent to client
--Sent to client
Message received: How are you doing
--Sent to client
--Sent to client
Message received: Socket programming can be fun!

--Sent to client\n
--Sent to client *
