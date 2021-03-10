<p align="right">
  ⭐ &nbsp;&nbsp;<strong>the project to show your appreciation.</strong> :arrow_upper_right:
</p>

<p align="right">
  <a href="http://godoc.org/github.com/rocketlaunchr/remember-go"><img src="http://godoc.org/github.com/rocketlaunchr/remember-go?status.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/rocketlaunchr/remember-go"><img src="https://goreportcard.com/badge/github.com/rocketlaunchr/remember-go" /></a>
  <a href="https://gocover.io/github.com/rocketlaunchr/remember-go"><img src="http://gocover.io/_badge/github.com/rocketlaunchr/remember-go" /></a>
</p>

<p align="center">
<img src="https://github.com/rocketlaunchr/remember-go/raw/master/assets/logo.png" alt="dataframe-go" />
</p>

# Cache Slow Database Queries

This package is used to cache the results of slow database queries in memory or Redis.
It can be used to cache any form of data (i.e it's not limited to caching database results). A Redis and in-memory storage driver is provided.

See [Article](https://medium.com/@rocketlaunchr.cloud/caching-slow-database-queries-1085d308a0c9) for further details including a tutorial.

The package is **production ready** and the API is stable. A variant of this package has been used in production for over 4 years.

## Installation

```
go get -u github.com/rocketlaunchr/remember-go
```

## Create a Key

Let’s assume the query’s argument is an arbitrary `search` term and a `page` number for pagination.

### CreateKeyStruct

CreateKeyStruct can generate a JSON based key by providing a struct.

```go
type Key struct {
    Search string
    Page   int `json:"page"`
}

var key string = remember.CreateKeyStruct(Key{"golang", 2})
```

### CreateKey

CreateKey provides more flexibility to generate keys:

```go
// Key will be "search-golang-2"
key :=  remember.CreateKey(false, "-", "search-x-y", "search", "golang", 2)
```

## Initialize the Storage Driver

### In-Memory

```go
import "github.com/rocketlaunchr/remember-go/memory"

var ms = memory.NewMemoryStore(10 * time.Minute)
```

### Redis

The Redis storage driver relies on Gary Burd’s excellent [Redis client library](https://github.com/gomodule/redigo/).

```go
import red "github.com/rocketlaunchr/remember-go/redis"
import "github.com/gomodule/redigo/redis"

var rs = red.NewRedisStore(&redis.Pool{
    Dial: func() (redis.Conn, error) {
        return redis.Dial("tcp", "localhost:6379")
    },
})
```

### Memcached

An experimental (and untested) memcached driver is provided.
It relies on Brad Fitzpatrick's [memcache driver](https://godoc.org/github.com/bradfitz/gomemcache/memcache).

### Ristretto

DGraph's [Ristretto](https://github.com/dgraph-io/ristretto) is a fast, fixed size, in-memory cache with a dual focus on throughput and hit ratio performance.

### Nocache

This driver is for testing purposes. It does not cache any data.

## Create a SlowRetrieve Function

The package initially checks if data exists in the cache. If it doesn’t, then it elegantly fetches the data directly from the database by calling the `SlowRetrieve` function. It then saves the data into the cache so that next time it doesn’t have to refetch it from the database.

```go
type Result struct {
    Title string
}

slowQuery := func(ctx context.Context) (interface{}, error) {
    results := []Result{}

    stmt := `
        SELECT title
        FROM books WHERE title LIKE ?
        ORDER BY title LIMIT ?, 20
    `

    rows, _ := db.QueryContext(ctx, stmt, search, (page-1)*20)

    for rows.Next() {
        var title string
        rows.Scan(&title)
        results = append(results, Result{title})
    }

    return results, nil
}
```

## Usage

```go
key := remember.CreateKeyStruct(Key{"golang", 2})
exp := 10*time.Minute

results, found, err := remember.Cache(ctx, ms, key, exp, slowQuery, remember.Options{GobRegister: false})

return results.([]Result) // Type assert in order to use
```

## Gob Register Errors

The Redis storage driver stores the data in a `gob` encoded form. You have to register with the [`gob`](https://golang.org/pkg/encoding/gob/) package the data type returned by the `SlowRetrieve` function. It can be done inside a `func init()`. Alternatively, you can set the `GobRegister` option to true. This will slightly impact concurrency performance however.

## Other useful packages

- [awesome-svelte](https://github.com/rocketlaunchr/awesome-svelte) - Resources for killing react
- [dataframe-go](https://github.com/rocketlaunchr/dataframe-go) - For statistics, machine-learning, and data manipulation/exploration
- [dbq](https://github.com/rocketlaunchr/dbq) - Zero boilerplate database operations for Go
- [electron-alert](https://github.com/rocketlaunchr/electron-alert) - SweetAlert2 for Electron Applications
- [google-search](https://github.com/rocketlaunchr/google-search) - Scrape google search results
- [igo](https://github.com/rocketlaunchr/igo) - A Go transpiler with cool new syntax such as fordefer (defer for for-loops)
- [mysql-go](https://github.com/rocketlaunchr/mysql-go) - Properly cancel slow MySQL queries
- [react](https://github.com/rocketlaunchr/react) - Build front end applications using Go
- [testing-go](https://github.com/rocketlaunchr/testing-go) - Testing framework for unit testing

#

### Legal Information

The license is a modified MIT license. Refer to `LICENSE` file for more details.

**© 2019-21 PJ Engineering and Business Solutions Pty. Ltd.**

### Final Notes

Feel free to enhance features by issuing pull-requests.
