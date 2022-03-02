// Licensed to You under the Apache License, Version 2.0.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"gitlab.pgre.dell.com/enterprise/telemetryservice/internal/databus"
	"gitlab.pgre.dell.com/enterprise/telemetryservice/internal/messagebus/stomp"
)

var configStrings = map[string]string{
	"mbhost": "activemq",
	"mbport": "61613",
	"URL":    "http://localhost:8086",
	"DBName": "poweredge_telemetry_metrics",
}

func handleGroups(writeAPI api.WriteAPI, groupsChan chan *databus.DataGroup) {
	for group := range groupsChan {
		for _, value := range group.Values {
			floatVal, _ := strconv.ParseFloat(value.Value, 64)

			timestamp, err := time.Parse(time.RFC3339, value.Timestamp)
			fmt.Printf("Value: %#v     Timestamp: '%v'  -->  '%v'\n", value, value.Timestamp, timestamp)
			if err != nil {
				log.Printf("Error parsing timestamp for  point %s: (%s) %v", value.Context+"_"+value.ID, value.Timestamp, err)
				continue
			}

			p := write.NewPointWithMeasurement("telemetry").
				AddTag("ServiceTag", value.System).
				AddTag("FQDD", value.Context).
				AddTag("Label", value.Label).
				AddField(value.ID, floatVal).
				SetTime(timestamp)

			// automatically batches things behind the scenes
			writeAPI.WritePoint(p)
		}
	}
}

func getEnvSettings() {
	fmt.Printf("Environment dump: %#v\n", os.Environ())
	mbHost := os.Getenv("MESSAGEBUS_HOST")
	if len(mbHost) > 0 {
		configStrings["mbhost"] = mbHost
	}
	mbPort := os.Getenv("MESSAGEBUS_PORT")
	if len(mbPort) > 0 {
		configStrings["mbport"] = mbPort
	}
	configStrings["URL"] = os.Getenv("INFLUXDB_URL")
	configStrings["DBName"] = os.Getenv("INFLUXDB_DB")
	configStrings["User"] = os.Getenv("INFLUXDB_USER")
	configStrings["Pass"] = os.Getenv("INFLUXDB_PASS")
	configStrings["Org"] = os.Getenv("INFLUX_ORG")
	configStrings["Bucket"] = os.Getenv("INFLUX_BUCKET")
	configStrings["Token"] = os.Getenv("INFLUX_TOKEN")
}

func main() {
	ctx := context.Background()

	//Gather configuration from environment variables
	getEnvSettings()

	dbClient := new(databus.DataBusClient)
	stompPort, _ := strconv.Atoi(configStrings["mbport"])
	for {
		mb, err := stomp.NewStompMessageBus(configStrings["mbhost"], stompPort)
		if err != nil {
			log.Printf("Could not connect to message bus: %s", err)
			time.Sleep(5 * time.Second)
		} else {
			dbClient.Bus = mb
			defer mb.Close()
			break
		}
	}

	groupsIn := make(chan *databus.DataGroup, 10)
	dbClient.Subscribe("/influx")
	dbClient.Get("/influx")
	go dbClient.GetGroup(groupsIn, "/influx")

	if configStrings["Pass"] == "" {
		log.Fatalf("must specify influx user/pass")
	}

	for {
		time.Sleep(5 * time.Second)

		client := influxdb2.NewClientWithOptions(
			configStrings["URL"],
			configStrings["Token"],
			influxdb2.DefaultOptions().SetBatchSize(5000),
		)
		writeAPI := client.WriteAPI(configStrings["Org"], configStrings["Bucket"]) // async, non-blocking

		go func(writeAPI api.WriteAPI) {
			for err := range writeAPI.Errors() {
				fmt.Printf("async write error: %s\n", err)
			}
		}(writeAPI)

		// TODO: Sensitive, for debugging only, remove once it works
		log.Printf("Connected to influx db(%s) org(%s) bucket(%s) at URL (%s) using (%s)\n", configStrings["DBName"], configStrings["Org"], configStrings["Bucket"], configStrings["URL"], configStrings["Token"])

		timedCtx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		ok, err := client.Ping(timedCtx)
		cancel()
		if !ok || err != nil {
			log.Printf("influx ping return = (%t): %s\n", ok, err)
			client.Close()
			continue
		}

		log.Printf("influx ping return = (%t)\n", ok)
		defer client.Close()
		handleGroups(writeAPI, groupsIn)
	}
}
