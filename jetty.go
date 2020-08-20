package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	//"net/http"
	"log"
	"net/url"
	"strings"
	"sync"
)

var (
	URL     string
	err     error
	d       net.Dialer
	Threads int
	verbose bool
	wg      sync.WaitGroup
	River   chan int    = make(chan int, 1000)
	Results chan result = make(chan result, 500)
)

type Ports []int

func (p *Ports) String() string {

	return "Ports Slice" // Not sure what this is for? Needed for interface satisfaction though.
}

func (p *Ports) Set(s string) (err error) {
	go func(err error) error {
		for err == nil {
			list := strings.Split(s, "-")
			iStart, err := strconv.Atoi(list[0])
			iEnd, err := strconv.Atoi(list[1])
			for ; iStart <= iEnd; iStart++ {
				if iStart < 10 {
					iStart = 10
				}
				*p = append(*p, iStart)
				// Send each port to river struct
				River <- iStart
			}
			return err
		}
		return err
	}(err)
	return err
}

var Target struct {
	harbour Ports
	addrs   []string
}

type result struct {
	port int
	err  error
}

func main() {

	flag.StringVar(&URL, "u", "", "The URL you want to scan.")
	flag.IntVar(&Threads, "t", 10, "The amount of threads you want to scan with.")
	flag.Var(&Target.harbour, "p", "What ports do you want to scan.('-p 0-1024')")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output.")
	flag.Parse()

	if _, err := url.Parse(URL); err != nil {
		fmt.Printf("Error ocurred in parsing your URL.\nPlease Try again.")
	}

	fmt.Printf("Resolving URL...\n")

	// TODO add context package usage maybe?

	Target.addrs, err = net.LookupHost(URL)
	d.Timeout = 500 * time.Millisecond // Wait 500ms for response, if not give up.
	if err != nil {
		fmt.Printf("Error occurred while trying to resolve the URL\n Error: %s\n", err)
	}

	// River channel to send ports on for go routines to pick up when they can.
	defer close(River)
	defer close(Results)
	defer wg.Done()
	// for i := 0; i < len(Target.harbour); i++ {
	// 	River <- Target.harbour[i]
	// }
	for i := 0; i < Threads; i++ {
		go ring(River, Results)
		wg.Add(1)
	}
	for i := 0; i < len(Target.harbour); i++ {
		res := <-Results
		if !verbose {
			if res.err == nil {
				fmt.Printf("Port %d open\n", res.port)
			}
		} else {
			if res.err == nil {
				log.Printf("Port %d open\n", res.port)
			} else {
				log.Printf("Port %d closed", res.port)
			}
		}
	}
}

// ring attempts to connect to the addr:port
func ring(r chan int, results chan result) {
	// TODO add sending result on channel
	for port := range r {
		addr := Target.addrs[0] + ":" + strconv.Itoa(port) // TODO FIX THIS!!!
		conn, err := d.Dial("tcp", addr)
		if err != nil {
			//TODO pick random ip in the list of ips, to test if its a single ip that is down. Very low-priority.
			results <- result{port, err}
		} else {
			results <- result{port, nil}
			conn.Close()
		}
	}
}

// Create Function to print both ports as they are scanned as done above; and to print a nice summary at the end
// Also make a cool ASCII banner when it starts
