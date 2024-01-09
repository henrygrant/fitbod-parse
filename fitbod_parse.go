package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type FitbodRecord struct {
	Timestamp         string
	ExerciseName      string
	Reps              int
	Weight					  float64
	Duration          float64
	Distance          float64
	Incline           float64
	Resistance        float64
	IsWarmup          bool
	Note              string
	Multiplier        float64
}

func Unmarshal(reader *csv.Reader, v interface{}) error {
	record, err := reader.Read()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(v).Elem()
	if s.NumField() != len(record) {
		return &FieldMismatch{s.NumField(), len(record)}
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		switch f.Type().String() {
		case "string":
			f.SetString(record[i])
		case "bool":
			bval, err := strconv.ParseBool(record[i])
			if err != nil {
				return err
			}
			f.SetBool(bval)
		case "int":
			ival, err := strconv.ParseInt(strings.TrimSpace(record[i]), 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		case "float64":
			fval, err := strconv.ParseFloat(strings.TrimSpace(record[i]), 10)
			if err != nil {
				return err
			}
			f.SetFloat(fval)
		default:
			return &UnsupportedType{f.Type().String()}
		}
	}
	return nil
}

type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}


func main() {
	filePath := "./data/WorkoutExport.csv"

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		if(lineNum == 1) {
			continue
		}
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		reader.Comma = ','
		var record FitbodRecord
		for {
			err := Unmarshal(reader, &record)
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", record)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}