# Jetty

Turbo fast port scanner, made in Go.
Runs on average 3x faster than nmap, does not miss any ports.(at least in my testing :D)


## Flags

`-u` Enter URL

`-p` Enter ports to scan. Must be done in format `x-y` (to be fixed.)

`-t` How many threads you want to use for scanning. Defaults to 20, I would recommend around 40 if you have a good internet connection. 

`-v` Enable verbose logging.

`-timeout` How long to wait for the server to respond to the request. Defaults to 500, adjust for your connection.



## Installation


Install with `go get -u github.com/orsetii/jetty`

Or build from source: 
```
git clone https://github.com/orsetii/jetty

cd $GOPATH/src/github.com/jetty

go install
```
