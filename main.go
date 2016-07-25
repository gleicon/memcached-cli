package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/chzyer/readline"
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

var completer = readline.NewPrefixCompleter(
	/*
			readline.PcItem("stats",
			readline.PcItem(""),
			readline.PcItem("slabs"),
			readline.PcItem("malloc"),
			readline.PcItem("items"),
			readline.PcItem("detail"),
			readline.PcItem("sizes"),
			readline.PcItem("reset"),
		),
	*/
	readline.PcItem("get"),
	readline.PcItem("set"),
	readline.PcItem("add"),
	readline.PcItem("replace"),
	readline.PcItem("append"),
	readline.PcItem("prepend"),
	readline.PcItem("incr"),
	readline.PcItem("decr"),
	readline.PcItem("delete"),
	readline.PcItem("flush_all"),
	/*
		readline.PcItem("version"),
		readline.PcItem("verbosity"),
	*/
	readline.PcItem("quit"),
)

var mc *memcache.Client

func main() {

	memcachedServer := flag.String("s", "localhost:11211", "memcached server addr:port")
	flag.Usage = func() {
		fmt.Println("Usage: memcached-cli -s localhost:11211")
		os.Exit(1)
	}
	flag.Parse()

	prompt := fmt.Sprintf("%s\033[31m»\033[0m ", *memcachedServer)
	l, err := readline.NewEx(&readline.Config{
		Prompt:          prompt,
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	mc = memcache.New(*memcachedServer)

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		cmd := ParseQuotedArgs(line)
		if err != nil {
			println(err.Error())
		}
		ll := strings.ToLower(cmd[0])
		switch {
		case strings.HasPrefix(ll, "get"):
			var it *memcache.Item
			var err error
			if it, err = cmdGet(cmd); err != nil {
				println(err.Error())
				break
			}
			if cmd[len(cmd)-1] == "--verbose" {
				PrintMemcachedItem(it)
			} else {
				println(string(it.Value))
			}
		case strings.HasPrefix(ll, "set"):
			if err := cmdWriteData(cmd); err != nil {
				println(err.Error())
				break
			}
		case strings.HasPrefix(ll, "add"):
			if err := cmdWriteData(cmd); err != nil {
				println(err.Error())
				break
			}
		case strings.HasPrefix(ll, "replace"):
			if err := cmdWriteData(cmd); err != nil {
				println(err.Error())
				break
			}
		case strings.HasPrefix(ll, "incr"):
			var newval uint64
			var err error
			if newval, err = cmdIncr(cmd); err != nil {
				println(err.Error())
				break
			}
			println(newval)
		case strings.HasPrefix(ll, "decr"):
			var newval uint64
			var err error
			if newval, err = cmdDecr(cmd); err != nil {
				println(err.Error())
				break
			}
			println(newval)
		case strings.HasPrefix(ll, "delete"):
			if err := cmdDelete(cmd); err != nil {
				println(err.Error())
				break
			}
		case strings.HasPrefix(ll, "flush_all"):
			println("This will wipe all memcached content. Are you sure ? (yes/no)")
			var yn string
			fmt.Scanln(&yn)
			if err != nil {
				break
			}

			if strings.ToLower(yn) == "no" {
				println("aborted")
				break
			}
			if err := cmdFlushAll(cmd); err != nil {
				println(err.Error())
				break
			}
		case ll == "help":
			usage(l.Stderr())
		case ll == "quit":
			goto exit
		case line == "":
		default:
			log.Println("Command not found:", strconv.Quote(line))
		}
	}
exit:
}
