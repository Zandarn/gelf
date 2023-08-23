# Go client for GELF protocol

Only UDP.

# Example
```go
package main

import (
	"fmt"
	"log"

	"github.com/Zandarn/gelf"
)

func main() {
	client := gelf.GelfClient()
	client.Config.Endpoint = "192.168.1.1:12201"

	gelfMessage := gelf.Message{
		Message:  "Example message",
		Level:    gelf.LogLevel.WARNING,
		Facility: gelf.LogFacility.LOCAL0,
		Extra: map[string]string{
			"Extra-key":   "extra_value",
			"Another key": "another value",
		},
	}

	length, err := client.SendMessage(gelfMessage)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	fmt.Println(length)
}

```
