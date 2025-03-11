package api

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SetHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func HandleDmList(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	cid := r.URL.Query().Get("cid")
	if cid == "" {
		http.Error(w, "oid parameter is required", http.StatusBadRequest)
		return
	}

	url := "https://api.bilibili.com/x/v1/dm/list.so?oid=" + cid
	method := "GET"
	payload := strings.NewReader("oid=" + cid)

	fmt.Println(url, method, payload)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Origin", "https://bilibili.com")
	req.Header.Add("Referer", "https://bilibili.com/")

	res, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contentEncoding := res.Header.Get("Content-Encoding")
	if contentEncoding == "gzip" {
		reader, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if contentEncoding == "deflate" {
		reader := flate.NewReader(bytes.NewReader(body))
		defer reader.Close()
		body, err = io.ReadAll(reader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(body)
}
