package data

import (
	"errors"
	"time"
)

// https://threedots.tech/post/safer-enums-in-go/
type ChannelType struct {
	channelType string
}

var (
	Unknown  = ChannelType{channelType: ""}
	WhatsApp = ChannelType{channelType: "whatsapp"}
)

func (cht ChannelType) String() string {
	return cht.channelType
}

func ChannelTypeFromString(s string) (ChannelType, error) {
	switch s {
	case WhatsApp.channelType:
		return WhatsApp, nil
	}

	return Unknown, errors.New("unknown channel type: " + s)
}

// use extension interface pattern to replace inheritance with composition
// see https://medium.com/swlh/what-is-the-extension-interface-pattern-in-golang-ce852dcecaec
type TenantChannelConfig struct {
	Channel ChannelType
}

type ChannelConfigWhatsAppNumbers struct {
	PhoneNumber string
	Language    string
}

type TenantChannelConfigWhatsApp struct {
	TenantChannelConfig
	AccountSid string
	Numbers    []ChannelConfigWhatsAppNumbers
}

// represents configuration of platform tenant
// for now this means tenant channels, later NLPs, human providers, etc.
type TenantConfig struct {
	TenantId string
	Name     string
	Desc     string
	Channels []TenantChannelConfig
}

type CustomerProfile struct {
	CustomerId  string
	PhoneNumber string
	Email       string
	Name        string
	Surname     string
}

type CustomerSession struct {
	SessionId  string
	CustomerId string
	Channel    ChannelType
}

type SessionMessage struct {
	SessionId   string
	CreatedAt   time.Time
	MessageText string
}
