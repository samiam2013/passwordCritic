package types

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// DownloadLists encapsulates code for downloading common password list(s) and provides a returned list of the sizes and paths to downloaded files
func DownloadLists() (map[int]string, error) {

	// need to download from my repository's file in binary formate (!! at a permalink !!)
	// 	trying github first with a HEAD request, then downloading
	//	failing over to gitlab ... or panic
	// need to slice the file n times for (k, k*10 .. k *10^n) sizes
	//	and save files with names like 1000.txt, 10000.txt

	sources := map[string]string{
		"github": "https://raw.githubusercontent.com/samiam2013/10-million-password/1b6edff51b6f0f6587a12a659d041ffa18a4d65b/10-million-password-list-top-1000000.txt",
		"gitlab": "https://gitlab.com/samiam2013/10-million-password/-/raw/main/10-million-password-list-top-1000000.txt",
	}

	found := ""
	for name, link := range sources {
		path := "../cache/" + name + ".txt"
		bytesWritten, err := dlFileTo(path, link)
		if err != nil {
			log.Printf("failed getting %s link %s: %s",
				name, link, err.Error())
		} else {
			// breaking here stops unneccesary retries/downloads
			found = path
			break
		}
		fmt.Println("bytes written:", bytesWritten)
	}

	created := map[int]string{}
	if found != "" {
		lengths := []int{1000, 10_000, 100_000}
		for _, length := range lengths {
			// split the file into this many lines and save a copy
			newPath := "../cache/" + strconv.Itoa(length) + ".txt"
			err := copyAndTrunc(found, newPath, length)
			if err != nil {
				return nil, fmt.Errorf("could not copy new file over %s", err.Error())
			}
			created[length] = newPath
		}
	} else {
		return nil, fmt.Errorf("failed at all retries to download list")
	}

	return created, nil
}

// copyAndTrunc copies the 10MM list to a new destination with n lines
func copyAndTrunc(orig, dest string, n int) (err error) {
	var origFH, destFH *os.File
	if origFH, err = os.Open(orig); err != nil {
		return err
	}
	defer origFH.Close()
	if destFH, err = os.Create(dest); err != nil {
		return err
	}
	defer destFH.Close()

	// find the nth line and truncate the destination file
	s := bufio.NewScanner(origFH)
	lines := 0
	for s.Scan() && lines < n {
		if _, err = destFH.Write(append(s.Bytes(), []byte("\n")...)); err != nil {
			return err
		}
		lines++
	}
	return nil
}

func dlFileTo(filepath, url string) (int64, error) {

	out, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	written, err := io.Copy(out, resp.Body)

	return written, err
}
