package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Log holds the default fields from an simple IIS Setup
type Log struct {
	Date          time.Time
	Time          string
	SIP           string
	CsMethod      string
	CsURIStem     string
	CsURIQuery    string
	SPort         string
	CsUsername    string
	CIP           string
	CsUserAgent   string
	CsReferer     string
	ScStatus      int
	ScSubstatus   int
	ScWin32Status int
	TimeTaken     int
	Customer      string
	Filename      string
}

func main() {
	files, err := ioutil.ReadDir("../sample")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var logs []Log
		fmt.Println(file.Name())
		filepath := "../sample/" + file.Name()
		if file, err := os.Open(filepath); err == nil {
			// make sure it gets closed
			defer file.Close()

			// create a new scanner and read the file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text() // read in the first line

				// Skip reading in comments, all lines which are not comments should
				// be a line in the log file
				if (strings.Contains(line, "#")) == false {
					data := strings.Fields(line)

					//convert strings to other types

					// LoadLocation uses http://golang.org/pkg/time/#LoadLocation.
					// icann names - get alist of these... we can use... for documentation
					// https://www.iana.org/time-zones
					location, err4 := time.LoadLocation("UTC")
					if err4 != nil {
						fmt.Printf("LoadLocation : %s", err4)
					}

					//https://www.iana.org/time-zones
					const shortFormlayout = "2006-01-02"
					iDate, err1 := time.ParseInLocation(shortFormlayout, data[0], location)

					if err1 != nil {
						log.Fatal(err1)
					}

					iscStatus, err2 := strconv.Atoi(data[11])
					iscSubstatus, err3 := strconv.Atoi(data[12])
					iscWin32Status, err4 := strconv.Atoi(data[13])
					itimeTaken, err5 := strconv.Atoi(data[14])

					if err2 != nil {
						log.Fatal(err2)
					}
					if err3 != nil {
						log.Fatal(err3)
					}
					if err4 != nil {
						log.Fatal(err4)
					}
					if err5 != nil {
						log.Fatal(err5)
					}
					// create a Log record
					row := Log{Date: iDate,
						Time:          data[1],
						SIP:           data[2],
						CsMethod:      data[3],
						CsURIStem:     data[4],
						CsURIQuery:    data[5],
						SPort:         data[6],
						CsUsername:    data[7],
						CIP:           data[8],
						CsUserAgent:   data[9],
						CsReferer:     data[10],
						ScStatus:      iscStatus,
						ScSubstatus:   iscSubstatus,
						ScWin32Status: iscWin32Status,
						TimeTaken:     itimeTaken,
						Customer:      "CUSTOMER_NAME",
						Filename:      file.Name(),
					}

					// Add the log record (row) to the logs slice
					logs = append(logs, row)
				}
			}

			// check for errors
			if err = scanner.Err(); err != nil {
				log.Fatal(err)
			}

			// Send data as tightly packed binary form via web service to cloud
			printJSON(logs[:1])

		} else {
			log.Fatal(err)
		}
	}
}

func printJSON(l []Log) {
	json, err := json.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print the Results
	fmt.Println(string(json))
}
