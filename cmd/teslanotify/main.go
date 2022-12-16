package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"text/template"

	"github.com/hnesland/teslanotify"
)

const (
	ENV_DEBUG        = "DEBUG"
	ENV_NTFY_URL     = "NTFY_URL"
	ENV_NTFY_TOPIC   = "NTFY_TOPIC"
	ENV_NTFY_MSG     = "NTFY_MSG"
	ENV_MQTT_HOST    = "MQTT_HOST"
	ENV_MQTT_PORT    = "MQTT_PORT"
	ENV_TESLA_STATES = "TESLA_STATES"
	ENV_TESLA_CARID  = "TESLA_CAR_ID"
)

var (
	debug             = false
	mqttHost          = "mosquitto"
	mqttPort          = "1883"
	ntfyURL           = "https://ntfy.sh/"
	ntfyTopic         = "teslas"
	ntfyMsg           = "Car is {{.State}}"
	teslaStates       = []string{"charging"}
	currentTeslaState = ""
	carID             = 1
)

var logger = log.Default()
var tmpl *template.Template

type ntfyTemplate struct {
	State string
}

func main() {
	var err error

	if s, ok := os.LookupEnv(ENV_DEBUG); ok {
		debug = s == "1"
	}

	if s, ok := os.LookupEnv(ENV_NTFY_URL); ok {
		ntfyURL = s
	}

	if s, ok := os.LookupEnv(ENV_NTFY_TOPIC); ok {
		ntfyTopic = s
	}

	if s, ok := os.LookupEnv(ENV_NTFY_MSG); ok {
		ntfyMsg = s
	}

	if s, ok := os.LookupEnv(ENV_TESLA_STATES); ok {
		teslaStates = strings.Split(s, ",")
	}

	if s, ok := os.LookupEnv(ENV_MQTT_HOST); ok {
		mqttHost = s
	}

	if s, ok := os.LookupEnv(ENV_MQTT_PORT); ok {
		mqttPort = s
	}

	if s, ok := os.LookupEnv(ENV_TESLA_CARID); ok {
		carID, err = strconv.Atoi(s)
		if err != nil {
			logger.Println("Error: " + err.Error())
			os.Exit(1)
		}
	}

	if debug {
		logger.Println("Starting")
	}

	// Parse the template from environment
	tmpl, err = template.New("msg").Parse(ntfyMsg)
	if err != nil {
		logger.Println("Error: " + err.Error())
		os.Exit(1)
	}

	svc := teslanotify.Service{
		Debug:         debug,
		Log:           logger,
		MQTTHost:      mqttHost,
		MQTTPort:      mqttPort,
		OnStateChange: onStateChange,
		CarID:         carID,
	}

	err = svc.Connect()
	if err != nil {
		logger.Println("Error: " + err.Error())
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		logger.Printf("Got signal %s", sig.String())
		return
	}
}

// onStateChange handles the different state changes from tesla mate
func onStateChange(s string) error {
	if debug {
		logger.Printf("State=%s\n", s)
	}

	// Skip notification handler if we don't have a state yet
	if currentTeslaState == "" {
		currentTeslaState = s
		return nil
	}

	// Skip notification handler if the state hasn't changed
	if s == currentTeslaState {
		return nil
	}

	// Loop through the wanted states, and notify if the change matches
	for _, ts := range teslaStates {
		if ts == s {
			var msg bytes.Buffer
			err := tmpl.Execute(&msg, ntfyTemplate{State: s})
			if err != nil {
				return err
			}

			err = ntfy(msg.String())
			if err != nil {
				return err
			}
		}
	}

	currentTeslaState = s
	return nil
}

// ntfy sends a notification to a ntfy.sh-instance
func ntfy(msg string) error {
	if debug {
		logger.Printf("Ntfy Msg=%s\n", msg)
	}
	_, err := http.Post(ntfyURL+ntfyTopic, "text/plain", strings.NewReader(msg))
	if err != nil {
		return err
	}

	return nil
}
