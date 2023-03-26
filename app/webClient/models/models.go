package models

type SendCommandInput struct {
	Payload []byte
	Ip      string
	Port    uint16
}

type SendCommandReturn struct {
	Payload []byte
}
