# Go WebSockets Chat Application.


The following features have been covered: Templates, Web sockets, OAuth2, Tracing(TDD),File, request handling and more


## Needed resources
```
> go get github.com/stretchr/gomniauth
> go get github.com/stretchr/objx
> go get github.com/gorilla/websocket
```
## OAuth2 keys:


Do some setup:


```
gomniauth.WithProviders(
		facebook.New("key", "secret",
			"http://localhost:**portListen**/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:**portListen**/auth/callback/github"),
		google.New("key", "secret",
			"http://localhost:**portListen**/auth/callback/google"),
	)
```
To get OAuth2 configured, replace respective key and secret in above code `main.go`. Key and secret got from the Provider's developer console.

## How do i get it running?

```
> go build -o chat
> ./chat --addr=":Serve at what port?"
> ./chat --addr=":8000"
```

Visit (http://localhost:3000/chat) to access application

## Nginx Configuration
Getting `Is not web socket upgrade` then have to upgrade the connection headers. Serving via nginx therefore do as follows


In **/etc/nginx/site-available/chat.conf** have the following as your nginx configuration file to upgrade connection header from http to socket.


```
upstream websocket {
        server localhost:8000;
    }

    server {
        proxy_read_timeout 86400s;
        proxy_send_timeout 86400s;
        server_name chatapp.net;
        listen 3000;
        location / {
            proxy_pass http://websocket;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
        }
    }
```

Change port listening at to appropriate port. 



To reload nginx, run cmd `sudo systemctl reload nginx.service`


## Sample output 


if `r.tracer = trace.New(os.Stdout)` in `main.go`


```
> 2018/04/21 16:36:44 Starting to serve at :8000
> New client joined
> New client left
> New client joined
> New client joined
> Message received: Hi room
> --Sent to client
> --Sent to client
> Message received: Hello eric. What are you coding?
> --Sent to client
> --Sent to client
> New client left
> New client left
```
## Required Resources
- Go Programming Blueprints by Mat Ryer




