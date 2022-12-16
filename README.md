# teslanotify | Receive notifications when car state changes

Simple addon to [TeslaMate](https://github.com/adriankumpf/teslamate) that uses [ntfy.sh](https://ntfy.sh) to notify when car state changes.

## Usage
### Environment Variables
| Variable     | Description | Default |
| ------------ | ----------- | -------- |
| DEBUG        | Outputs more verbose logging if set to `1`. | `0` |
| NTFY_URL     | ntfy instance. | `https://ntfy.sh/` |
| NTFY_TOPIC   | ntfy topic. You want to change this. | `teslas` |
| NTFY_MSG     | Notification message template. This is the message your notification device receives. | `Car is {{.State}}` |
| MQTT_HOST    | The host of the TeslaMate MQTT server. | `moquitto` |
| MQTT_PORT    | The port of the TeslaMate MQTT server. | `1883` |
| TESLA_STATES | Comma separated list of states to notify about. Can be `charging`, `driving`, `suspended`, `online` or `offline` |`charging` |
| TESLA_CAR_ID | The TeslaMate ID of the car. | `1` |

### Setup

To use together with TeslaMate under docker-compose, you can add this to the services section of your docker-compose file. A prebuilt image of this repo is available on Docker Hub.

```yaml
services:

  teslanotify:
    image: hnesland/teslanotify:latest
    restart: always
    environment:
      - NTFY_TOPIC=secret_tesla_topic
      - TESLA_STATES=charging
```

Or you can run the binary anywhere you like, it reads configuration from the environment variables like above.

Build the binary with `go build -o teslanotify ./cmd/teslanotify` and run it like:

```bash
export NTFY_TOPIC=secret_tesla_topic
export MQTT_HOST=127.0.0.1
export MQTT_PORT=1883
export TESLA_STATES=charging
export DEBUG=1
./teslanotify
```
