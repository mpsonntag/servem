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

	"github.com/docopt/docopt-go"
	"github.com/gorilla/mux"
)

const usage = `Tiny server for serving up files locally

Usage:
  servem [--conf <dir>] [--listen <address>]
  servem -h | --help
  servem --howto
  servem --version

Options:
  -h --help           Show this screen.
  --howto             Print description howto use this server.
  --version           Print version.
  --conf <dir>        Path to directory where the required goodies.yml can be found. default: .
  --listen <address>  Port to listen at [default: :8083]
  `

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
	var port = ":3030"

	args, err := docopt.Parse(usage, nil, true, "v1.0.0", false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] parsing cli arguments: %s\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] cli arguments: %v\n", args)

	// Use individual port if provided. Crude check whether leading colon is available.
	if p, ok := args["--listen"]; ok {
		port = p.(string)
		if string(port[0]) != ":" {
			port = ":"+ port
		}
		fmt.Fprintf(os.Stdout, "[Warmup] Using port: '%s'\n", port)
	}

	fmt.Fprintln(os.Stdout, "[Warmup] Registering routes")
	router := mux.NewRouter()
	registerRoutes(router)

	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	fmt.Fprintln(os.Stdout, "[Start] Listen and serve")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] Server startup: %v\n", err)
		os.Exit(-1)
	}
}
