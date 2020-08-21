# jetty

Turbo fast port scanner.

#TODO
Ordered by priority

Abstract scanning sections into different file. - Create Function to print both ports as they are scanned as done above; and to print a nice summary at the end
In scanning function, attempt to use ipv6.

Map services to ports(via what x port is usually used for)
Add checking for 'filtered ports'  - https://nmap.org/book/man-port-scanning-basics.html
Add functionality to have nmap-like port specification.
Add stdin piping scanning functionality
make a cool ASCII banner when it starts
Add nmap script functionality - this will be very hard.
Port flags to getopt package for better shit https://godoc.org/github.com/pborman/getopt


Test change in memory usage from 1k cap channel to 65k cap channel. River channel to be exact.