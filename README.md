# BroadlinkAC2MQTT
Control your broadlink-based air conditioner using Home Assistant

![Image](image.png)

## Advantages

* Small application size (~10.2 Mb docker, ~8.2 Mb Windows Standalone)
* Easy to install and use
* Support all platforms
* Parallel independent air conditioning support.
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
      auto_discovery_topic: homeassistant   # Optional
      auto_discovery_topic_retain: false    # Default: true
      certificate_authority: "./config/cert/ca.crt" # Optional. CA certificate in CRT format.
      skip_cert_cn_check: false                     # Default: true. Don’t verify if the common name in the server certificate matches the value of broker.
      
      ## Authorization using client certificates
      certificate_client: "./config/cert/client.crt"  # Optional
      key-client: "./config/cert/client.key"          # Optional
    
    devices:
      - ip: 192.168.1.12
        mac: 34ea345b0fd4   # Only this format is supported
        name: Childroom AC
        port: 80 
      - ip: 192.168.1.18
        mac: 34ea346b0mks   # Only this format is supported
        name: Bedroom AC
        port: 80 

```

## Installation

### Docker Compose

```
    version: '3.5'
    services:
      broadlinkac2mqtt:
        image: "ghcr.io/artemvladimirov/broadlinkac2mqtt:latest"
        container_name: "broadlinkac2mqtt"
        restart: "on-failure"
        volumes:
            - /PATH_TO_YOUR_CONFIG:/config     

```

### Docker

```
   docker run -d --name="broadlinkac2mqtt" -v /PATH_TO_YOUR_CONFIG:/config --restart always ghcr.io/artemvladimirov/broadlinkac2mqtt:latest   
```

### Standalone application

Download application from releases or build it with command "go build". Then you can run a program. The config folder must be located in the program folder

## Support

To motivate the developer, click on the STAR ⭐. I will be very happy!
