package rest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/evoila/BPM-Client/model"
)

const blob = "/blob/"

func DownloadBlob(url, uuid string, boshBlob model.BoshBlob) {

	var path = BuildPath([]string{url, blob, uuid, "?filename=", boshBlob.Name})

	fmt.Println("Downloading", boshBlob.Name)

	response, err := http.Get(path)

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	fmt.Println("Download ResponseCode:", response.StatusCode)
	fmt.Println("Download response Headers", response.Header)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(boshBlob.Name)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		panic(err)
	}

	time.Sleep(10000000000)
	fmt.Println(n, "bytes downloaded.")
}
