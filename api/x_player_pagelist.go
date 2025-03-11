package api

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
)

func HandlePlayerPagelist(w http.ResponseWriter, r *http.Request) {
	SetHeaders(w)

	bvid := r.URL.Query().Get("bvid")
	if bvid == "" {
		http.Error(w, "bvid parameter is required", http.StatusBadRequest)
		return
	}

	values := url.Values{}
	values.Add("bvid", bvid)
	url := "https://api.bilibili.com/x/player/pagelist"
	url += "?" + values.Encode()

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
