package types

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// DownloadLists encapsulates code for downloading common password list(s) and provides a returned list of the sizes and paths to downloaded files
func DownloadLists() (pwLists map[int]string, err error) {
	const dmCommonFolderURL string = "https://raw.githubusercontent.com/danielmiessler/SecLists/" +
		"aa0eb72f3871b01372596a34fb5378910df50073/Passwords/Common-Credentials/"
	pwLists = map[int]string{
		// 1_000_000: dmCommonFolderURL + "10-million-password-list-top-1000000.txt",
		100_000: dmCommonFolderURL + "10-million-password-list-top-100000.txt",
		10_000:  dmCommonFolderURL + "10-million-password-list-top-10000.txt",
		1000:    dmCommonFolderURL + "10-million-password-list-top-1000.txt",
	}

	waitChan := make(chan bool, len(pwLists))

	for listSize, url := range pwLists {
		log.Print("downloading url: ", url)
		newFilename := strconv.Itoa(listSize) + ".txt"
		newFilePath := CacheFolderPath + string(os.PathSeparator) + newFilename

		go dlConcurrent(newFilePath, url, waitChan)

		pwLists[listSize] = newFilePath
	}

	for i := 0; i < len(pwLists); i++ {
		log.Printf("finished? %v", <-waitChan)
	}

	return
}

func dlConcurrent(filepath, url string, finished chan bool) {
	_, err := dlFileTo(filepath, url)
	if err != nil {
		log.Print("could not download file: " + err.Error())
		finished <- false
	}

	defer func() {
		finished <- true
	}()
}

func dlFileTo(filepath, url string) (written int64, err error) {
	var out *os.File
	if out, err = os.Create(filepath); err != nil {
		return
	}
	defer out.Close()

	var resp *http.Response
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	written, err = io.Copy(out, resp.Body)

	return
}
