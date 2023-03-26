# BroadlinkAC2MQTT
Control your broadlink-based air conditioner using Home Assistant

![Image](image.png)

## Advantages

    * Small application size (~10.2 Mb docker, ~8.2 Mb Windows Standalone)
    * Easy to install and use
    * Support all platforms
    * Parallel independent air conditioning support.  I
      If one air conditioner is offline, it will not affect the rest!

## Configuration

You must specify the mqtt and air conditioner settings in the config.yml file in the config folder.

Example of config.yml 

```
    service:
      update_interval: 10 # In seconds. Default: 10
      log_level: error    # Supported: info, disabled, fatal, debug, error. Default: error
    
    mqtt:
      port: 1883                            # Default: 1883
      host: 192.168.1.36                    # Required
      user: admin                           # Optional  
      password: password                    # Optional    
      client_id: airac                      # Default: broadlinkac
      topic_prefix: aircon                  # Default: airac
      auto_discovery_topic: homeassistant   # Default: homeassistant
      auto_discovery_topic_retain: false    # Default: true
      auto_discovery: true                  # Default: true
    
    devices:
      - ip: 192.168.1.12
        mac: 34ea345b0fd4   # Only this format is supported
        name: Childroom AC
        port: 80 

```

## Installation

### Docker Compose

```
    version: '3.5'
    services:
      broadlinkac2mqtt:
        image: "ghcr.io/artvladimirov/broadlinkac2mqtt:latest"
        container_name: "broadlinkac2mqtt"
        restart: "on-failure"
        volumes:
            - /PATH_TO_YOUR_CONFIG:/config     

```

### Standalone application

Download application from releases or build it with command 

```
    go build
```

## Support

To motivate the developer, click on the STAR ⭐. I will be very happy!