package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Mode string

const (
	charMode Mode = "char"
	lineMode Mode = "line"
)

func main() {
	flag.Usage = func() {
		fmt.Println(`go-wc: the wc commandline method but in go`)
		flag.PrintDefaults()
	}
	// set up the args for the commands, same as wc
	var lineC, charC bool
	flag.BoolVar(&lineC, "l", false, "print the newline counts")
	flag.BoolVar(&charC, "c", false, "print the character counts")
	flag.Parse()

	// check if the input is piped in or not
	readStat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	mode := lineMode
	if charC {
		mode = charMode
	}
	var reader io.Reader
	if readStat.Mode()&os.ModeCharDevice == 0 {
		// read from the standard input and read using a buffer for better efficiency.
		reader = bufio.NewReader(os.Stdin)
	} else {
		// fetch the filepath from the args and read the file.
		fileData, readErr := os.ReadFile(flag.Arg(0))
		if readErr != nil {
			if errors.Is(readErr, os.ErrNotExist) {
				_, _ = os.Stdout.Write([]byte("file does not exist \n"))
				os.Exit(0)
			}
			panic(readErr)
		}
		reader = strings.NewReader(string(fileData))
	}
	countFromBuffer(reader, mode)
}

// countFromBuffer counts the number of lines/characters in the input based on the mode
func countFromBuffer(r io.Reader, mode Mode) {
	// create a buffer
	st := make([]byte, 32*1024)
	count := 0
	var err error
	for {
		byteIndex := 0
		byteIndex, err = r.Read(st)
		if err != nil {
			break
		}
		if mode == lineMode {
			count = count + bytes.Count(st[:byteIndex], []byte("\n"))
		} else {
			count = count + len(st[:byteIndex])
		}
	}
	if err != nil {
		if err == io.EOF {
			fmt.Println(count)
			os.Exit(0)
		}
		os.Exit(2)
	}
}
