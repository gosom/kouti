## Hello world with kouti

An example on how you can build a web server using kouti.

Build and start the server:

```
go build
./hello-world
```

Then in another terminal try:

```
curl http://localhost:8080
curl http://localhost:8080/panic
```

and view the logs of the server

The servers terminates gracefully (try that using Ctrl+C)

