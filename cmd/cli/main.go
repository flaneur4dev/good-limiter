package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const invalidArgs = "invalid arguments"

var (
	add    = flag.String("add", "", "add subnet (CIDR notation) to black/white list in the format of subnet=list")
	del    = flag.String("del", "", "delete subnet (CIDR notation) from black/white list in the format of subnet=list")
	drop   = flag.String("drop", "", "drop buckets in the format of login@ip")
	server = flag.String("s", ":50051", "server address in the format of host:port")
	help   = flag.Bool("h", false, "print help information")
)

func printHelp() {
	txt := `Utility for rate limiter configuration.
	Usage: limiter-cli [-add] [-del] [-drop] [-h] [-s]
	Example: limiter-cli -add 12.168.0.0/24=black -del 10.18.0.64/26=white -drop login_42@92.168.0.101 -s localhost:50050`

	fmt.Println(txt)
}

type result struct {
	sl []string
}

func (r *result) Add(s string) {
	r.sl = append(r.sl, s)
}

func (r *result) String() string {
	if r.sl == nil {
		return "type -h to get help information"
	}

	return strings.Join(r.sl, "\n")
}

func main() {
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	res := new(result)

	c, err := newClient(*server)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1) //nolint:gocritic
	}
	defer c.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if *add != "" {
		s := "[add]: "
		t := strings.Split(*add, "=")
		if len(t) != 2 {
			s += invalidArgs
		} else {
			s += c.AddNet(ctx, t[0], t[1])
		}

		res.Add(s)
	}

	if *del != "" {
		s := "[del]: "
		t := strings.Split(*del, "=")
		if len(t) != 2 {
			s += invalidArgs
		} else {
			s += c.DeleteNet(ctx, t[0], t[1])
		}

		res.Add(s)
	}

	if *drop != "" {
		s := "[drop]: "
		t := strings.Split(*drop, "@")
		if len(t) != 2 {
			s += invalidArgs
		} else {
			s += c.DropBucket(ctx, t[0], t[1])
		}

		res.Add(s)
	}

	fmt.Println(res)
}
