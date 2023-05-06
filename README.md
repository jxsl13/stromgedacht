# stromgedacht

This is a Go client implementation for the https://www.stromgedacht.de/api-info/ api.


I am in no way or form associated with stromgedacht.de


## fetch the dependency
```shell
go get github.com/jxsl13/stromgedacht@latest
```


## example usage


```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/jxsl13/stromgedacht/client"
)

func main() {
	c, _ := client.New()

	state, err := c.GetNow("68309") // Mannheim
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(state, "", " ")
	fmt.Println(string(data))
}

```