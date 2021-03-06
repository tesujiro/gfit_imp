package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
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
	colWeight    = 27
)

const (
	BloodPressure BpNoteType = iota
	Weight
)

type BpNote struct {
	bpNoteType BpNoteType
	timeNanos  int64
	value      float32
	highValue  float32
	lowValue   float32
}

func time2unix(yyyymmdd, hhmm, defaultHhmm string) (int64, error) {
	if hhmm == "" {
		hhmm = defaultHhmm
	}
	t, err := time.Parse("2006/1/2 15:04", yyyymmdd+" "+hhmm)
	if err != nil {
		return 0, err
	}
	return t.UnixNano(), nil
}

func parseBloodPressure(date, hhmm, highStr, lowStr, defaultHhmm string) (BpNote, error) {
	unixtime, err := time2unix(date, hhmm, defaultHhmm)
	if err != nil {
		return BpNote{}, fmt.Errorf("time format error: %v\n", err)
	}
	//fmt.Printf("MorningBP(%v,%v,%v)\n", unixtime, line[colMorningBP], line[colMorningBP+1])

	high, err := strconv.ParseFloat(highStr, 32)
	if err != nil {
		return BpNote{}, fmt.Errorf("parseFloat error: %s", highStr)
	}
	low, err := strconv.ParseFloat(lowStr, 32)
	if err != nil {
		return BpNote{}, fmt.Errorf("parseFloat error: %s", lowStr)
	}

	bp := BpNote{
		bpNoteType: BloodPressure,
		timeNanos:  unixtime,
		highValue:  float32(high),
		lowValue:   float32(low),
	}
	return bp, nil
}

func parseLine(line []string) ([]BpNote, error) {
	bpn := []BpNote{}
	date := line[0]
	if line[colMorningBP] != "" {
		col := colMorningBP
		hhmm := line[col-1]
		highStr := line[col]
		lowStr := line[col+1]

		bp, err := parseBloodPressure(date, hhmm, highStr, lowStr, defaultMorningTime)
		if err != nil {
			return bpn, err
		}
		bpn = append(bpn, bp)
	}
	if line[colNightBP] != "" {
		col := colNightBP
		hhmm := line[col-1]
		highStr := line[col]
		lowStr := line[col+1]

		bp, err := parseBloodPressure(date, hhmm, highStr, lowStr, defaultNightTime)
		if err != nil {
			return bpn, err
		}
		bpn = append(bpn, bp)
	}
	if line[colWeight] != "" {
		unixtime, err := time2unix(date, "", defaultMorningTime)
		if err != nil {
			return bpn, fmt.Errorf("time format error: %v\n", err)
		}
		weight, err := strconv.ParseFloat(line[colWeight], 32)
		if err != nil {
			return bpn, fmt.Errorf("parseFloat error: %s", line[colWeight])
		}
		bp := BpNote{
			bpNoteType: Weight,
			timeNanos:  unixtime,
			value:      float32(weight),
		}
		bpn = append(bpn, bp)
	}
	return bpn, nil
}

func csv2json(filepath string) error {
	if filepath == "" {
		return nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return errors.New("Cannot open file: " + filepath)
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
		if err != nil {
			return fmt.Errorf("csv read error: %v\n", err)
		}
		//fmt.Println(line)

		/*
			for i, col := range line {
				fmt.Printf("line[%d]\t:%v\n", i, col)
			}
		*/

		lineNumber++
		if lineNumber == 1 {
			continue
		}

		//if len(line) != 3 {
		//return nil, fmt.Errorf("csv format error: %v\n", line)
		//}

		bpns, err := parseLine(line)
		if err != nil {
			return err
		}
		_ = bpns
		/*
			for _, bpn := range bpns {
				fmt.Println(bpn)
			}
		*/
	}
	return nil
}

func main_() int {
	if len(os.Args) < 2 {
		fmt.Printf("error len(os.Args):%v", len(os.Args))
		return 1
	}

	jsMap := getBpRequest()
	//fmt.Println(jsMap)
	js, err := json.Marshal(jsMap)
	if err != nil {
		fmt.Printf("err : %s\n", err)
	}
	fmt.Printf("%s\n", js)

	err = csv2json(os.Args[1])
	if err != nil {
		fmt.Printf("err : %s\n", err)
	}
	return 0
}

func main() {
	os.Exit(main_())
}
