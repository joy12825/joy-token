# redis
GoFrame `gredis.Adapter` implements using `go-redis`.

# Installation
```
go get -u github.com/joy12825/gf/contrib/nosql/redis/v2
```

Commonly imported at top of `main.go`:
```go
package main

import (
	_ "github.com/joy12825/gf/contrib/nosql/redis/v2"

	// Other imported packages.
)

func main() {
	// Main logics.
}
```
