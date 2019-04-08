Caching Slow Database Queries [![GoDoc](http://godoc.org/github.com/rocketlaunchr/remember-go?status.svg)](http://godoc.org/github.com/rocketlaunchr/remember-go) [![Go Report Card](https://goreportcard.com/badge/github.com/rocketlaunchr/remember-go)](https://goreportcard.com/report/github.com/rocketlaunchr/remember-go)
===============

This package is used to cache the results of slow database queries in memory or Redis.
It can be used to cache any form of data. A Redis and in-memory storage driver is provided.

See [Article](https://medium.com/@rocketlaunchr.cloud/caching-slow-database-queries-1085d308a0c9) for further details.

The package is **production ready** and the API is stable. A variant of this package has been used in production for over 1.5 years.
Once the community creates more storage drivers, version 1.0.0 will be tagged. It is recommended your package manager locks to a commit id instead of the master branch directly.


## Installation

```
go get -u github.com/rocketlaunchr/remember-go
```


## QuickStart

```go
import (
	"github.com/gomodule/redigo/redis"
	"github.com/rocketlaunchr/remember-go"
	"github.com/rocketlaunchr/remember-go/memory"
	red "github.com/rocketlaunchr/remember-go/redis"
)

var ms = memory.NewMemoryStore(10 * time.Minute) // In-memory cache
var rs = red.NewRedisStore(&redis.Pool{
	Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", "localhost:6379")
	},
})

type Key struct {
	Page int
	Loop bool `json:"loop"`
}



slowQuery := func(ctx context.Context) (interface{}, error) {
	type data struct {
		Str string
		Int int
	}

	return data{"asdgasdg", 5}, nil
}


key := remember.CreateKeyStruct(&K{5, true})

out, found, err := remember.Cache(context.Background(), rs, key, 10*time.Minute, slowQuery, remember.Options{GobRegister: true})

fmt.Println(out, found, err)

```

#

### Legal Information

The license is a modified MIT license. Refer to `LICENSE` file for more details.

**© 2018 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests.

**Star** the project to show your appreciation.