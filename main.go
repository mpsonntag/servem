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

func genericServe(w http.ResponseWriter, r *http.Request, filePath string, fromRoute string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Error] serving %s: %v", fromRoute, err)
		fmt.Fprintf(w, "Cannot serve %s\n", fromRoute)
		return
	}
	fmt.Fprintf(os.Stdout, "[Serve: %s] %s\n", fromRoute, filePath)
	http.ServeContent(w, r, filePath, time.Now(), bytes.NewReader(content))
}

func serveImageFile(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["file"]
	path := fmt.Sprintf("%s/img/%s", serveDirectory, filename)

	genericServe(w, r, path, "image")
}

func serveFontsFile(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["file"]
	path := fmt.Sprintf("%s/fonts/%s", serveDirectory, filename)

	genericServe(w, r, path, "fonts")
}

func serveAnyFile(w http.ResponseWriter, r *http.Request) {
	dir := mux.Vars(r)["dir"]
	filepath := mux.Vars(r)["file"]
	path := fmt.Sprintf("%s/%s/%s", serveDirectory, dir, filepath)

	genericServe(w, r, path, "dynFile")
}

func serveDeepDir(w http.ResponseWriter, r *http.Request) {
	dir := mux.Vars(r)["dir"]
	deepdir := mux.Vars(r)["deepdir"]
	file := mux.Vars(r)["file"]
	fmt.Fprintf(os.Stdout, "Trying to serve %s\n", dir)
	path := fmt.Sprintf("%s/%s/%s/%s", serveDirectory, dir, deepdir, file)

	genericServe(w, r, path, "dynFile")
}

// TODO server variable depth url
func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", root)
	r.HandleFunc("/img/{file}", serveImageFile)
	r.HandleFunc(`/fonts/{file}`, serveFontsFile)
	r.HandleFunc(`/{dir}/{file}`, serveAnyFile)
	r.HandleFunc(`/{dir}/{deepdir}/{file}`, serveDeepDir)
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
