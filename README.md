# memcached-cli

	redis-cli inspired memcached client 
	history and verbose mode for get

## installing
	$ go install github.com/gleicon/memcached-cli

## build
	$ make

## running
	$ memcached-cli -s server:port

## Implemented commands
	
	get <key> [--verbose]
	set <key> <value> [expiration in secs]
	add <key> <value> [expiration in secs]
	replace <key> <value> [expiration in secs]
	append <key> <value> [expiration in secs]
	prepend <key> <value> [expiration in secs]
	incr <key> <delta>
	decr <key> <delta>
	delete <key>
	flush_all (asks for confirmation)
	quit


gleicon 2016
