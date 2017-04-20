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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const usage = `Tiny server for serving up files locally

Usage:
  servem [--conf <dir>] [--listen <address>] [--serve <dir>]
  servem -h | --help
  servem --howto
  servem --version

Options:
  -h --help           Show this screen.
  --howto             Print description howto use this server.
  --version           Print version.
  --conf <dir>        Path to directory where the required 'goodies.yml' can be found. default: .
  --serve <dir>       Path to directory where files should be served from.
  --listen <address>  Port to listen at [default: :8083]
  `

var serveDirectory = "/home/msonntag/Chaos/work/gnode-bootstrap-theme/build"
var port = ":3030"

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Servem running")
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	filename := "bootstrap.css"
	path := serveDirectory + "/css/" + filename

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] serving css: %v", err)
		fmt.Fprintln(w, "Cannot serve css")
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(content))
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	filename := "bootstrap.js"
	path := serveDirectory + "/js/" + filename

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] serving js: %v", err)
		fmt.Fprintln(w, "Cannot serve js")
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(content))
}

func serveImageFile(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["file"]
	path := serveDirectory + "/img/" + filename

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] serving image: %v\n", err)
		fmt.Fprintln(w, "Cannot serve image")
		return
	}
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(content))
}

func serveBuildFile(w http.ResponseWriter, r *http.Request) {
	filepath := mux.Vars(r)["remainder"]
	path := serveDirectory + "/" + filepath

	fmt.Fprintf(os.Stdout, "[Server] serving file: %s\n", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] serving build file: %v\n", err)
		fmt.Fprintln(w, "Cannot serve requested file")
		return
	}
	http.ServeContent(w, r, filepath, time.Now(), bytes.NewReader(content))
}

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/servecss", serveCSS)
	r.HandleFunc("/servejs", serveJS)
	r.HandleFunc("/img/{file}", serveImageFile)
	r.HandleFunc(`/build/{remainder:[a-zA-Z0-9=\-\/]*}`, serveBuildFile)
}

func main() {
	args, err := docopt.Parse(usage, nil, true, "v1.0.0", false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] parsing cli arguments: '%s', abort...\n\n", err.Error())
		os.Exit(-1)
	}
	fmt.Fprintf(os.Stdout, "[Warmup] cli arguments: %v\n", args)

	// Use individual port if provided. Crude check whether leading colon is available.
	if p, ok := args["--listen"]; ok {
		port = p.(string)
		if string(port[0]) != ":" {
			port = ":" + port
		}
		fmt.Fprintf(os.Stdout, "[Warmup] using port: '%s'\n", port)
	}

	if srv, ok := args["--serve"]; ok && srv != nil {
		serveDirectory = srv.(string)
		fmt.Fprintf(os.Stdout, "[Warmup] serving files from directory: '%s'\n", serveDirectory)
	}

	/*
		// Check whether goodies file was provided.
		goodies := "goodies.yml"
		path := "./" + goodies
		fmt.Fprintf(os.Stdout, "[Warmup] checking goodies file at '%s'\n", path)

		stats, err := os.Stat(path)
		if err != nil {
			fmt.Fprint(os.Stderr, "\n[Error] goodies file not found, abort...\n\n")
			os.Exit(-1)
		}
		if stats.Size() == 0 {
			fmt.Fprint(os.Stderr, "\n[Error] goodies file empty, abort...\n\n")
			os.Exit(-1)
		}

		fmt.Fprintf(os.Stdout, "[Warmup] file size: %v\n", stats.Size())
	*/
	fmt.Fprintln(os.Stdout, "[Warmup] registering routes")
	router := mux.NewRouter()
	registerRoutes(router)

	handler := handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Authorization"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET"}),
	)(router)

	server := http.Server{
		Addr:    port,
		Handler: handler,
	}

	fmt.Fprintln(os.Stdout, "[Start] Listen and serve")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[Error] Server startup: '%v', abort...\n\n", err)
		os.Exit(-1)
	}
}
