# cashier
Simple thread-safe in-memory TLRU caching library for Go

**If all you need is a thread-safe map for caching, we strongly encourage reviewing
[sync.Map](https://golang.org/pkg/sync/#Map) for your use case**

cashier is an in-memory key-value TLRU cache that is thread-safe and has the ability to store any object as the value.
There are feature toggles for both time-awareness and LRU functionality

## Installation

```
go get github.com/BradLugo/cashier
```

### Quick Start

```go
import (
	"fmt"
	"time"

    cashier
)

func main() {
    // Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c := cashier.New(cashier.NoItemLimit, 5*time.Minute, 10*time.Minute)

	// Set the value of the key "foo" to "bar", with the default expiration time
	c.Set("foo", "bar", cashier.DefaultExpiration)

	// Set the value of the key "baz" to 42, with no expiration time
	// (the item won't be removed until it is re-set, or removed using
	// c.Delete("baz")
	c.Set("baz", 42, cashier.NoExpiration)

	// Get the string associated with the key "foo" from the cache
	foo, found := c.Get("foo")
	if found {
		fmt.Println(foo)
	}

	// Since Go is statically typed, and cache values can be anything, type
	// assertion is needed when values are being passed to functions that don't
	// take arbitrary types, (i.e. interface{}). The simplest way to do this for
	// values which will only be used once--e.g. for passing to another
	// function--is:
	foo, found := c.Get("foo")
	if found {
		MyFunction(foo.(string))
	}

	// This gets tedious if the value is used several times in the same function.
	// You might do either of the following instead:
	if x, found := c.Get("foo"); found {
		foo := x.(string)
		// ...
	}
	// or
	var foo string
	if x, found := c.Get("foo"); found {
		foo = x.(string)
	}
	// ...
	// foo can then be passed around freely as a string

	// Want performance? Store pointers!
	c.Set("foo", &MyStruct, cache.DefaultExpiration)
	if x, found := c.Get("foo"); found {
		foo := x.(*MyStruct)
		// ...
	}
}
```

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the
[tags on this repository](https://github.com/BradLugo/cashier/tags). 

## License

This project is licensed under the Apache-2.0 License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

Huge inspiration for this project and much of the initial implementation came from patrickmn's
[go-cache](https://github.com/patrickmn/go-cache) and the [LRU library](https://github.com/golang/groupcache/tree/master/lru) from golang's [groupcache ](https://github.com/golang/groupcache)
