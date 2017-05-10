package main

import (
	"os"
	"os/signal"
)

type hook struct {
	handlers []func()
}

func (hook *hook) Add(f func()) {
	hook.handlers = append(hook.handlers, f)
}

func createHook() *hook {
	h := &hook{
		handlers: make([]func(), 0),
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		for _, i := range h.handlers {
			i()
		}
		os.Exit(0)
	}()
	return h
}

var instance = createHook()

func GetHook() *hook {
	return instance
}
