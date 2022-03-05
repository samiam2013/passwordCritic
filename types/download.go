package types

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadLists encapsulates code for downloading common password list(s) and provides a returned list of the sizes and paths to downloaded files
func DownloadLists() (map[int]string, error) {

	// need to download from my repository's file in binary formate (!! at a permalink !!)
	// 	trying github first with a HEAD request, then downloading
	//	failing over to gitlab ... or panic
	// need to slice the file n times for (k, k*10 .. k *10^n) sizes
	//	and save files with names like 1000.txt, 10000.txt

	return map[int]string{}, fmt.Errorf("not implemented")
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
