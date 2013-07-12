# dns-notify

Small utility to easily send DNS NOTIFY requests.

## Install

First you have to [install and setup Go](http://golang.org/doc/install).

Then run

	go get github.com/abh/dns-notify

That should put the binary into $GOPATH/bin/dns-notify

## Usage

Specify the domain to notify with the -domain parameter and then the servers to notify.

    dns-notify -domain=example.com 127.0.0.1 10.0.0.1 192.168.0.1:5053 [2001:1::2]:53

Optional parameters:

* -verbose

A bit more verbose output.

* -quiet

Only output on errors.

* -timeout

How long to wait for responses, defaults to 2000 milliseconds (2 seconds).