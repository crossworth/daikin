package types

import (
	"fmt"
)

type Mode int

func (m Mode) String() string {
	switch m {
	case 0:
		return "Automático"
	case 2:
		return "Desumidificar"
	case 3:
		return "Resfriar"
	case 4:
		return "Aquecer"
	case 6:
		return "Ventilar"
	default:
		return fmt.Sprintf("%d", m)
	}
}

type Fan int

func (f Fan) String() string {
	switch f {
	case 3:
		return "Baixa"
	case 4:
		return "Média-Baixa"
	case 5:
		return "Média"
	case 6:
		return "Média-Alta"
	case 7:
		return "Alta"
	case 17:
		return "Automático"
	case 18:
		return "Silencioso"
	default:
		return fmt.Sprintf("%d", f)
	}
}

type Port struct {
	Power         int     `json:"power"`
	Mode          Mode    `json:"mode"`
	Temperature   float64 `json:"temperature"`
	Fan           Fan     `json:"fan"`
	HSwing        int     `json:"h_swing"`
	VSwing        int     `json:"v_swing"`
	Coanda        int     `json:"coanda"`
	Econo         int     `json:"econo"`
	Powerchill    int     `json:"powerchill"`
	GoodSleep     int     `json:"good_sleep"`
	Streamer      int     `json:"streamer"`
	OutQuite      int     `json:"out_quite"`
	OnTimerSet    int     `json:"on_timer_set"`
	OnTimerValue  int     `json:"on_timer_value"`
	OffTimerSet   int     `json:"off_timer_set"`
	OffTimerValue int     `json:"off_timer_value"`
	Sensors       Sensors `json:"sensors"`
	RstR          int     `json:"rst_r"`
	FWVer         string  `json:"fw_ver"`
}

type Sensors struct {
	RoomTemp float64 `json:"room_temp"`
	OutTemp  float64 `json:"out_temp"`
}

type PortState struct {
	Power       *int     `json:"power,omitempty"`
	Mode        *Mode    `json:"mode,omitempty"`
	Temperature *float64 `json:"temperature,omitempty"`
	Fan         *Fan     `json:"fan,omitempty"`
	VSwing      *int     `json:"v_swing,omitempty"`
	Coanda      *int     `json:"coanda,omitempty"`
	Econo       *int     `json:"econo,omitempty"`
	Powerchill  *int     `json:"powerchill,omitempty"`
}
