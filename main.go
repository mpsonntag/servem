// Copyright (c) 2017, Michael Sonntag (michael.p.sonntag@gmail.com)
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted under the terms of the BSD License.
// (https://github.com/mpsonntag/servem/blob/master/LICENSE)

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Servem running")
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	filename := "bootstrap.css"
	path := "/home/msonntag/Chaos/work/gnode-bootstrap-theme/build/css/" + filename

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] serving css: %v", err)
		fmt.Fprintln(w, "Cannot serve css")
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(content))
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/servecss", serveCSS)
}

func main() {
	const port = ":3030"

	fmt.Println("[Starting server] Registering routes")
	router := mux.NewRouter()
	registerRoutes(router)

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Println("[Starting server] Listen and serve")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] Server startup: %v\n", err)
		os.Exit(-1)
	}
}
