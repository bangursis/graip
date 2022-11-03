package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/gookit/validate"
)

type Job struct {
	Name     string `json:"name" validate:"required"`
	Resource string `json:"resource" validate:"required|url"`
	Period   int    `json:"schedule_every_X_minutes" validate:"required"`
	SaveTo   string `json:"save_to" validate:"required"`
}

var jobs []Job

func init() {
	config, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	defer config.Close()

	data, err := io.ReadAll(config)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &jobs); err != nil {
		panic(err)
	}

	for i := range jobs {
		v := validate.Struct(&jobs[i])

		if !v.Validate() {
			panic(v.Errors)
		}
	}
}

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	qCh, errCh := make(chan struct{}), make(chan error)
	for _, j := range jobs {
		os.Create(j.SaveTo)
		go runner(qCh, errCh, j)
	}

	for {
		select {
		case err := <-errCh:
			fmt.Println(err)
		case <-sigCh:
			for _ = range jobs {
				qCh <- struct{}{}
			}
			os.Exit(1)
		}

	}
}
