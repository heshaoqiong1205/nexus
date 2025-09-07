package things

import (
	"encoding/json"
	"fmt"
	"time"
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

type Video struct {
	Flip       *bool `json:"flip"`
	OSD        *bool `json:"osd"`
	Brightness *int  `json:"brightness"`
	Sharpness  *int  `json:"sharpness"`
}

func (v *Video) Validate() error {
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

func (v *Video) Equal(other *Video) bool {
	if other == nil {
		return v == nil
	}
	if v == nil {
		return false
	}

	// Check Flip
	if v.Flip == nil || other.Flip == nil {
		if v.Flip != other.Flip {
			return false
		}
	} else if *v.Flip != *other.Flip {
		return false
	}

	// Check OSD
	if v.OSD == nil || other.OSD == nil {
		if v.OSD != other.OSD {
			return false
		}
	} else if *v.OSD != *other.OSD {
		return false
	}

	// Check Brightness
	if v.Brightness == nil || other.Brightness == nil {
		if v.Brightness != other.Brightness {
			return false
		}
	} else if *v.Brightness != *other.Brightness {
		return false
	}

	// Check Sharpness
	if v.Sharpness == nil || other.Sharpness == nil {
		if v.Sharpness != other.Sharpness {
			return false
		}
	} else if *v.Sharpness != *other.Sharpness {
		return false
	}

	return true
}

type Storage struct {
	Mode     *string `json:"mode"`
	Capacity *int    `json:"capacity"`
	Status   *bool   `json:"status"`
}

func (s *Storage) Validate() error {
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

func (s *Storage) Equal(other *Storage) bool {
	if other == nil {
		return s == nil
	}
	if s == nil {
		return false
	}

	// Check Mode
	if s.Mode == nil || other.Mode == nil {
		if s.Mode != other.Mode {
			return false
		}
	} else if *s.Mode != *other.Mode {
		return false
	}

	// Check Capacity
	if s.Capacity == nil || other.Capacity == nil {
		if s.Capacity != other.Capacity {
			return false
		}
	} else if *s.Capacity != *other.Capacity {
		return false
	}

	// Check Status
	if s.Status == nil || other.Status == nil {
		if s.Status != other.Status {
			return false
		}
	} else if *s.Status != *other.Status {
		return false
	}

	return true
}

type Record struct {
	Mode     *string `json:"mode"`
	Duration *int    `json:"duration"`
}

func (r *Record) Validate() error {
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

func (r *Record) Equal(other *Record) bool {
	if other == nil {
		return r == nil
	}
	if r == nil {
		return false
	}

	// Check Mode
	if r.Mode == nil || other.Mode == nil {
		if r.Mode != other.Mode {
			return false
		}
	} else if *r.Mode != *other.Mode {
		return false
	}

	// Check Duration
	if r.Duration == nil || other.Duration == nil {
		if r.Duration != other.Duration {
			return false
		}
	} else if *r.Duration != *other.Duration {
		return false
	}

	return true
}

type VMD struct {
	Status      *bool `json:"status"`
	Sensitivity *int  `json:"sensitivity"`
	Area        []int `json:"area"`
}

func (v *VMD) Validate() error {
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

func (v *VMD) Equal(other *VMD) bool {
	if other == nil {
		return v == nil
	}
	if v == nil {
		return false
	}

	// Check Status
	if v.Status == nil || other.Status == nil {
		if v.Status != other.Status {
			return false
		}
	} else if *v.Status != *other.Status {
		return false
	}

	// Check Sensitivity
	if v.Sensitivity == nil || other.Sensitivity == nil {
		if v.Sensitivity != other.Sensitivity {
			return false
		}
	} else if *v.Sensitivity != *other.Sensitivity {
		return false
	}

	// Check Area
	if len(v.Area) != len(other.Area) {
		return false
	}
	for i := range v.Area {
		if v.Area[i] != other.Area[i] {
			return false
		}
	}

	return true
}

type Detection struct {
	Status      *bool `json:"status"`
	Sensitivity *int  `json:"sensitivity"`
}

func (d *Detection) Validate() error {
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

func (d *Detection) Equal(other *Detection) bool {
	if other == nil {
		return d == nil
	}
	if d == nil {
		return false
	}

	// Check Status
	if d.Status == nil || other.Status == nil {
		if d.Status != other.Status {
			return false
		}
	} else if *d.Status != *other.Status {
		return false
	}

	// Check Sensitivity
	if d.Sensitivity == nil || other.Sensitivity == nil {
		if d.Sensitivity != other.Sensitivity {
			return false
		}
	} else if *d.Sensitivity != *other.Sensitivity {
		return false
	}

	return true
}

type Point struct {
	SN *int `json:"sn"`
	X  *int `json:"x"`
	Y  *int `json:"y"`
}

func (p *Point) Validate() error {
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

func (p *Point) Equal(other *Point) bool {
	if other == nil {
		return p == nil
	}
	if p == nil {
		return false
	}

	// Check SN
	if p.SN == nil || other.SN == nil {
		if p.SN != other.SN {
			return false
		}
	} else if *p.SN != *other.SN {
		return false
	}

	// Check X
	if p.X == nil || other.X == nil {
		if p.X != other.X {
			return false
		}
	} else if *p.X != *other.X {
		return false
	}

	// Check Y
	if p.Y == nil || other.Y == nil {
		if p.Y != other.Y {
			return false
		}
	} else if *p.Y != *other.Y {
		return false
	}

	return true
}

type Cruise struct {
	Status      *string `json:"status"`
	PresetPoint []Point  `json:"preset_point"`
	Route       []int    `json:"route"`
}

func (c *Cruise) Validate() error {
	if c.Status == nil {
		return fmt.Errorf("cruise status cannot be nil")
	}
	if c.PresetPoint == nil {
		return fmt.Errorf("cruise preset point cannot be nil")
	}
	if c.Route == nil {
		return fmt.Errorf("cruise route cannot be nil")
	}
	for _, point := range c.PresetPoint {
		if err := point.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cruise) Equal(other *Cruise) bool {
	if other == nil {
		return c == nil
	}
	if c == nil {
		return false
	}

	// Check Status
	if c.Status == nil || other.Status == nil {
		if c.Status != other.Status {
			return false
		}
	} else if *c.Status != *other.Status {
		return false
	}

	// Check PresetPoint
	if len(c.PresetPoint) != len(other.PresetPoint) {
		return false
	}
	for i := range c.PresetPoint {
		if !c.PresetPoint[i].Equal(&other.PresetPoint[i]) {
			return false
		}
	}

	// Check Route
	if len(c.Route) != len(other.Route) {
		return false
	}
	for i := range c.Route {
		if c.Route[i] != other.Route[i] {
			return false
		}
	}

	return true
}

type Siren struct {
	Duration *int `json:"duration"`
	Volume   *int `json:"volume"`
}

func (s *Siren) Validate() error {
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

func (s *Siren) Equal(other *Siren) bool {
	if other == nil {
		return s == nil
	}
	if s == nil {
		return false
	}

	// Check Duration
	if s.Duration == nil || other.Duration == nil {
		if s.Duration != other.Duration {
			return false
		}
	} else if *s.Duration != *other.Duration {
		return false
	}

	// Check Volume
	if s.Volume == nil || other.Volume == nil {
		if s.Volume != other.Volume {
			return false
		}
	} else if *s.Volume != *other.Volume {
		return false
	}

	return true
}

type State struct {
	Video            *Video     `json:"video"`
	Storage          *Storage   `json:"storage"`
	Record           *Record    `json:"record"`
	MotionDetection  *VMD       `json:"motion_detection"`
	DecibelDetection *Detection `json:"decibel_detection"`
	Cruise           *Cruise    `json:"cruise"`
	Siren            *Siren     `json:"siren"`
	Volume           *int       `json:"volume"`
	PrivacyMode      *bool      `json:"privacy_mode"`
	NightVision      *bool      `json:"night_vision"`
	MotionTracking   *bool      `json:"motion_tracking"`
}

func (state *State) Validate() error {
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
	if state.Record != nil {
		if err := state.Record.Validate(); err != nil {
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
	if state.Cruise != nil {
		if err := state.Cruise.Validate(); err != nil {
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

func (state *State) Marshal() ([]byte, error) {
	if err := state.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(state)
}

type DesiredVideo struct {
	Video
	ID        int64  `json:"id"`
	Timestamp int64   `json:"timestamp"`
}

type DesiredStorage struct {
	Storage
	ID        int64  `json:"id"`
	Timestamp int64   `json:"timestamp"`
}

type DesiredRecord struct {
	Record
	ID         int64  `json:"id"`
	Timestamp  int64   `json:"timestamp"`

}

type DesiredVMD struct {
	VMD
	ID         int64  `json:"id"`
	Timestamp  int64   `json:"timestamp"`
}

type DesiredDecibelDetection struct {
	Detection
	ID         int64  `json:"id"`
	Timestamp  int64   `json:"timestamp"`
}

type DesiredCruise struct {
	Cruise
	ID         int64  `json:"id"`
	Timestamp  int64   `json:"timestamp"`
}

type DesiredSiren struct {
	Siren
	ID         int64  `json:"id"`
	Timestamp  int64   `json:"timestamp"`
}

type DesiredVolume struct {
	Value     *int    `json:"value"`
	ID        int64   `json:"id"`
	Timestamp int64   `json:"timestamp"`
}

type DesiredBoolean struct {
	Value     *bool   `json:"value"`
	ID        int64   `json:"id"`
	Timestamp int64   `json:"timestamp"`
}


type DesiredState struct {
	Video            *DesiredVideo     `json:"video"`
	Storage          *DesiredStorage   `json:"storage"`
	Record           *DesiredRecord   `json:"record"`
	MotionDetection  *DesiredVMD       `json:"motion_detection"`
	DecibelDetection *DesiredDecibelDetection `json:"decibel_detection"`
	Cruise           *DesiredCruise   `json:"cruise"`
	Siren            *DesiredSiren     `json:"siren"`
	Volume           *DesiredVolume     `json:"volume"`
	PrivacyMode      *DesiredBoolean    `json:"privacy_mode"`
	NightVision      *DesiredBoolean    `json:"night_vision"`
	MotionTracking   *DesiredBoolean    `json:"motion_tracking"`
}

func (state *DesiredState) Confirm(lastID int64) {
	if state.Video != nil && state.Video.ID <= lastID {
		state.Video = nil
	}
	if state.Storage != nil && state.Storage.ID <= lastID {
		state.Storage = nil
	}
	if state.Record != nil && state.Record.ID <= lastID {
		state.Record = nil
	}
	if state.MotionDetection != nil && state.MotionDetection.ID <= lastID {
		state.MotionDetection = nil
	}
	if state.DecibelDetection != nil && state.DecibelDetection.ID <= lastID {
		state.DecibelDetection = nil
	}
	if state.Cruise != nil && state.Cruise.ID <= lastID {
		state.Cruise = nil
	}
	if state.Siren != nil && state.Siren.ID <= lastID {
		state.Siren = nil
	}
	if state.Volume != nil && state.Volume.ID <= lastID {
		state.Volume = nil
	}
	if state.PrivacyMode != nil && state.PrivacyMode.ID <= lastID {
		state.PrivacyMode = nil
	}
	if state.NightVision != nil && state.NightVision.ID <= lastID {
		state.NightVision = nil
	}
	if state.MotionTracking != nil && state.MotionTracking.ID <= lastID {
		state.MotionTracking = nil
	}
}

func NewDesiredState(startID int64, origin, update *State) (int64, DesiredState) {
	state := DesiredState{}
	ID := startID
	if update.Video != nil && !update.Video.Equal(origin.Video) {
		ID++
		// Assuming DesiredVideo has the same field names as Video
		// Modify this to match the actual DesiredVideo struct definition
		state.Video.Flip = update.Video.Flip
		state.Video.OSD = update.Video.OSD
		state.Video.Brightness = update.Video.Brightness
		state.Video.Sharpness = update.Video.Sharpness
		state.Video.ID =  ID
		state.Video.Timestamp = time.Now().UnixMilli()
	}

	if update.Storage != nil && !update.Storage.Equal(origin.Storage) {
		ID++
		// Assuming DesiredStorage has the same field names as Storage
		// Modify this to match the actual DesiredStorage struct definition
		state.Storage.Capacity = update.Storage.Capacity
		state.Storage.Mode = update.Storage.Mode
		state.Storage.Status = update.Storage.Status
		state.Storage.ID =  ID
		state.Storage.Timestamp = time.Now().UnixMilli()
	}

	if update.MotionDetection != nil && !update.MotionDetection.Equal(origin.MotionDetection) {
		ID++
		// Assuming DesiredMotionDetection has the same field names as MotionDetection
		// Modify this to match the actual DesiredMotionDetection struct definition
		state.MotionDetection.Status = update.MotionDetection.Status
		state.MotionDetection.Sensitivity = update.MotionDetection.Sensitivity
		state.MotionDetection.ID = ID
		state.MotionDetection.Timestamp = time.Now().UnixMilli()
	}

	if update.DecibelDetection != nil && !update.DecibelDetection.Equal(origin.DecibelDetection) {
		ID++
		// Assuming DesiredDecibelDetection has the same field names as DecibelDetection
		// Modify this to match the actual DesiredDecibelDetection struct definition
		state.DecibelDetection.Status = update.DecibelDetection.Status
		state.DecibelDetection.Sensitivity = update.DecibelDetection.Sensitivity
		state.DecibelDetection.ID = ID
		state.DecibelDetection.Timestamp = time.Now().UnixMilli()
	}


	if update.Record != nil && !update.Record.Equal(origin.Record) {
		ID++
		// Assuming DesiredRecord has the same field names as Record
		// Modify this to match the actual DesiredRecord struct definition
		state.Record.Duration = update.Record.Duration
		state.Record.Mode = update.Record.Mode
		state.Record.ID =  ID
		state.Record.Timestamp = time.Now().UnixMilli()
	}

	if update.Cruise != nil && !update.Cruise.Equal(origin.Cruise) {
		ID++
		// Assuming DesiredCruise has the same field names as Cruise
		// Modify this to match the actual DesiredCruise struct definition
		state.Cruise.Status = update.Cruise.Status
		state.Cruise.PresetPoint = update.Cruise.PresetPoint
		state.Cruise.Route = update.Cruise.Route
		state.Cruise.ID =  ID
		state.Cruise.Timestamp = time.Now().UnixMilli()

	}

	if update.DecibelDetection != nil && !update.DecibelDetection.Equal(origin.DecibelDetection) {
		ID++
		// Assuming DesiredDecibelDetection has the same field names as DecibelDetection
		// Modify this to match the actual DesiredDecibelDetection struct definition
		state.DecibelDetection.Status = update.DecibelDetection.Status
		state.DecibelDetection.Sensitivity = update.DecibelDetection.Sensitivity
		state.DecibelDetection.ID = ID
		state.DecibelDetection.Timestamp = time.Now().UnixMilli()
	}

	if update.MotionDetection != nil && !update.MotionDetection.Equal(origin.MotionDetection) {
		ID++
		// Assuming DesiredMotionDetection has the same field names as MotionDetection
		// Modify this to match the actual DesiredMotionDetection struct definition
		state.MotionDetection.Status = update.MotionDetection.Status
		state.MotionDetection.Sensitivity = update.MotionDetection.Sensitivity
		state.MotionDetection.ID = ID
		state.MotionDetection.Timestamp = time.Now().UnixMilli()
	}

	if update.MotionTracking != nil && (*update.MotionTracking != *origin.MotionTracking) {
		ID++
		// Assuming DesiredMotionTracking has the same field names as MotionTracking
		// Modify this to match the actual DesiredMotionTracking struct definition
		state.MotionTracking.Value = update.MotionTracking
		state.MotionTracking.ID = ID
		state.MotionTracking.Timestamp = time.Now().UnixMilli()
	}

	if update.NightVision != nil && (*update.NightVision != *origin.NightVision) {
		ID++
		// Assuming DesiredNightVision has the same field names as NightVision
		// Modify this to match the actual DesiredNightVision struct definition
		state.NightVision.Value = update.NightVision
		state.NightVision.ID = ID
		state.NightVision.Timestamp = time.Now().UnixMilli()
	}

	if update.PrivacyMode != nil && (*update.PrivacyMode != *origin.PrivacyMode) {
		ID++
		// Assuming DesiredPrivacyMode has the same field names as PrivacyMode
		// Modify this to match the actual DesiredPrivacyMode struct definition
		state.PrivacyMode.Value = update.PrivacyMode
		state.PrivacyMode.ID = ID
		state.PrivacyMode.Timestamp = time.Now().UnixMilli()
	}
	return ID, state
}

func (state *DesiredState) Marshal() ([]byte, error) {
	return json.Marshal(state)
}
