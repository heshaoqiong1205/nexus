package things

type PanTilt struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type PTZSpeed struct {
	PanTilt PanTilt `json:"pan_tilt"`
	Zoom    float32 `json:"zoom"`
}

type SetPresetResponse struct {
	BaseResponse
	PresetID *string `json:"preset_id"`
}

type GetPresetsResponse struct {
	BaseResponse
	Presets []string `json:"presets"`
}

type PresetTourSpot struct {
	PresetID string    `json:"preset_id"`
	Speed    *PTZSpeed `json:"speed"`
	StayTime int       `json:"stay_time"`
}

type PresetTour struct {
	TourID string           `json:"tour_id"`
	Status *string          `json:"status"`
	Spots  []PresetTourSpot `json:"spots"`
}

type GetPresetToursResponse struct {
	BaseResponse
	Tours []PresetTour `json:"tours"`
}

type CreatePresetTourResponse struct {
	BaseResponse
	TourID string `json:"tour_id"`
}

type FloatRange struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

type PTZConfig struct {
	FixedHomePosition      bool       `json:"fixed_home_position"`
	PositionSpace          FloatRange `json:"position_space"`
	ZoomSpace              FloatRange `json:"zoom_space"`
	PanTiltSpeedSpace      PTZSpeed   `json:"pan_tilt_speed_space"`
	ZoomSpeedSpace         PTZSpeed   `json:"zoom_speed_space"`
	MaximumNumberOfPresets int        `json:"maximum_number_of_presets"`
}

type PTZConfigResponse struct {
	BaseResponse
	PTZConfig PTZConfig `json:"ptz_config"`
}

type PTZService struct {
	thingService ThingService
}

func NewIoTService(thingService ThingService) *PTZService {
	return &PTZService{
		thingService: thingService,
	}
}

func (ptz *PTZService) GetPTZConfig(deviceID string, channelIdx int) PTZConfigResponse {
	payload := struct {
		ChannelIdx int `json:"channel_idx"`
	}{
		ChannelIdx: channelIdx,
	}
	var result PTZConfigResponse
	ptz.thingService.RPC(deviceID, "ptz_get_config", payload, &result)
	return result
}

func (ptz *PTZService) ContinuousMove(deviceID string, channelIdx int, velocity PTZSpeed, timeout *int) BaseResponse {
	payload := struct {
		ChannelIdx int      `json:"channel_idx"`
		Velocity   PTZSpeed `json:"velocity"`
		Timeout    *int     `json:"timeout"`
	}{
		ChannelIdx: channelIdx,
		Velocity:   velocity,
		Timeout:    timeout,
	}

	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_get_config", payload, &result)
	return result
}

func (ptz *PTZService) Stop(deviceID string, channelIdx int) BaseResponse {
	payload := struct {
		ChannelIdx int `json:"channel_idx"`
	}{
		ChannelIdx: channelIdx,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_stop", payload, &result)
	return result
}

func (ptz *PTZService) SetPreset(deviceID string, channelIdx int) SetPresetResponse {
	payload := struct {
		ChannelIdx int `json:"channel_idx"`
	}{
		ChannelIdx: channelIdx,
	}
	var result SetPresetResponse
	ptz.thingService.RPC(deviceID, "ptz_set_preset", payload, &result)
	return result
}

func (ptz *PTZService) DeletePreset(deviceID string, channelIdx int, presetID string) BaseResponse {
	payload := struct {
		ChannelIdx int    `json:"channel_idx"`
		PresetID   string `json:"preset_id"`
	}{
		ChannelIdx: channelIdx,
		PresetID:   presetID,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_delete_preset", payload, &result)
	return result
}

func (ptz *PTZService) GetPresets(deviceID string, channelIdx int) GetPresetsResponse {
	payload := struct {
		ChannelIdx int `json:"channel_idx"`
	}{
		ChannelIdx: channelIdx,
	}
	var result GetPresetsResponse
	ptz.thingService.RPC(deviceID, "ptz_get_presets", payload, &result)
	return result
}

func (ptz *PTZService) GoToPreset(deviceID string, channelIdx int, presetID string) BaseResponse {
	payload := struct {
		ChannelIdx int    `json:"channel_idx"`
		PresetID   string `json:"preset_id"`
	}{
		ChannelIdx: channelIdx,
		PresetID:   presetID,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_goto_preset", payload, &result)
	return result
}

func (ptz *PTZService) CreatePresetTour(deviceID string, channelIdx int, spots []PresetTourSpot) CreatePresetTourResponse {
	payload := struct {
		ChannelIdx int              `json:"channel_idx"`
		Spots      []PresetTourSpot `json:"spots"`
	}{
		ChannelIdx: channelIdx,
		Spots:      spots,
	}
	var result CreatePresetTourResponse
	ptz.thingService.RPC(deviceID, "ptz_create_preset_tour", payload, &result)
	return result
}

func (ptz *PTZService) DeletePresetTour(deviceID string, channelIdx int, tourID string) BaseResponse {
	payload := struct {
		ChannelIdx int    `json:"channel_idx"`
		TourID     string `json:"tour_id"`
	}{
		ChannelIdx: channelIdx,
		TourID:     tourID,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_delete_preset_tour", payload, &result)
	return result
}

func (ptz *PTZService) StartPresetTour(deviceID string, channelIdx int, tourID string) BaseResponse {
	payload := struct {
		ChannelIdx int    `json:"channel_idx"`
		TourID     string `json:"tour_id"`
	}{
		ChannelIdx: channelIdx,
		TourID:     tourID,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_start_preset_tour", payload, &result)
	return result
}

func (ptz *PTZService) StopPresetTour(deviceID string, channelIdx int, tourID string) BaseResponse {
	payload := struct {
		ChannelIdx int    `json:"channel_idx"`
		TourID     string `json:"tour_id"`
	}{
		ChannelIdx: channelIdx,
		TourID:     tourID,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_stop_preset_tour", payload, &result)
	return result
}

func (ptz *PTZService) GetPresetTours(deviceID string, channelIdx int) GetPresetToursResponse {
	payload := struct {
		ChannelIdx int `json:"channel_idx"`
	}{
		ChannelIdx: channelIdx,
	}
	var result GetPresetToursResponse
	ptz.thingService.RPC(deviceID, "ptz_get_preset_tours", payload, &result)
	return result
}

func (ptz *PTZService) SetPresetTour(deviceID string, channelIdx int, presetTour PresetTour) BaseResponse {
	payload := struct {
		ChannelIdx int        `json:"channel_idx"`
		PresetTour PresetTour `json:"preset_tour"`
	}{
		ChannelIdx: channelIdx,
		PresetTour: presetTour,
	}
	var result BaseResponse
	ptz.thingService.RPC(deviceID, "ptz_set_preset_tour", payload, &result)
	return result
}
