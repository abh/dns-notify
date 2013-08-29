# dns-notify

Small utility to easily send DNS NOTIFY requests.

## Install

First you have to [install and setup Go](http://golang.org/doc/install).

Then run

	go get github.com/abh/dns-notify

That should put the binary into $GOPATH/bin/dns-notify

## API

If started with the `-listen` parameter the program runs an HTTP server with an API
for sending notifies.

The path for the POST requests is`/api/v1/notify/example.com`. The server will use the
servers specified on startup, but the domain is taken from the API request.

    curl -XPOST -s http://localhost:8765/api/v1/example.com

The response is JSON and includes the result of each NOTIFY sent.

## Usage


### Send a notify mode
 
Specify the domain to notify with the -domain parameter and then the servers to notify.

    dns-notify -domain=example.com 127.0.0.1 10.0.0.1 192.168.0.1:5053 [2001:1::2]:53

### Daemon mode

To start dns-notify in daemon mode, specify the `-listen` parameter
instead of `-domain`.

   dns-notify -listen=10.0.0.1:8050 127.0.0.1 10.0.0.1

### Optional parameters

* -verbose

A bit more verbose output.

* -quiet

Only output on errors.

* -timeout

How long to wait for responses, defaults to 2000 milliseconds (2 seconds).

## Error handling

Errors are reported, but requests are not retried.
