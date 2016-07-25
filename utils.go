package main

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/bradfitz/gomemcache/memcache"
)

func PrintMemcachedItem(it *memcache.Item) {
	fmt.Printf("Key: %s\nValue: %s (%v)\nFlags: %v\nExpiration: %d\n", it.Key, string(it.Value), it.Value, it.Flags, it.Expiration)
}

func ParseQuotedArgs(line string) []string {
	// inspipasted from https://play.golang.org/p/ztqfYiPSlv
	lastQuote := rune(0)
	f := func(c rune) bool {
		switch {
		case c == lastQuote:
			lastQuote = rune(0)
			return false
		case lastQuote != rune(0):
			return false
		case unicode.In(c, unicode.Quotation_Mark):
			lastQuote = c
			return false
		default:
			return unicode.IsSpace(c)

		}
	}

	m := strings.FieldsFunc(line, f)
	return m
}
