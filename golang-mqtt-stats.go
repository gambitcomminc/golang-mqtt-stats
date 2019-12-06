// Copyright (c) 2019 Gambit Communications, Inc.

package main

import (
	"fmt"
	"os"
	"time"
	"github.com/pborman/getopt/v2"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	optHost := getopt.StringLong("host", 'h', "", "broker host")
	optPort := getopt.StringLong("port", 'p', "1883", "port (default 1883)")
	optTopic := getopt.StringLong("topic", 't', "", "topic")
	optQos := getopt.IntLong("qos", 'q', 0, "QOS (default 0)")
	optUser := getopt.StringLong("user", 'u', "", "username")
	optPassword := getopt.StringLong("password", 'P', "", "password")
	optStore := getopt.StringLong("store", 's', ":memory:", "store directory (default :memory)")
	optHelp := getopt.BoolLong("help", 0, "Help")
	getopt.Parse()

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	if *optHost == "" {
		fmt.Println("ERROR: Need --host ")
		getopt.Usage()
		os.Exit(-1)
	}

	if *optTopic == "" {
		fmt.Println("ERROR: Need --topic ")
		getopt.Usage()
		os.Exit(-1)
	}

	broker := *optHost + ":" + *optPort
	fmt.Println("Subscribing to topic " + *optTopic + " on " + broker)

	mqttOpts := MQTT.NewClientOptions()
	mqttOpts.AddBroker(broker)
	mqttOpts.SetUsername(*optUser)
	mqttOpts.SetPassword(*optPassword)
	if *optStore != ":memory:" {
		mqttOpts.SetStore(MQTT.NewFileStore(*optStore))
	}

	inChan := make(chan [2]string)

	mqttOpts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		inChan <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := MQTT.NewClient(mqttOpts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(*optTopic, byte(*optQos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// single threaded counter
	receiveCount := 0
	lastReportedCount := 0
	lastReportedTime := time.Now()
	countBytes := 0
	reportSeconds := 5
	reportInterval := time.Second * time.Duration(reportSeconds)
	for {
		incoming := <-inChan
		receiveCount++
		countBytes = countBytes + len(incoming[1])

		// every 5 seconds, print perSecond
		now := time.Now()
		elapsed := now.Sub(lastReportedTime)
		if elapsed >= reportInterval { 
			msgPerSecond := float64 (receiveCount - lastReportedCount) / elapsed.Seconds()
			bytesPerSec := float64 (countBytes) / elapsed.Seconds()
			bytesPerMsg := bytesPerSec / msgPerSecond

			fmt.Printf ("elapsed %.5s: msgs/s %.1f, bytes/s %.1f, bytes/msg %.1f\n", elapsed.String(), msgPerSecond, bytesPerSec, bytesPerMsg)
			lastReportedTime = now
			lastReportedCount = receiveCount
			countBytes = 0
//			fmt.Printf("topic: %s msg: %s\n", incoming[0], incoming[1])
		}
	}

	client.Disconnect(500)
}
