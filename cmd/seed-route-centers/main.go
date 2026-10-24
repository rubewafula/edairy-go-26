// Command seed-route-centers makes sure every route has at least one center.
// Routes without a center get a default "<route name> Center". Safe to re-run.
//
//	go run ./cmd/seed-route-centers
//
// The API also runs this on startup, so this command is only needed for one-off runs.
package main

import (
	"fmt"
	"os"

	"github.com/rubewafula/edairy-go-26/internal/db"
	"github.com/rubewafula/edairy-go-26/internal/initializers"
	"github.com/rubewafula/edairy-go-26/internal/services"
)

func main() {
	initializers.LoadEnvVariables()
	db.ConnectToDatabase()

	n, err := services.EnsureRouteCenters()
	if err != nil {
		fmt.Fprintln(os.Stderr, "seed route centers:", err)
		os.Exit(1)
	}
	fmt.Printf("Seeded %d default route center(s). Every route now has at least one center.\n", n)
}
