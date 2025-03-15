## Installation

As a library

```shell
go get github.com/faizal-asep-outlook/env
```

## Usage

create a struct type with default values ​​to retrieve and store data from the server environment variables:

```go
type Config struct {
	Port          string `env:"APP_PORT" default:"8080"`
}
```

Then in your Go app you can do something like

```go
package main

import (
    "log"
    "github.com/faizal-asep-outlook/env"
)

func main() {
	cfg := Config{}
	err := env.Parse(&cfg)
  if err != nil {
		log.Fatal("Error parse config")
	}
  log.Println(cfg.Port)

  // now do something 
}
```
