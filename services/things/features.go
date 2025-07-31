package things

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

type ResolutionRatio struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

func (r *ResolutionRatio) validate() error {
	if r.Width <= 0 {
		return errors.New("invalid resolution ratio width")
	}
	if r.Height <= 0 {
		return errors.New("invalid resolution ratio height")
	}
	return nil
}

// video feature
// The sampling rate is fixed at 90,000
type VideoFeature struct {
	Num              int               `json:"num"`
	Codecs           []string          `json:"codecs"`
	ResolutionRatios []ResolutionRatio `json:"resolution_ratios"`
	Streams          []int             `json:"streams"`
}

func (v *VideoFeature) validate() error {
	log.Printf("Validating video feature: %+v\n", v)
	if v.Num <= 0 {
		return errors.New("invalid video feature num")
	}
	if len(v.Codecs) == 0 {
		return errors.New("invalid video feature codecs")
	}
	if len(v.ResolutionRatios) == 0 {
		return errors.New("invalid video feature resolution ratios")
	}
	for _, codec := range v.Codecs {
		if codec != "H264" && codec != "H265" && codec != "VP8" && codec != "VP9" {
			return errors.New("invalid video feature codec: " + codec)
		}
	}
	for _, ratio := range v.ResolutionRatios {
		if err := ratio.validate(); err != nil {
			return err
		}
	}
	return nil
}

type AudioFeature struct {
	Codecs       []string `json:"codecs"`
	SamplingRate int      `json:"sampling_rate"`
	Channels     int      `json:"channels"`
	Bits         int      `json:"bits"`
}

func (a *AudioFeature) validate() error {
	if a.SamplingRate != 8000 && a.SamplingRate != 16000 && a.SamplingRate != 32000 && a.SamplingRate != 44100 && a.SamplingRate != 48000 {
		return errors.New("invalid audio feature sampling rate: " + strconv.Itoa(a.SamplingRate))
	}
	if a.Channels <= 0 {
		return errors.New("invalid audio feature channels: " + strconv.Itoa(a.Channels))
	}
	if a.Bits <= 0 {
		return errors.New("invalid audio feature bits: " + strconv.Itoa(a.Bits))
	}
	if len(a.Codecs) == 0 {
		return errors.New("audio feature codecs is empty")
	}
	for _, codec := range a.Codecs {
		if codec != "PCMA" && codec != "PCMU" && codec != "PCM" && codec != "OPUS" && codec != "AAC" {
			return errors.New("invalid audio feature codec: " + codec)
		}
	}
	return nil
}

type Features struct {
	P2P          *string       `json:"p2p"`
	WebRTC       []string      `json:"webrtc"`
	UPNP         *string       `json:"upnp"`
	AI           *string       `json:"ai"`
	VideoFeature *VideoFeature `json:"video_feature"`
	AudioFeature *AudioFeature `json:"audio_feature"`
}

func (f *Features) Validate() error {
	if f.VideoFeature != nil {
		if err := f.VideoFeature.validate(); err != nil {
			return err
		}
	}
	if f.AudioFeature != nil {
		if err := f.AudioFeature.validate(); err != nil {
			return err
		}
	}
	if f.P2P != nil {
		if *f.P2P != "standard" && *f.P2P != "enhance" {
			return errors.New("invalid p2p feature: " + *f.P2P)
		}
	}
	if len(f.WebRTC) == 0 {
		return errors.New("invalid webrtc feature")
	}
	for _, webrtc := range f.WebRTC {
		if webrtc != "SRTP" && webrtc != "DC" {
			return errors.New("invalid webrtc feature: " + webrtc)
		}
	}
	if f.AI != nil {
		if *f.AI != "local" && *f.AI != "remote" {
			return errors.New("invalid ai feature: " + *f.AI)
		}
	}
	if f.UPNP != nil {
		if *f.UPNP != "IDGV1" && *f.UPNP != "IDGV2" {
			return errors.New("invalid upnp feature: " + *f.UPNP)
		}
	}
	return nil
}

func (f *Features) Marshal() ([]byte, error) {
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(f)
}

func (f *Features) ValidateRequirements(requiredFeatures []string) error {
	for _, feature := range requiredFeatures {
		switch feature {
		case "p2p":
			if f.P2P == nil {
				return errors.New("missing p2p feature")
			}
		case "webrtc":
			if len(f.WebRTC) == 0 {
				return errors.New("missing webrtc feature")
			}
		case "upnp":
			if f.UPNP == nil {
				return errors.New("missing upnp feature")
			}
		case "ai":
			if f.AI == nil {
				return errors.New("missing ai feature")
			}
		default:
			return errors.New("unknown feature: " + feature)
		}
	}
	return nil
}
