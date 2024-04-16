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
    "time"

    "github.com/jxsl13/stromgedacht/client"
)

func main() {
    c, _ := client.New()

    // get current state
    state, err := c.GetNow("68309") // Mannheim
    if err != nil {
        panic(err)
    }

    stateStr, _ := json.MarshalIndent(state, "", " ")
    fmt.Println(string(stateStr))  // { "state": -1 }

    // get forecast
    forecast, err := c.GetForecast("68309", time.Time{}, time.Time{}) // Mannheim, default time range
    if err != nil {
        panic(err)
    }

    loadStr, _ := json.MarshalIndent(forecast.Load[0], "", " ")
    residualLoadStr, _ := json.MarshalIndent(forecast.ResidualLoad[0], "", " ")
    renewableEnergyStr, _ := json.MarshalIndent(forecast.RenewableEnergy[0], "", " ")
    superGreenThresholdStr, _ := json.MarshalIndent(forecast.SuperGreenThreshold[0], "", " ")

    fmt.Println(string(loadStr))  // { "dateTime": "2024-04-14T00:00:00Z", "value": 4207 }
    fmt.Println(string(residualLoadStr))  // { "dateTime": "2024-04-14T00:00:00Z", "value": 2511.419 }
    fmt.Println(string(renewableEnergyStr))// { "dateTime": "2024-04-14T00:00:00Z", "value": 1695.581 }
    fmt.Println(string(superGreenThresholdStr))// { "dateTime": "2024-04-14T00:00:00Z", "value": 4000 }
}

```