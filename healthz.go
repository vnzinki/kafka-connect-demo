package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Connector struct {
	Status struct {
		Name      string `json:"name"`
		Connector struct {
			State    string `json:"state"`
			WorkerID string `json:"worker_id"`
		} `json:"connector"`
		Tasks []struct {
			ID       int    `json:"id"`
			State    string `json:"state"`
			WorkerID string `json:"worker_id"`
		} `json:"tasks"`
		Type string `json:"type"`
	} `json:"status"`
}

type ConnectorStatus map[string]Connector

func main() {
	connectURI := flag.String("url", "", "Connect URI")
	flag.Parse()

	resp, err := http.Get(fmt.Sprintf("%s/connectors?expand=status", *connectURI))
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	ConnectorStatus := ConnectorStatus{}
	json.Unmarshal([]byte(resBody), &ConnectorStatus)

	for _, connector := range ConnectorStatus {
		if connector.Status.Connector.State == "RUNNING" {
			fmt.Printf("%s is RUNNING\n", connector.Status.Name)
		} else {
			fmt.Printf("%s is NOT RUNNING - Restarting\n", connector.Status.Name)

			restartURI := fmt.Sprintf("%s/connectors/%s/restart", *connectURI, connector.Status.Name)

			_, err = http.Post(restartURI, "application/json", nil)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			continue // Only restart one connector, stop checking tasks
		}

		for _, tasks := range connector.Status.Tasks {
			if tasks.State == "RUNNING" {
				fmt.Printf("%s task %d is RUNNING\n", connector.Status.Name, tasks.ID)
			} else {
				fmt.Printf("%s task %d is NOT RUNNING - Restarting\n", connector.Status.Name, tasks.ID)

				restartURI := fmt.Sprintf("%s/connectors/%s/tasks/%d/restart", *connectURI, connector.Status.Name, tasks.ID)

				_, err := http.Post(restartURI, "application/json", bytes.NewBuffer([]byte{}))
				if err != nil {
					fmt.Printf("%s\n", err)
					os.Exit(1)
				}
			}
		}
	}

}
