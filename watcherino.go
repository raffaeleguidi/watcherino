// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"os/exec"
	"time"
	"path/filepath"
	"gopkg.in/alecthomas/kingpin.v2"
)

func Execute(command string, folder string, file string, event string) {
	log.Println("Executing", command, folder, file, event)
	cmd := exec.Command(command, folder, file, event)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Println("Command finished with error: %v", err)
	} else {
		log.Printf("Command executed successfully")
	}
}

var timer *time.Timer
var second bool

func Delay (command string, folder string, file string, event string, delay int) {
	if (second) {
		timer.Reset(time.Second * time.Duration(delay)) 
		return
	} else {
		second = true
	}
	t := time.NewTimer(time.Second * time.Duration(delay))
	go func() {
		timer = t;
        <-t.C
        second = false;
        Execute(command, folder, file, event)
    }()
}

func Watcher(folder string, pattern string, command string, delay int) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				write := event.Op&fsnotify.Write == fsnotify.Write
				create := event.Op&fsnotify.Create == fsnotify.Create
				if (write || create) {
				    filename := event.Name[len(folder)+1:1+len(event.Name)-len(folder)]
					matched, err := filepath.Match(pattern, filename)
					if err != nil {
			        	log.Println(err)
					} else {
						eventType := "WRITE"
						if (create) { eventType = "CREATE" }
						if (matched) {
							Delay(command, folder, filename, eventType, delay)
						}
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

var (
	command = kingpin.Arg("command", "Command to execute").Required().String()
	folder  = kingpin.Arg("folder", "Folder to watch for changes").Default(".").String()
	pattern  = kingpin.Flag("pattern", "Pattern filter").Default("*").String()
	delay   = kingpin.Flag("delay", "Number seconds to wait from last change").Default("5").Int()
)

func main() {
	kingpin.Version("0.0.1")
	kingpin.Parse()
	log.Println("Watcherino executing", *command, "for changes in folder", *folder, "pattern", *pattern, "with a", *delay, "seconds delay" )
	Watcher(*folder, *pattern, *command, *delay)
}
