package main

import (
	"fmt"
	"net/http"

	"github.com/Jinof/dm-server/api"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.SetHeaders(w)
		http.NotFound(w, r)
	})
	http.HandleFunc("/api/x_player_pagelist", api.HandlePlayerPagelist)
	http.HandleFunc("/api/x_v1_dm_list.so", api.HandleDmList)
	fmt.Println("Server started at :8080")
	err := http.ListenAndServe(":8080", nil)
	fmt.Println(err)
}
