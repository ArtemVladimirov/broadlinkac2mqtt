package models

const (
	DeviceClassClimate string = "climate"
	DeviceClassSwitch  string = "switch"
)

type ConfigMqtt struct {
	Broker                   string
	User                     *string
	Password                 *string
	ClientId                 string
	TopicPrefix              string
	AutoDiscoveryTopic       *string
	AutoDiscoveryTopicRetain bool
}

type ClimateDiscoveryTopic struct {
	FanModeCommandTopic     string                     `json:"fan_mode_command_topic" example:"aircon/34ea345b0fd4/fan_mode/set"`
	SwingModeCommandTopic   string                     `json:"swing_mode_command_topic" example:"aircon/34ea345b0fd4/swing_mode/set"`
	SwingModes              []string                   `json:"swing_modes"` // 'on' 'off'
	TempStep                float32                    `json:"temp_step" example:"0.5"`
	TemperatureStateTopic   string                     `json:"temperature_state_topic" example:"aircon/34ea345b0fd4/temp/value"`
	TemperatureCommandTopic string                     `json:"temperature_command_topic" example:"aircon/34ea345b0fd4/temp/set"`
	Precision               float32                    `json:"precision" example:"0.5"`
	CurrentTemperatureTopic string                     `json:"current_temperature_topic" example:"aircon/34ea345b0fd4/current_temp/value"` // Temperature in the room
	Device                  DiscoveryTopicDevice       `json:"device"`
	ModeCommandTopic        string                     `json:"mode_command_topic" example:"aircon/34ea345b0fd4/mode/set"`
	ModeStateTopic          string                     `json:"mode_state_topic" example:"aircon/34ea345b0fd4/mode/value"`
	Modes                   []string                   `json:"modes"` // [“auto”, “off”, “cool”, “heat”, “dry”, “fan_only”]
	Name                    *string                    `json:"name"`
	FanModes                []string                   `json:"fan_modes"` // : [“auto”, “low”, “medium”, “high”]
	SwingModeStateTopic     string                     `json:"swing_mode_state_topic" example:"aircon/34ea345b0fd4/swing_mode/value"`
	FanModeStateTopic       string                     `json:"fan_mode_state_topic" example:"aircon/34ea345b0fd4/fan_mode/value"`
	UniqueId                string                     `json:"unique_id" example:"34ea345b0fd4"`
	MaxTemp                 float32                    `json:"max_temp" example:"32.0"`
	MinTemp                 float32                    `json:"min_temp" example:"16.0"`
	Availability            DiscoveryTopicAvailability `json:"availability"`
	Icon                    string                     `json:"icon"`
	TemperatureUnit         string                     `json:"temperature_unit"` // C or F
}

type SwitchDiscoveryTopic struct {
	Device       DiscoveryTopicDevice       `json:"device"`
	Name         string                     `json:"name" example:"childroom"`
	UniqueId     string                     `json:"unique_id" example:"34ea345b0fd4"`
	StateTopic   string                     `json:"state_topic" example:"aircon/34ea345b0fd4/display/switch"`
	CommandTopic string                     `json:"command_topic" example:"aircon/34ea345b0fd4/display/switch/set"`
	Availability DiscoveryTopicAvailability `json:"availability"`
	Icon         string                     `json:"icon"`
}

type DiscoveryTopicDevice struct {
	Model string `json:"model" example:"Aircon"`
	Mf    string `json:"mf" example:"Broadlink"`
	Sw    string `json:"sw" example:"1.1.3"`
	Ids   string `json:"ids" example:"34ea345b0fd4"`
	Name  string `json:"name" example:"childroom"`
}

type DiscoveryTopicAvailability struct {
	PayloadAvailable    string `json:"payload_available" example:"online"`
	PayloadNotAvailable string `json:"payload_not_available" example:"offline"`
	Topic               string `json:"topic" example:"aircon/34ea345b0fd4/availability/value"`
}

type PublishClimateDiscoveryTopicInput struct {
	Topic ClimateDiscoveryTopic
}

type PublishSwitchDiscoveryTopicInput struct {
	Topic SwitchDiscoveryTopic
}

type PublishAmbientTempInput struct {
	Mac         string
	Temperature float32
}

type PublishTemperatureInput struct {
	Mac         string
	Temperature float32
}

type PublishModeInput struct {
	Mac  string
	Mode string
}

type PublishSwingModeInput struct {
	Mac       string
	SwingMode string
}

type PublishFanModeInput struct {
	Mac     string
	FanMode string
}

type PublishAvailabilityInput struct {
	Mac          string
	Availability string
}
type PublishDisplaySwitchInput struct {
	Mac    string
	Status string
}
