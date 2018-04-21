# Go WebSockets Chat Application.
## How do i get it running?

```
> go build -o chat
> ./chat --addr=":Serve at what port?"
> ./chat --addr=":8000"
```
###### Sample output if `r.tracer = trace.New(os.Stdout)`\n

`Output here`

## Getting `Is not web socket upgrade` then have to upgrade the connection headers. Serving via nginx therefore do as follows

###### In **/etc/nginx/site-available/chat.conf** have the following as your nginx configuration file to upgrade connection header from http to socket.


```
upstream websocket {
        server localhost:8000;
    }

    server {
        proxy_read_timeout 86400s;
        proxy_send_timeout 86400s;
        server_name chatapp.net;
>        listen 3000;
        location / {
            proxy_pass http://websocket;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }
    }
```

Change port listenig at to appropriate port. 




