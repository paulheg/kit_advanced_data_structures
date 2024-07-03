package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"

	"github.com/paulheg/kit_advanced_data_structures/internal/bitvector"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var verbose = flag.Bool("verbose", false, "add more parameters to the output")

func main() {

	flag.Parse()

	// if we want to record a CPU profile an output filepath will be set
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// parse the commandline args these are the input and output paths
	files := flag.Args()
	if len(files) != 2 {
		log.Fatal("wrong number of arguments, need [input] [output]")
	}

	inputPath := files[0]
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatal("could not open input file:", err)
	}
	defer inputFile.Close()

	outputPath := files[1]
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatal("could not create output file:", err)
	}
	defer outputFile.Close()

	// here the actual processing begins
	err = bitvector.ProcessFile(inputFile, outputFile, os.Stdout, *verbose)
	if err != nil {
		log.Fatal("error processing file:", err)
	}
}
