package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tuotoo/qrcode"
)

func main() {
	intf := flag.String("interface", ":8080", "The networkinterface and port to run on")
	flag.Parse()
	http.HandleFunc("/", decodeBin)
	//http.HandleFunc("/base64", decodeBase64)
	log.Println(http.ListenAndServe(*intf, nil))
}

func testFunc() {
	/*fi, err := os.Open("test.png")
	if err != nil {
		log.Fatal(err)
	}
	defer fi.Close()
	info, err := readImage(fi)
	log.Println(info)*/
}

func decodeBase64(w http.ResponseWriter, r *http.Request) {
	str64 := r.URL.Query().Get("qr")
	b, err := base64.StdEncoding.DecodeString(str64)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fi := bytes.NewReader(b)
	str, err := readImage(fi)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, str)

}

func decodeBin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, files := range r.MultipartForm.File {
		for _, file := range files {
			log.Println(file.Filename)
			f, err := file.Open()
			if err != nil {
				log.Println(err)
				continue
			}
			defer f.Close()
			str, err := readImage(f)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, str+"\n")
		}
	}
}

func readImage(fi io.Reader) (string, error) {
	qrinfo, err := qrcode.Decode(fi)
	if err != nil {
		return "", err
	}
	return qrinfo.Content, nil

}
