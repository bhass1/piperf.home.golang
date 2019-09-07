package main

import (
	"fmt"
	"os"
	"os/exec"
	"log"
	"time"
	"github.com/eclipse/paho.mqtt.golang"
)

var debug_print = func(message string) {
	t := time.Now()
	fmt.Println(t.Format(time.RFC850) + " " + message)
}

var bar mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	debug_print("Connection lost")
}

var baz mqtt.ReconnectHandler = func(client mqtt.Client, copts *mqtt.ClientOptions) {
	debug_print("Reconnecting...")
}

var connect_mqtt = func(c mqtt.Client) {
	debug_print("Connecting...")
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	debug_print("Ok!")
}

var shutdown_mqtt = func(c mqtt.Client) {
	debug_print("Shutting down...")
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
	opts.SetPingTimeout(1 * time.Second)
	opts.SetConnectionLostHandler(bar)
	opts.SetAutoReconnect(true)
	opts.SetReconnectingHandler(baz)

	c := mqtt.NewClient(opts)
	connect_mqtt(c)
	for {
		if retry {
			retry = false
			connect_mqtt(c)
		}

		//debug_print("Iperf...")
		out, err := exec.Command("./do_iperf.sh").Output()
		if err != nil {
			fmt.Println("iperf failure:")
			fmt.Println(err)
		} else {
			c.Publish(topic, 0, false, out)
			//debug_print("published")
		}
		time.Sleep(1*time.Second)
	}

	shutdown_mqtt(c)
}
