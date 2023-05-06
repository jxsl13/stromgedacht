# stromgedacht

This is a Go client implementation for the https://www.stromgedacht.de/api-info/ api.


I am in no way or form associated with stromgedacht.de


## fetch te dependency
```shell
go get github.com/jxsl13/stromgedacht@latest
```


## example usage


```go
import (
    "github.com/jxsl13/stromgedacht/client"
    "encoding/json"
    "fmt"
)


func main() {
    c, _ := client.New()

    state, err := c.GetNow()
    if err != nil {
        panic(err)
    }

    data, _ := json.MarshalIndent(state, "", " ")
    fmt.Println(string(data))
}

```