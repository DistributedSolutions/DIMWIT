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
	"github.com/DistributedSolutions/DIMWIT/provider"
	"github.com/DistributedSolutions/DIMWIT/torrent"
	"github.com/DistributedSolutions/DIMWIT/util"
	"github.com/DistributedSolutions/DIMWIT/writeHelper"
)

var _ = log.Prefix()

var CloseCalls []func()

func GrabFlagsAndRun() {
	var (
		fct             = flag.String("fct", "fake", "Factom Client Type: 'fake', 'dumb")
		cdbt            = flag.String("lvl2", "Map", "Constructor DB Type: 'Map', 'Bolt', 'LDB'")
		persistFactomDB = flag.Bool("p", false, "Use persistent factom storage")
	)
	flag.Parse()

	err := StartEngine(*fct, *cdbt, *persistFactomDB)
	if err != nil {
		log.Printf("Error: Failed to start: %s", err.Error())
	}
}

// StartEngine is the main start, that launches the appropriate go routines and handles closing.
func StartEngine(factomClientType string, lvl2CacheType string, persistent bool) error {
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
		if persistent {
			factomClient = lite.NewBoltFakeDumbLite()
		} else {
			factomClient = lite.NewMapFakeDumbLite()
		}
	default:
		return fmt.Errorf("Level 2 Cache Type given not valid. Found '%s', expected either: 'dumb', 'fake'", factomClientType)
	}

	var lvl2Cache database.IDatabase
	switch lvl2CacheType {
	case "Bolt":
		lvl2Cache = database.NewBoltDB(util.GetHomeDir() + constants.HIDDEN_DIR + constants.LVL2_CACHE)
	case "LDB":
	case "Map":
		lvl2Cache = database.NewMapDB()
	default:
		return fmt.Errorf("DBType given not valid. Found '%s', expected either: Bolt, Map, LDB", lvl2CacheType)
	}

	// Starts SQL InterfaceDB
	sql, err := constructor.NewSqlWriter()
	if err != nil {
		return err
	}

	// Construtor -> Updates level 2 cache
	con, err := constructor.NewContructor(lvl2Cache, sql)
	if err != nil {
		return err
	}
	con.SetReader(factomClient)
	CloseCalls = append(CloseCalls, con.InterruptClose)

	// Write Helper
	wh, err := writeHelper.NewWriterHelper(con, factomClient)
	if err != nil {
		return nil
	}

	// Provider -> Serves API
	prov, err := provider.NewProvider(lvl2Cache, wh, factomClient)
	if err != nil {
		return err
	}
	CloseCalls = append(CloseCalls, prov.Close)

	// Torrent Client
	torClient, err := torrent.NewTorrentClient()
	if err != nil {
		return err
	}
	CloseCalls = append(CloseCalls, torClient.Close)
	prov.Router.HandleFunc("/stream", torClient.HandleStream)

	//sets interface between api and torrent client up
	prov.TorrentClientInterface.SetClient(torClient)

	// Start Go Routines
	go con.StartConstructor()
	go prov.Serve()

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

	// Set WholeState
	w := new(WholeState)
	w.Constructor = con
	w.FactomClient = factomClient
	w.Provider = prov
	w.TorrentClient = torClient
	w.WriteHelper = wh

	// Run the Control
	Control(w)
	return nil
}
