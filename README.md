# Jetty

Turbo fast port scanner, made in Go.
Runs on average 3x faster than nmap, does not miss any ports.(at least in my testing :D)


## Flags

`-u` Enter URL

`-p` Enter ports to scan. Must be done in format `x-y` (to be fixed.)

`-t` How many threads you want to scan. Defaults to 20, I would recommend around 40 if you have a good internet connection. 

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

## TODO
Ordered by priority

Create mode to pipe in x urls and scan each one.

In scanning function, attempt to use ipv6.(20% faster apparently)

Map services to ports(via what x port is usually used for)

Add checking for 'filtered ports'  - https://nmap.org/book/man-port-scanning-basics.html

Add functionality to have nmap-like port specification.
Port flags to getopt package for better shit https://godoc.org/github.com/pborman/getopt

Add stdin piping scanning functionality

make a cool ASCII banner when it starts

Add nmap script functionality


Test change in memory usage from 1k cap channel to 65k cap channel. River channel to be exact.