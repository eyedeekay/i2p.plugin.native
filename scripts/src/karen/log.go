package main

import "log"

func Printf(format string, v ...interface{}) {
	if *verbose {
		if v == nil {
			log.Printf(format)
		} else {
			log.Printf(format, v...)
		}
	}
}

func Println(v ...interface{}) {
	if *verbose {
		log.Println(v...)
	}
}
