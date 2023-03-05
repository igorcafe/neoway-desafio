package main

import (
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	// PPROF=cpu ativa o profiling de CPU
	if os.Getenv("PPROF") == "cpu" {
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Panicf("failed to create cpu.prof: %v", err)
		}
		defer f.Close()
		log.Println("pprof: profiling CPU")
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(os.Args) != 2 {
	}
	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		log.Panicln("Missing env variable BASE_PATH")
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Panicln("Missing env variable DATABASE_URL")
	}

	runner, err := NewRunner(dbUrl, basePath)
	if err != nil {
		log.Panicf("failed to create Runner: %v", err)
	}
	defer runner.Close()

	err = runner.PrepareDatabase()
	if err != nil {
		log.Panic(err)
	}

	runner.ProcessLines()
}
