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
		log.Panicf("Missing positional arg FILE. Try %s some_file.txt", os.Args[0])
	}
	basePath := os.Args[1]
	dbUrl := os.Getenv("DATABASE_URL")

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
