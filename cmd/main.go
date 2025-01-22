// cmd/main.go
package main

import (
	"github.com/Tushar7890/RetailPulse/internal/api"
)

func main() {
	router := api.SetupRouter()
	router.Run(":8080")
}
