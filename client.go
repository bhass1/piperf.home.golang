package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"time"
	"github.com/eclipse/paho.mqtt.golang"
)

var bar mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println("Connection lost")
	retry = true
}

var connect_mqtt = func(c mqtt.Client) {
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
}

var shutdown = func(c mqtt.Client) {
	fmt.Printf("Shutting down...\n")
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}


var topic string = "piperf/result"
var broker string = "aws_mqtt_broker:1883"
var cid string = "piperf_home"
var retry = false

func main() {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(cid)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(foo)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectionLostHandler(bar)

	c := mqtt.NewClient(opts)
	connect_mqtt(c)
	for {
		if retry {
			retry = false
			connect_mqtt(c)
		}

		out, err := exec.Command("./do_iperf.sh").Output()
		if err != nil {
			fmt.Println("iperf failure:")
			fmt.Println(err)
		} else {
			c.Publish(topic, 0, false, out)
		}
		time.Sleep(1*time.Second)
	}

	shutdown(c)
}
