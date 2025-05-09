package things

type OSDService struct {
	thingService ThingService
}

type color struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type OSDConfig struct {
	Enabled         bool    `json:"enabled"`
	Position        string  `json:"position"`
	FontSize        *int    `json:"font_size"`
	FontColor       *color  `json:"font_color"`
	DateFormat      *string `json:"date_format"`
	TimeFormat      *string `json:"time_format"`
	BackgroundColor *color  `json:"background_color"`
	FontStyle       *string `json:"font_style"`
}

type GetOSDConfigResponse struct {
	BaseResponse
	OSDConfig *OSDConfig `json:"osd_config"`
}

func NewOSDService(thingService ThingService) *OSDService {
	return &OSDService{
		thingService: thingService,
	}
}

func (osd *OSDService) SetOSD(deviceID string, channelIdx int, config OSDConfig) BaseResponse {
	payload := struct {
		ChannelIdx int       `json:"channel_idx"`
		Config     OSDConfig `json:"osd"`
	}{
		ChannelIdx: channelIdx,
		Config:     config,
	}

	var result BaseResponse
	osd.thingService.RPC(deviceID, "set_osd", payload, &result)
	return result
}

func (osd *OSDService) GetOSD(deviceID string, channelIdx int) GetOSDConfigResponse {
	payload := struct {
		ChannelIdx int `json:"channel_idx"`
	}{
		ChannelIdx: channelIdx,
	}

	var result GetOSDConfigResponse
	osd.thingService.RPC(deviceID, "get_osd", payload, &result)
	return result
}
