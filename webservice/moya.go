// Moya is a test application will will recieve json
// and post the data to mongodb.
// TODO:  Should be clean the data here or change the schema for our own
//        database?
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

// Log holds the default fields from an simple IIS Setup
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
	Customer      string    // The unique customers name for grouping and viewing.
	Filename      string    // The file which contained the logs.
}

// This only writes out to the web browser...
// We can do a handler for each diff. data set
func handler(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	if len(data) < -1 {
		fmt.Printf("%s", data)
	}
	var logs []Log
	err = json.Unmarshal(data, &logs)
	if err != nil {
		panic(err)
	}

	printJSON(logs)

	// err = InsertToDatabase(h)
	// if err != nil {
	// 	panic(err)
	// }
}

// TODO Add Basic Auth to check for user and password
// TODO Add Rate Limiting and Throttling
func main() {
	http.HandleFunc("/", handler)

	// TODO set port in a config
	http.ListenAndServe(":8080", nil)
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

// TODO Set Mongodb in config, return ACK on success, true / false
// the database could be down
func InsertToDatabase(data Log) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	// Database name and collection, these need to be
	// hardcoded because we need them to always be the
	// same to have consitency across the applications
	c := session.DB("customername").C("loghist")

	err = c.Insert(&data)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
