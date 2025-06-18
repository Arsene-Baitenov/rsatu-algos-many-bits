package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Arsene-Baitenov/rsatu-algos-many-bits/engine"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Person string

const (
	P Person = "P"
	S Person = "S"
)

type Speech string

const (
	Know     Speech = "знаю"
	DontKnow Speech = "не знаю"
	Stop     Speech = "стоп"
)

func main() {
	logLevelFlag := flag.String("log-level", "disabled", "log level: trace, debug, info, warn, error, fatal, panic")
	flag.Parse()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	logLevel, err := zerolog.ParseLevel(strings.ToLower(*logLevelFlag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level: %s\n", *logLevelFlag)
		os.Exit(1)
	}
	zerolog.SetGlobalLevel(logLevel)

	var n uint64
	fmt.Print("Введите n: ")
	fmt.Fscan(os.Stdin, &n)

	eng := engine.New(n)

	currPerson := P
	var speech Speech
	reader := bufio.NewReader(os.Stdin)

loop:
	for {
		fmt.Printf("%v: ", currPerson)
		input, _ := reader.ReadString('\n')
		speech = Speech(strings.TrimSpace(input))

		switch speech {
		case Know:
			if currPerson == P {
				fmt.Println(eng.ComputePairsByProds())
				currPerson = S
			} else {
				fmt.Println(eng.ComputePairsBySums())
				currPerson = P
			}
		case DontKnow:
			if currPerson == P {
				eng.FilterNonTrivialProds()
				currPerson = S
			} else {
				eng.FilterNonTrivialSums()
				currPerson = P
			}
		case Stop:
			break loop
		default:
			fmt.Println("Некорректная реплика")
		}
	}
}
