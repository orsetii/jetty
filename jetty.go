package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	URL     string
	err     error
	d       net.Dialer
	Threads int
	verbose bool
	River   chan int    = make(chan int, 1000)
	Results chan result = make(chan result, 65535)
	timeout int
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

// Struct holding key data about the target which the program is pointed at via flags.
var Target struct {
	harbour Ports
	addrs   []string
	Total   struct {
		openPorts   Ports
		closedPorts Ports
		ipsUsed     []string
		timeTaken   time.Duration
	}
}

type result struct {
	port int
	err  error
}

// TODO add context package usage maybe?

func init() {
	flag.StringVar(&URL, "u", "", "The URL you want to scan.") //TODO make this into a list flag.
	flag.IntVar(&Threads, "t", 10, "How many threads you want to scan. Defaults to 20, I would recommend around 40 if you have a good internet connection. ")
	flag.Var(&Target.harbour, "p", "Enter ports to scan. Must be done in format `x-y` (to be fixed.)")
	flag.BoolVar(&verbose, "v", false, "Enable verbose output.") // Maybe change this so dont have to pass it to functions
	flag.IntVar(&timeout, "timeout", 350, "Time to wait for response from target.")
}

func main() {
	start := time.Now()
	flag.Parse()
	fmt.Printf(`

    ___  _______  _________  _________    ___    ___ 
   |\  \|\  ___ \|\___   ___\\___   ___\ |\  \  /  /|
   \ \  \ \   __/\|___ \  \_\|___ \  \_| \ \  \/  / /
 __ \ \  \ \  \_|/__  \ \  \     \ \  \   \ \    / / 
|\  \\_\  \ \  \_|\ \  \ \  \     \ \  \   \/  /  /  
\ \________\ \_______\  \ \__\     \ \__\__/  / /    
 \|________|\|_______|   \|__|      \|__|\___/ /     
                                        \|___|/      
                                                     
                                                     

`)
	fmt.Printf("Resolving URL...\n\n")

	d.Timeout = time.Millisecond * time.Duration(timeout) // Wait 500ms for response, if not give up.
	Target.addrs = resolve([]string{URL})
	if len(Target.addrs) < 1 {
		log.Fatalf("\nCould not resolve URL to any addresses.\nPlease try again.\n")
	}
	for _, ip := range Target.addrs {
		fmt.Printf("IP found: %s\n", ip)
	}
	// River channel to send ports on for go routines to pick up when they can.
	defer close(River)
	defer close(Results)
	fmt.Printf("\nScanning...\n\n")
	for i := 0; i < Threads; i++ {
		go ring(River, Results)
	}
	for i := 0; i < len(Target.harbour); i++ {
		res := <-Results
		if res.err == nil {
			Target.Total.openPorts = append(Target.Total.openPorts, res.port)
			if verbose {
				log.Printf("Port %d open", res.port)
			} else {
				fmt.Printf("Port %d open", res.port)
			}
		} else {
			Target.Total.closedPorts = append(Target.Total.closedPorts, res.port)
			if verbose {
				log.Printf("Port %d closed", res.port)
			}
		}
	}
	Target.Total.timeTaken = time.Since(start)
	scorecard()
}

// ring attempts to connect to the addr:port
func ring(r chan int, results chan result) {
	for port := range r {
		addr := Target.addrs[0] + ":" + strconv.Itoa(port)
		conn, err := d.Dial("tcp", addr)
		if err != nil {
			results <- result{port, err}
		} else {
			results <- result{port, nil}
			conn.Close()
		}
	}
}

func resolve(urls []string) (Ips []string) { // TODO ADD IPS USED TO TARGET.TOTAL STRUCT AS AND WHEN INSIDE THIS FUNC
	for _, u := range urls {
		uResolved, err := net.LookupHost(u)
		if err != nil && verbose {
			log.Printf("Error occurred while trying to resolve URL: %s\n Error: %s\n", u, err)
			continue
		}
		Ips = append(Ips, uResolved...)
	}
	return Ips
}

// Prints Results in appropiate fashion.
func scorecard() {
	fmt.Printf("-----------------\nTime Taken %.2f seconds \n%d Ports scanned\n%d Ports Open:\n", Target.Total.timeTaken.Seconds(), len(Target.harbour), len(Target.Total.openPorts))
	for _, open := range Target.Total.openPorts {
		fmt.Printf("	Port %d open\n", open)
	}
}
