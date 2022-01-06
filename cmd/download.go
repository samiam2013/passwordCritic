package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
)

/* DownloadLists encapsulates code for downloading common password list(s)
 * 	and provides a returned list of the sizes and paths to downloaded files */
func DownloadLists() (pwLists map[int]string, err error) {
	const dmCommonFolderURL string = "https://raw.githubusercontent.com/danielmiessler" +
		"/SecLists/master/Passwords/Common-Credentials/"
	pwLists = map[int]string{
		1_000_000: dmCommonFolderURL + "10-million-password-list-top-1000000.txt",
		100_000:   dmCommonFolderURL + "10-million-password-list-top-100000.txt",
		10_000:    dmCommonFolderURL + "10-million-password-list-top-10000.txt",
		1000:      dmCommonFolderURL + "10-million-password-list-top-1000.txt",
	}

	for listSize, url := range pwLists {
		newFilename := strconv.Itoa(listSize) + ".txt"
		newFilePath := CacheFolderPath + string(os.PathSeparator) + newFilename
		_, err = dlFileTo(newFilePath, url) // TODO fail on size too small instead of ignoring
		if err != nil {
			return // TODO is there another option? another source to fall back on?
		} else {
			pwLists[listSize] = newFilePath
		}
	}

	return
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
