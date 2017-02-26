package engine

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DistributedSolutions/DIMWIT/common/constants"
	"github.com/DistributedSolutions/DIMWIT/constructor"
	"github.com/DistributedSolutions/DIMWIT/database"
	"github.com/DistributedSolutions/DIMWIT/factom-lite"
)

var _ = log.Prefix()

var CloseCalls []func()

func GrabFlagsAndRun() {
	var (
		fct  = flag.String("fct", "fake", "Factom Client Type: 'fake', 'dumb")
		cdbt = flag.String("lvl2", "Map", "Constructor DB Type: 'Map', 'Bolt', 'LDB'")
	)
	flag.Parse()

	StartEngine(*fct, *cdbt)
}

// StartEngine is the main start, that launches the appropriate go routines and handles closing.
func StartEngine(factomClientType string, lvl2CacheType string) error {
	log.Println("-- DIMWIT Engine Initiated -- ")
	log.Printf("%-20s: %s\n", "FactomClientType", factomClientType)
	log.Printf("%-20s: %s\n", "Level2CacheType", lvl2CacheType)

	CloseCalls = make([]func(), 0)
	// Factom-Lite Client
	var factomClient lite.FactomLite
	switch factomClientType {
	case "dumb":
		factomClient = lite.NewDumbLite()
	case "fake":
		factomClient = lite.NewFakeDumbLite()
	default:
		return fmt.Errorf("Level 2 Cache Type given not valid. Found '%s', expected either: 'dumb', 'fake'", factomClientType)
	}

	var lvl2Cache database.IDatabase
	switch lvl2CacheType {
	case "Bolt":
		lvl2Cache = database.NewBoltDB(constants.HIDDEN_DIR + constants.LVL2_CACHE)
	case "LDB":
	case "Map":
		lvl2Cache = database.NewMapDB()
	default:
		return fmt.Errorf("DBType given not valid. Found '%s', expected either: Bolt, Map, LDB", lvl2CacheType)
	}

	// Construtor -> Updates level 2 cache
	con, err := constructor.NewContructor(lvl2Cache)
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

	w := new(WholeState)
	w.Constructor = con
	w.FactomClient = factomClient

	// Run the Control
	Control(w)
	return nil
}
