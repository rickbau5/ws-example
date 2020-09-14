# ws-example
`ws-example` shows an simple example implementation of a Golang Websocket server. It is primarily an exploration of 
[`gorilla/websockets`](https://github.com/gorilla/websocket) and `nginx` to proxy websocket connections.

Using `nginx` as a proxy, we can run multiple websocket servers and distribute connections among them. 

See details below for running the application locally.

## Running Locally

### docker-compose
Simply running `docker-compose up` is enough to get the services up and functional. Connect a websocket client to
`http://localhost:8080/ws` to interact with the websocket server.

An `nginx` service will start which proxies websocket connections to `ws-backend-n`. `nginx` is used to provide
a single host to connect to (e.g. `localhost:8080`) but distribute requests across multiple backends.

The `ws-backend-n` services are a simple Golang application using [`gorilla/websocket`](https://github.com/gorilla/websocket)
that echoes back any message received but prepends the hostname to differentiate which backend host you are connected to.

Use `websocat` to connect as a client for demonstration purposes. It can be installed with Homebrew (`brew install websocat`).

For example:
```bash
$ websocat -E --ping-interval 5 --ping-timeout 10 ws://localhost:8080/ws
```
This gives an interactive shell to send and received messages to the backend. Type a message and press enter to see a response.

```bash
$ websocat -E --ping-interval 5 --ping-timeout 10 ws://localhost:8080/ws
[a9de44781391] hello :) type a message and press ENTER
hello
[a9de44781391] hello
```

Press `CTRL+C` to exit, with the `-E` this will close the connection immediately with the backend as well.

### Docker Swarm
Docker Swarm can be used to emulate an environment that would be more similar to production. The `docker-compose.swarm.yml`
can be used to deploy the service stack with one `nginx` service and, initially, one `ws-backend` service. From there,
the `ws-backend` service can be scaled arbitrarily and `nginx` will automatically proxy requests to all backends.

**Deploy `websockets` stack using `docker-compose.swarm.yml`**
```bash
$ docker stack deploy --compose-file=docker-compose.swarm.yml websockets
Creating service websockets_nginx
Creating service websockets_ws-backend
```

**View tasks in the `websockets` stack**
```bash
$ docker stack ps websockets
ID                  NAME                      IMAGE                  NODE                DESIRED STATE       CURRENT STATE                ERROR                       PORTS
w992fkm1363w        websockets_nginx.1        nginx:1.19.2           carb                Running             Running about a minute ago                               
54ipaz8sal94        websockets_ws-backend.1   golang:1.15.2-buster   carb                Running             Running about a minute ago
```

**Scale `ws-backend` service in stack**
```bash
$ docker service scale websockets_ws-backend=2
websockets_ws-backend scaled to 2
overall progress: 2 out of 2 tasks 
1/2: running   [==================================================>] 
2/2: running   [==================================================>] 
verify: Service converged 
```

**Use `websocat` to connect to the services**
```bash
$ websocat -E --ping-interval 5 --ping-timeout 10 ws://localhost:8080/ws
[886396d2edcc] hello :) type a message and press ENTER
asdf
[886396d2edcc] asdf
...
$ websocat -E --ping-interval 5 --ping-timeout 10 ws://localhost:8080/ws
[1eba9cf95b52] hello :) type a message and press ENTER
adsf
[1eba9cf95b52] adsf
```
Note that we have connected to two different hosts `886396d2edcc` and `1eba9cf95b52`. This indicates that `nginx`
is correctly proxying connections to an arbitrary number of backends and we are therefore naively load balancing the
websocket connections.

**View logs for `ws-backend` services**
```bash
docker service logs -f websockets_ws-backend
websockets_ws-backend.2.plyyp6ztzydh@carb    | Got message: asdf
websockets_ws-backend.1.54ipaz8sal94@carb    | Got message: adsf
```
