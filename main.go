/*
	Nothing too Complicated.

	Practice code Using
		- Requests
		- Regex
		- File Managing (??)

	Asks for an URL unless it recieves an argument, makes a request to the URL
	and uses regex to get the URL, lastly it makes a request to the URL with the
	file and copies the data into the recently made File.

*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

// This is what makes the name of the File
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func HandleError(err error, s *spinner.Spinner) {
	if err != nil {
		(s).Prefix = "Error... "
		(s).Stop()
		panic(err)
	}
}

func DownloadFile(path string, URL string, s *spinner.Spinner) {
	s.Prefix = "Starting to Download "
	resp, err := http.Get(URL)
	HandleError(err, s)
	defer resp.Body.Close()
	s.Prefix = "Downloading Data... "
	file, err := os.Create(path)
	HandleError(err, s)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	HandleError(err, s)
	s.Prefix = "Finished! "
	s.Stop()
	fmt.Println("Done!")
}

func main() {
	var URL string
	if len(os.Args) > 1 {
		URL = os.Args[1]
	} else {
		fmt.Printf("MediaFire URL: ")
		fmt.Scanln(&URL)
	}
	if len(URL) > 9 {
		s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
		s.Prefix = "Making Request "
		s.Start()
		resp, err := http.Get(URL)
		HandleError(err, s)
		s.Prefix = "Parsing Body "
		body, err := ioutil.ReadAll(resp.Body)
		HandleError(err, s)
		var FullBody string = string(body)
		s.Prefix = "Finding Download URL "
		regexp, _ := regexp.Compile(`(?m)https:\/\/download([0-9]*)\.mediafire.com\/(.*)\/(.*)\/(.*)\.[^"]+`)
		var DownloadURLMatch = regexp.FindString(FullBody)
		if DownloadURLMatch != "" {
			s.Prefix = "Generating File Name "
			splitted := strings.Split(DownloadURLMatch, ".")
			DownloadFile(fmt.Sprintf("%v.%v", RandomString(10), splitted[len(splitted)-1]), DownloadURLMatch, s)
		} else {
			s.Stop()
			fmt.Println("No matches found in URL")
		}
	} else {
		main()
	}
}
