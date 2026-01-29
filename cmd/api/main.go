package main

import (
	"fmt"

	"github.com/thebeyond-net/control-plane/internal/config"
)

func main() {
	if _, err := config.New(); err != nil {
		panic(err)
	}

	fmt.Println("Infrastructure Ready")
}
