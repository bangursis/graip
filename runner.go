package main

import (
	"io"
	"net/http"
	"os"
	"time"
)

func runner(qCh <-chan struct{}, errCh chan<- error, j *Job) {
	for {
		select {
		case <-time.Tick(time.Duration(j.Period) * time.Minute):
			resp, err := http.Get(j.Resource)
			if err != nil {
				errCh <- formatErr(HTTPGetError, j.Name, err)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errCh <- formatErr(ReadResponseError, j.Name, err)
			}
			resp.Body.Close()

			if err := os.WriteFile(j.SaveTo, body, os.ModeAppend); err != nil {
				errCh <- formatErr(WriteToFileError, j.Name, err)
			}

		case <-qCh:
			return
		}
	}
}
