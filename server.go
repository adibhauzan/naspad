package naspad

import (
	"log"
	"net/http"
)

func NewServer() *Driver {
	driver := NewDriver()

	// Set up a default RGroup with the base route path as "/"
	driver.RGroup.baseRoutePath = "/"

	driver.RGroup.BaseRoutePath()
	// Return the Driver which implements http.Handler
	return driver
}

// Run starts the HTTP server with the provided Driver
func (d *Driver) Run(addr string) {
	log.Printf("Listening and serving HTTP on %s", addr)
	if err := http.ListenAndServe(addr, d); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}