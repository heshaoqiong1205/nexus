package things

import (
	"encoding/json"
	"fmt"
)

type SubState interface {
	Validate() error
}

type location struct {
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func (l *location) Validate() error {
	if l.Latitude == nil || l.Longitude == nil {
		return fmt.Errorf("location latitude and longitude cannot be nil")
	}
	return nil
}

type video struct {
	Flip       *bool `json:"flip"`
	OSD        *bool `json:"osd"`
	Brightness *int  `json:"brightness"`
	Sharpness  *int  `json:"sharpness"`
}

func (v *video) Validate() error {
	if v.Flip == nil {
		return fmt.Errorf("video flip cannot be nil")
	}
	if v.OSD == nil {
		return fmt.Errorf("video osd cannot be nil")
	}
	if v.Brightness == nil {
		return fmt.Errorf("video brightness cannot be nil")
	}
	if v.Sharpness == nil {
		return fmt.Errorf("video sharpness cannot be nil")
	}
	if *v.Brightness < 0 || *v.Brightness > 100 {
		return fmt.Errorf("video brightness must be between 0 and 100")
	}
	return nil
}

type storage struct {
	Mode     *string `json:"mode"`
	Capacity *int    `json:"capacity"`
	Status   *bool   `json:"status"`
}

func (s *storage) Validate() error {
	if s.Mode == nil {
		return fmt.Errorf("storage mode cannot be nil")
	}
	if s.Capacity == nil {
		return fmt.Errorf("storage capacity cannot be nil")
	}
	if s.Status == nil {
		return fmt.Errorf("storage status cannot be nil")
	}
	if *s.Capacity < 0 {
		return fmt.Errorf("storage capacity must be greater than or equal to 0")
	}
	return nil
}

type recored struct {
	Mode     *string `json:"mode"`
	Duration *int    `json:"duration"`
}

func (r *recored) Validate() error {
	if r.Mode == nil {
		return fmt.Errorf("record mode cannot be nil")
	}
	if r.Duration == nil {
		return fmt.Errorf("record duration cannot be nil")
	}
	if *r.Duration < 0 {
		return fmt.Errorf("record duration must be greater than or equal to 0")
	}
	return nil
}

type vmd struct {
	Status      *bool `json:"status"`
	Sensitivity *int  `json:"sensitivity"`
	Area        []int `json:"area"`
}

func (v *vmd) Validate() error {
	if v.Status == nil {
		return fmt.Errorf("motion detection status cannot be nil")
	}
	if v.Sensitivity == nil {
		return fmt.Errorf("motion detection sensitivity cannot be nil")
	}
	if *v.Sensitivity < 0 || *v.Sensitivity > 100 {
		return fmt.Errorf("motion detection sensitivity must be between 0 and 100")
	}
	return nil
}

type detection struct {
	Status      *bool `json:"status"`
	Sensitivity *int  `json:"sensitivity"`
}

func (d *detection) Validate() error {
	if d.Status == nil {
		return fmt.Errorf("detection status cannot be nil")
	}
	if d.Sensitivity == nil {
		return fmt.Errorf("detection sensitivity cannot be nil")
	}
	if *d.Sensitivity < 0 || *d.Sensitivity > 100 {
		return fmt.Errorf("detection sensitivity must be between 0 and 100")
	}
	return nil
}

type point struct {
	SN *int `json:"sn"`
	X  *int `json:"x"`
	Y  *int `json:"y"`
}

func (p *point) Validate() error {
	if p.SN == nil {
		return fmt.Errorf("point sn cannot be nil")
	}
	if p.X == nil {
		return fmt.Errorf("point x cannot be nil")
	}
	if p.Y == nil {
		return fmt.Errorf("point y cannot be nil")
	}
	return nil
}

type crusise struct {
	Status      *string `json:"status"`
	PresetPoint []point `json:"preset_point"`
	Route       []int   `json:"route"`
}

func (c *crusise) Validate() error {
	if c.Status == nil {
		return fmt.Errorf("crusise status cannot be nil")
	}
	if c.PresetPoint == nil {
		return fmt.Errorf("crusise preset point cannot be nil")
	}
	if c.Route == nil {
		return fmt.Errorf("crusise route cannot be nil")
	}
	for _, point := range c.PresetPoint {
		if err := point.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type siren struct {
	Duration *int `json:"duration"`
	Volume   *int `json:"volume"`
}

func (s *siren) Validate() error {
	if s.Duration == nil {
		return fmt.Errorf("siren duration cannot be nil")
	}
	if s.Volume == nil {
		return fmt.Errorf("siren volume cannot be nil")
	}
	if *s.Duration < 0 {
		return fmt.Errorf("siren duration must be greater than or equal to 0")
	}
	if *s.Volume < 0 || *s.Volume > 100 {
		return fmt.Errorf("siren volume must be between 0 and 100")
	}
	return nil
}

type state struct {
	Video            *video     `json:"video"`
	Storage          *storage   `json:"storage"`
	Recored          *recored   `json:"recored"`
	MotionDetection  *vmd       `json:"motion_detection"`
	DecibelDetection *detection `json:"decibel_detection"`
	Crusise          *crusise   `json:"crusise"`
	Siren            *siren     `json:"siren"`
	Volume           *int       `json:"volume"`
	PrivacyMode      *bool      `json:"privacy_mode"`
	NightVision      *bool      `json:"night_vision"`
	MotionTracking   *bool      `json:"motion_tracking"`
}

func (state *state) Validate() error {
	if state.Video != nil {
		if err := state.Video.Validate(); err != nil {
			return err
		}
	}
	if state.Storage != nil {
		if err := state.Storage.Validate(); err != nil {
			return err
		}
	}
	if state.Recored != nil {
		if err := state.Recored.Validate(); err != nil {
			return err
		}
	}
	if state.MotionDetection != nil {
		if err := state.MotionDetection.Validate(); err != nil {
			return err
		}
	}
	if state.DecibelDetection != nil {
		if err := state.DecibelDetection.Validate(); err != nil {
			return err
		}
	}
	if state.Crusise != nil {
		if err := state.Crusise.Validate(); err != nil {
			return err
		}
	}
	if state.Siren != nil {
		if err := state.Siren.Validate(); err != nil {
			return err
		}
	}
	if state.Volume != nil {
		if *state.Volume < 0 || *state.Volume > 100 {
			return fmt.Errorf("volume must be between 0 and 100")
		}
	}
	return nil
}

func (state *state) Marshal() ([]byte, error) {
	if err := state.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(state)
}
