// Copyright craig nicholson. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Header contains the summary information about the logs
type Header struct {
	Logs         []Log  // The logs from IIS
	CustomerName string // The unique customers name for grouping and viewing.
	Filename     string // The file which contained the logs.
	Records      int    // The number of rows in the log file we just parsed.
}

// Log contains the default fields from an simple IIS Setup
type Log struct {
	Date          time.Time // The date on which the activity occurred.
	Time          string    // The time, in coordinated universal time (UTC), at which the activity occurred.
	SIP           string    // The IP address of the server on which the log file entry was generated.
	CsMethod      string    // The requested verb, for example, a GET method.
	CsURIStem     string    // The target of the verb, for example, Default.htm.
	CsURIQuery    string    // The query, if any, that the client was trying to perform. A Universal Resource Identifier (URI) query is necessary only for dynamic pages.
	SPort         string    // The server port number that is configured for the service.
	CsUsername    string    // The name of the authenticated user that accessed the server. Anonymous users are indicated by a hyphen.
	CIP           string    // The IP address of the client that made the request.
	CsUserAgent   string    // The browser type that the client used.
	CsReferer     string    // The site that the user last visited. This site provided a link to the current site.
	ScStatus      int       // The HTTP status code.
	ScSubstatus   int       // The substatus error code.
	ScWin32Status int       // The Windows status code.
	TimeTaken     int       // The length of time that the action took, in milliseconds.
}

func loadLog(data []string) Log {
	// I should also return error here too right? And test on failure

	// convert strings to other types
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
	}
	return row
}

func main() {
	// TODO: directory needs to be read from a .config file
	files, err := ioutil.ReadDir("../sample")
	if err != nil {
		log.Fatal(err)
	}

	// Load the history
	processed, errRL := loadHistory("History.dat")
	if errRL != nil {
		log.Fatal(errRL)
	}

	for _, file := range files {
		var logs []Log
		fileName := file.Name()

		filepath := "../sample/" + fileName
		if file, err := os.Open(filepath); err == nil {

			// make sure file gets closed
			defer file.Close()

			// Only parse files which have not been parsed.
			if !processed[fileName] {
				// create a new scanner and read the file line by line
				// I think scanner does have some type of limit on the lennth...
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text() // read in the first line

					// All lines which are not comments should
					// be a line in the log file.  Don't send the comments to home base.
					if (strings.Contains(line, "#")) == false {
						data := strings.Fields(line)
						row := loadLog(data)

						// Add the log record (row) to the logs slice
						logs = append(logs, row)
					}
				}

				// check for errors
				if err = scanner.Err(); err != nil {
					log.Fatal(err)
				}

				// Send data as tightly packed binary form via web service to cloud
				// printJSON(logs[:1])
				header := Header{CustomerName: "Tupoc", Filename: fileName, Logs: logs, Records: len(logs)}
				result := publishData("http://localhost:8080", header)

				if result == "200 OK" {
					// Log the file loaded successfully
					// Make sure we never reach this line of code if an error occurs, or
					// before the endpoint successfully loads the data to the database.
					fmt.Printf("Successful Load of %v : %v\n", fileName, result)
					// Note the append could fail also, and if it does we would re-load the file
					// again so we need a dedupe process on in moya.
					// alert appendToFile failed, and alert should mention file name so we
					// can clean out dupes.
					appendStringToFile("History.dat", fileName)
				}
			}
		} else {
			log.Fatal(err)
		}
	}
}

// publishData v0.1 to pump data to raw endpoint w/ no security....
func publishData(endpoint string, h Header) string {
	// gob encoding
	var data = dataOutAsByte(h)

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	return resp.Status
}

// loadHistory loads the files we have logged from previous
// successul processing and returns a map we can search
// to test if we have already processed the file.
func loadHistory(path string) (map[string]bool, error) {
	files := make(map[string]bool) // create the map

	// Open or create the new file in read only mode.
	// 0666 - chmod 666 means that all users can read and write but cannot execute
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// create a scanner so we can create a map of each filename in this fil
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Setting the value to true,
		// since we have already processed this file
		files[line] = true
	}
	return files, scanner.Err()
}

// appendStringToFile
func appendStringToFile(path, text string) error {
	// Open or create the file, if this is the first record
	// 0600 sets the file permisions to -rw-r--r--
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Write the string, and add in carrige return so we can read these back in as a map
	if _, err = f.WriteString(text + "\n"); err != nil {
		panic(err)
	}
	return nil
}

func dataOutAsByte(h Header) []byte {
	json, err := json.Marshal(h)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// Print the Results
	return json
}
