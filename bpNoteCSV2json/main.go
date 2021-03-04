package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type BpNoteType int

const (
	defaultMorningTime = "06:00"
	defaultNightTime   = "21:00"
)

const (
	colMorningBP = 2
	colNightBP   = 5
)

const (
	BloodPresure BpNoteType = iota
	Weight
)

type BpNote struct {
	bpNoteTyoe BpNoteType
	timeNanos  int64
	value      float32
}

func time2unix(yyyymmdd, hhmm string) (int64, error) {
	if hhmm == "" {
		hhmm = defaultMorningTime //TODO:
	}
	t, err := time.Parse("2006/01/0215:04", yyyymmdd+hhmm)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}

func readCsv(filepath string) ([]BpNote, error) {
	if filepath == "" {
		return nil, nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.New("Cannot open file: " + filepath)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string
	var lineNumber int

	for {
		line, err = reader.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(line)

		if lineNumber == 0 {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %v\n", err)
		}
		//if len(line) != 3 {
		//return nil, fmt.Errorf("csv format error: %v\n", line)
		//}

		date := line[0]
		switch {
		case line[colMorningBP] != "":
			fmt.Printf("MorningBP(%s,%s,%s)\n", line[colMorningBP-1], line[colMorningBP], line[colMorningBP+1])
			col := colMorningBP
			hhmm := line[col-1]
			unixtime, err := time2unix(date, hhmm)
			if err != nil {
				fmt.Printf("MorningBP(%s,%s,%s)\n", line[colMorningBP-1], line[colMorningBP], line[colMorningBP+1])
				return nil, fmt.Errorf("time format error: %v\n", err)
			}
			fmt.Printf("MorningBP(%s,%s,%s)\n", unixtime, line[colMorningBP], line[colMorningBP+1])
		}
		/*
			host := line[0]
			port_str := line[1]
			comment := line[2]
			port, err := strconv.Atoi(port_str)
			if err != nil {
				fmt.Println("port number error: " + port_str)
				return nil, fmt.Errorf("csv format error: %v\n", line)
			}
			serverList = append(serverList, server{
				host:    host,
				port:    port,
				comment: comment,
			})
		*/
	}
	return nil, nil
}

func main_() int {
	fmt.Println("main_()")
	if len(os.Args) < 2 {
		fmt.Printf("error len(os.Args):%v", len(os.Args))
		return 1
	}
	_, err := readCsv(os.Args[1])
	if err != nil {
		fmt.Printf("err : %s\n", err)
	}
	return 0
}

func main() {
	os.Exit(main_())
}
