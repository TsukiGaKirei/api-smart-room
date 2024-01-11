package auth

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// for publish
const (
	mqttBrokerHost     = "34.101.245.2"
	mqttBrokerPort     = 1883
	mqttBrokerUsername = "user1"
	mqttBrokerPassword = "qweasd123"
	mqttTopic          = "api-topic"
)

// PublishMessage publishes a message to the MQTT broker
func PublishMessage(data string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqttBrokerHost, mqttBrokerPort))
	opts.SetUsername(mqttBrokerUsername)
	opts.SetPassword(mqttBrokerPassword)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to MQTT broker: %v", token.Error())
	}
	defer client.Disconnect(0)

	token := client.Publish(mqttTopic, 0, false, data)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to publish message: %v", token.Error())
	}

	fmt.Println("Message published successfully.")
}
