package config

type Rtc struct {
	Zego *RtcZego `yaml:"zego" json:"zego"`
}

var RtcConfig = new(Rtc)
