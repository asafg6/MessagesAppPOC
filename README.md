# MessagesAppPOC
An "event driven" app POC, using HTTP2 server push and server sent events

You can read all about it at <a href="https://www.turtle-techies.com/post/dont-call-us-we-will-call-you/">https://www.turtle-techies.com/post/dont-call-us-we-will-call-you/</a>.

## How To Run It Locally

### Dependencies

* Redis listening on 6379 with no authentication (I recommend <a href="https://docs.docker.com/samples/library/redis/">docker</a>)
* Go 1.9
* npm 5.x and above

### Installing

```shell
git clone https://github.com/asafg6/MessagesAppPOC.git
cd MessagesAppPOC
go get github.com/go-redis/redis
go get github.com/asafg6/sse_handler
cd frontend
npm install
npm run build

```


### Running

on the project directory

```shell
go run main.go
```

Open your browser at https://localhost:8080 (Ignore the not secure warning)

Push a message to redis
```shell
 redis-cli publish events '{"id": 1, "data": "an important message", "type": "blue"}'
 ```


