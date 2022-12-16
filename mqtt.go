package teslanotify

import (
	"errors"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TeslaMateCarState = "teslamate/cars/%d/state"
)

var (
	ErrMQTTHostMissing = errors.New("MQTT host missing")
	ErrMQTTPortMissing = errors.New("MQTT port missing")
)

type Service struct {
	Debug         bool
	Log           *log.Logger
	MQTTHost      string
	MQTTPort      string
	CarID         int
	OnStateChange func(state string) error
}

func (s *Service) Connect() error {
	if len(s.MQTTHost) == 0 {
		return ErrMQTTHostMissing
	}

	if len(s.MQTTPort) == 0 {
		return ErrMQTTPortMissing
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", s.MQTTHost, s.MQTTPort))
	opts.SetClientID("teslanotify")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetDefaultPublishHandler(s.messagePubHandler)
	opts.SetOnConnectHandler(s.connectHandler)
	opts.SetConnectionLostHandler(s.connectLostHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	token := client.Subscribe(fmt.Sprintf(TeslaMateCarState, s.CarID), 1, nil)
	token.Wait()

	return nil
}

func (s *Service) messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	switch msg.Topic() {
	case fmt.Sprintf(TeslaMateCarState, s.CarID):
		err := s.OnStateChange(string(msg.Payload()))
		if err != nil {
			s.Log.Printf("State handler error: %v\n", err)
		}
	}
}

func (s *Service) connectHandler(client mqtt.Client) {
	if s.Debug {
		s.Log.Printf("Connected to MQTT %s:%s\n", s.MQTTHost, s.MQTTPort)
	}
}

func (s *Service) connectLostHandler(client mqtt.Client, err error) {
	if s.Debug {
		s.Log.Printf("Connection to MQTT %s:%s lost: %v\n", s.MQTTHost, s.MQTTPort, err)
	}
}
