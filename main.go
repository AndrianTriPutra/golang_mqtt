package main

import (
	"fmt"
	"log"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	opts   = MQTT.NewClientOptions() //for mqtt client options
	client = MQTT.NewClient(opts)    //for mqtt client

	dns          string
	clientid     string
	connection   bool
	TopicsubConf string
	user         = "kecapi/manggis"
	topic        = "kecapi/test"
)

func main() {
	ConnectMqtt()

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			Timestamp := get_time()
			data := "1234|" + Timestamp
			publish_packet(topic, data)
		}
	}
}
func ConnectMqtt() {
	broker := "broker.hivemq.com:1883" //"broker.emqx.io:1883" //"test.mosquitto.org:1883"
	//user := " "
	//pass := " "

	dns = "tcp://" + broker
	clientid = "gondril/" + user

	opts = MQTT.NewClientOptions()
	opts.AddBroker(dns).SetClientID(clientid)
	//opts.SetUsername(user).SetPassword(pass)
	opts.SetKeepAlive(15 * time.Second)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(15 * time.Second)
	opts.SetAutoReconnect(true).SetMaxReconnectInterval(15 * time.Second)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	opts.SetWill(user+"/wills", "good-bye!", 0, true)

	client = MQTT.NewClient(opts)
	tokenConn := client.Connect()
	if tokenConn.WaitTimeout(15*time.Second) && tokenConn.Error() != nil { //WaitTimeout(10*time.Second)
		log.Println("MQTT Not Connected")
	}
	time.Sleep(1 * time.Second)
	if !connection {
		log.Println("MQTT NOT Connected")
	}

}

var connectHandler MQTT.OnConnectHandler = func(client MQTT.Client) {
	log.Println("MQTT Connected")
	connection = true
	mysubscribe()
}

var connectLostHandler MQTT.ConnectionLostHandler = func(client MQTT.Client, err error) {
	fmt.Printf("Connection MQTT lost: %v", err)
}

func mysubscribe() {
	TopicsubConf = topic
	tokenSubsConf := client.Subscribe(TopicsubConf, 0, messagePubHandler)
	if tokenSubsConf.WaitTimeout(5*time.Second) && tokenSubsConf.Error() != nil {
		//printWarning("SubsConf on subscriber Err")
	} else {
		//printNotif("SubsConf on subscriber succes")
	}
}

func publish_packet(topic, message string) {
	Sending := client.Publish(topic, 1, true, message)
	if Sending.WaitTimeout(3*time.Second) != true {
		log.Printf("send %s failed", topic)
		fmt.Println()
	}
}

var messagePubHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func get_time() string {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	timeutcplus := time.Now().In(loc)

	hhmmssplus := fmt.Sprintf("%02d:%02d:%02d", timeutcplus.Hour(), timeutcplus.Minute(), timeutcplus.Second())
	yymmddplus := fmt.Sprintf("%02d/%02d/%02d", timeutcplus.Day(), timeutcplus.Month(), timeutcplus.Year())
	timestamp := hhmmssplus + " " + yymmddplus
	//log.Printf("timestamp:%s", timestamp)
	return timestamp
}
