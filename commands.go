package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
)

func cmdGet(cmd []string) (*memcache.Item, error) {
	key := cmd[1]
	if len(key) < 1 {
		return nil, errors.New("Key is required: get <key> [--verbose]")
	}
	it, err := mc.Get(key)
	if err != nil {
		return nil, errors.New("Error fetching data: " + err.Error())
	}
	return it, nil
}

func cmdDelete(cmd []string) error {
	key := cmd[1]
	if len(cmd) < 2 || len(key) < 1 {
		return errors.New("Key is required: delete <key>")
	}
	err := mc.Delete(key)
	if err != nil {
		return errors.New("Error fetching data: " + err.Error())
	}
	return nil
}

/* set, add, replace, append, prepend */
func cmdWriteData(cmd []string) error {
	var err error
	val := cmd[1:]
	if len(cmd[0]) < 1 || len(val) < 2 {
		return errors.New("Key is required: set|add|replace|append|prepend| <key> <value> [expiration]")
	}

	exp := 0
	if len(val) == 3 {
		exp, err = strconv.Atoi(val[2])
		if err != nil {
			exp = 0
		}
	}
	cleanVal := strings.Trim(val[1], "\"")
	pl := &memcache.Item{Key: val[0], Value: []byte(cleanVal), Expiration: int32(exp)}
	switch cmd[0] {
	case "set":
		err = mc.Set(pl)
		break
	case "add":
		err = mc.Add(pl)
		break
	case "replace":
		err = mc.Replace(pl)
		break
		/*
			case "append":
				err = mc.Append(pl)
				break
			case "prepend":
				err = mc.Prepend(pl)
				break
		*/
	default:
		log.Println(val[0])
	}
	if err != nil {
		errors.New("Error seting data: " + err.Error())
	}
	return nil
}

func cmdIncr(cmd []string) (uint64, error) {
	var err error
	val := cmd[1:]
	if len(cmd[0]) < 1 || len(val) < 2 {
		return 0, errors.New("Key is required: incr <key> <value>")
	}
	delta := 0
	if len(val) == 3 {
		delta, err = strconv.Atoi(val[1])
		if err != nil {
			delta = 0
		}
	}
	it, err := mc.Increment(val[0], uint64(delta))
	if err != nil {
		return 0, errors.New("Error fetching data: " + err.Error())
	}
	return it, nil
}

func cmdDecr(cmd []string) (uint64, error) {
	var err error
	val := cmd[1:]
	if len(cmd[0]) < 1 || len(val) < 2 {
		return 0, errors.New("Key is required: incr <key> <value>")
	}
	delta := 0
	if len(val) == 3 {
		delta, err = strconv.Atoi(val[1])
		if err != nil {
			delta = 0
		}
	}
	it, err := mc.Decrement(val[0], uint64(delta))
	if err != nil {
		return 0, errors.New("Error fetching data: " + err.Error())
	}
	return it, nil
}

func cmdFlushAll(cmd []string) error {
	err := mc.FlushAll()
	if err != nil {
		return errors.New("Error flushing data: " + err.Error())
	}
	return nil
}
