// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"os/exec"
	"time"
	"strconv"
)

func Execute(command string, folder string, file string, event string) {
	log.Println("Executing", command, folder, file, event)
	cmd := exec.Command(command, folder, file, event)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("command started")
	err = cmd.Wait()
	if err != nil {
		log.Println("command finished with error: %v", err)
	} else {
		log.Printf("command ended")
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

func Watcher(folder string, command string, delay int) {
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
				log.Println("event:", event)
				write := event.Op&fsnotify.Write == fsnotify.Write
				create := event.Op&fsnotify.Create == fsnotify.Create
				if (write || create){
					eventType := "WRITE"
					if (create) { eventType = "CREATE" }
					Delay(command, folder, event.Name, eventType, delay)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(folder)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main(){
	folder := "."
	if (len(os.Args)) > 1 {
		folder = os.Args[1]
	}
	command := "echo"
	if (len(os.Args)) > 2 {
		command = os.Args[2]
	}
	delay := 5
	if (len(os.Args)) > 3 {
		input, err := strconv.Atoi(os.Args[3])
		if (err != nil) {
			delay = input
		}
	}
	log.Println("Watcherino watching folder:", folder, "and in case executing:", command, "with delay", delay, "seconds")
	Watcher(folder, command, delay)
}