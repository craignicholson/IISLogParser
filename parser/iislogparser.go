package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir("../sample")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		fmt.Println(file.Name())
		filepath := "../sample/" + file.Name()
		if file, err := os.Open(filepath); err == nil {
			// make sure it gets closed
			defer file.Close()

			// create a new scanner and read the file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				//log.Println(scanner.Text())
				if (strings.Contains(line, "#")) == false {
					//2015-11-04 21:19:38
					//value := line[:19]
					//fmt.Println(value)
					data := strings.Fields(line)
					fmt.Println(data)

					//Get the http response and bytes [500 0 0 109]
					//get the position of the last '-'
					fmt.Println(len(data))
					fmt.Printf("Fields are: %q", data)
					//fmt.Println(data[15])

				}
			}

			// check for errors
			if err = scanner.Err(); err != nil {
				log.Fatal(err)
			}

		} else {
			log.Fatal(err)
		}
	}
}
