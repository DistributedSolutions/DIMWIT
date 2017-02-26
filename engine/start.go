package engine

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var _ = log.Prefix()

var CloseCalls []func()

func GrabFlagsAndRun() {
	var (
		fct  = flag.String("fct", "fake", "Factom Client Type: 'fake', 'dumb")
		cdbt = flag.String("condb", "Map", "Constructor DB Type: 'Map', 'Bolt', 'LDB'")
	)
	flag.Parse()

	StartEngine(*fct, *cdbt)
}

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
	/*c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Println("Safely Closing....")
		for _, f := range CloseCalls {
			f()
		}
		log.Println("Completed safe close")
		os.Exit(1)
	}()*/
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Safely Closing....")
		for _, f := range CloseCalls {
			f()
		}
		log.Println("Completed safe close")
		os.Exit(1)
	}()

	// Run the Control
	Control()
	return nil
}
