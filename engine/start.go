package engine

import (
	"os"
	"os/signal"

	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var CloseCalls []func()

// StartEngine is the main start, that launches the appropriate go routines and handles closing.
func StartEngine(factomClientType string, constructorDBType string) error {
	CloseCalls = make([]func(), 0)
	// Factom-Lite Client
	var factomClient lite.FactomLite
	switch factomClientType {
	case "dumb":
		factomClient = lite.NewDumbLite()
	case "fake":
		factomClient = lite.NewFakeDumbLite()
	}

	// Construtor -> Updates level 2 cache
	con, err := constructor.NewContructor(constructorDBType)
	if err != nil {
		return err
	}
	con.SetReader(factomClient)
	CloseCalls = append(CloseCalls, con.InterruptClose)

	// Start Go Routines
	go con.StartConstructor()

	// Safe Close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _, f := range CloseCalls {
			f()
		}
	}()

	return nil
}
