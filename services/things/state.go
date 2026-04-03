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

type StateMeta struct {
	Epoch     *int64 `json:"epoch,omitempty"`
	ID        *int64 `json:"id,omitempty"`
	Timestamp *int64 `json:"timestamp,omitempty"`
}

type Video struct {
	Flip       *bool `json:"flip"`
	OSD        *bool `json:"osd"`
	Brightness *int  `json:"brightness"`
	Sharpness  *int  `json:"sharpness"`
	StateMeta
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
	StateMeta
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
	StateMeta
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
	StateMeta
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
	StateMeta
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
	PresetPoint []Point `json:"preset_point"`
	Route       []int   `json:"route"`
	StateMeta
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
	StateMeta
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

type IntValue struct {
	Value *int `json:"value"`
	StateMeta
}

type BoolValue struct {
	Value *bool `json:"value"`
	StateMeta
}

type State struct {
	Video            *Video      `json:"video"`
	Storage          *Storage    `json:"storage"`
	Record           *Record     `json:"record"`
	MotionDetection  *VMD        `json:"motion_detection"`
	DecibelDetection *Detection  `json:"decibel_detection"`
	Cruise           *Cruise     `json:"cruise"`
	Siren            *Siren      `json:"siren"`
	Volume           *IntValue   `json:"volume"`
	PrivacyMode      *BoolValue  `json:"privacy_mode"`
	NightVision      *BoolValue  `json:"night_vision"`
	MotionTracking   *BoolValue  `json:"motion_tracking"`
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
	if state.Volume != nil && state.Volume.Value != nil {
		if *state.Volume.Value < 0 || *state.Volume.Value > 100 {
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
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredStorage struct {
	Storage
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredRecord struct {
	Record
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredVMD struct {
	VMD
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredDecibelDetection struct {
	Detection
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredCruise struct {
	Cruise
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredSiren struct {
	Siren
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredVolume struct {
	Value     *int  `json:"value"`
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredBoolean struct {
	Value     *bool `json:"value"`
	Epoch     int64 `json:"epoch"`
	ID        int64 `json:"id"`
	Timestamp int64 `json:"timestamp"`
}

type DesiredState struct {
	Video            *DesiredVideo            `json:"video"`
	Storage          *DesiredStorage          `json:"storage"`
	Record           *DesiredRecord           `json:"record"`
	MotionDetection  *DesiredVMD              `json:"motion_detection"`
	DecibelDetection *DesiredDecibelDetection `json:"decibel_detection"`
	Cruise           *DesiredCruise           `json:"cruise"`
	Siren            *DesiredSiren            `json:"siren"`
	Volume           *DesiredVolume           `json:"volume"`
	PrivacyMode      *DesiredBoolean          `json:"privacy_mode"`
	NightVision      *DesiredBoolean          `json:"night_vision"`
	MotionTracking   *DesiredBoolean          `json:"motion_tracking"`
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func intPtrEqual(a, b *int) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func intValue(value *IntValue) *int {
	if value == nil {
		return nil
	}
	return value.Value
}

func boolValue(value *BoolValue) *bool {
	if value == nil {
		return nil
	}
	return value.Value
}

func cloneVideo(video *Video) Video {
	if video == nil {
		return Video{}
	}
	return Video{
		Flip:       video.Flip,
		OSD:        video.OSD,
		Brightness: video.Brightness,
		Sharpness:  video.Sharpness,
	}
}

func cloneStorage(storage *Storage) Storage {
	if storage == nil {
		return Storage{}
	}
	return Storage{
		Mode:     storage.Mode,
		Capacity: storage.Capacity,
		Status:   storage.Status,
	}
}

func cloneRecord(record *Record) Record {
	if record == nil {
		return Record{}
	}
	return Record{
		Mode:     record.Mode,
		Duration: record.Duration,
	}
}

func cloneVMD(vmd *VMD) VMD {
	if vmd == nil {
		return VMD{}
	}
	area := append([]int(nil), vmd.Area...)
	return VMD{
		Status:      vmd.Status,
		Sensitivity: vmd.Sensitivity,
		Area:        area,
	}
}

func cloneDetection(detection *Detection) Detection {
	if detection == nil {
		return Detection{}
	}
	return Detection{
		Status:      detection.Status,
		Sensitivity: detection.Sensitivity,
	}
}

func cloneCruise(cruise *Cruise) Cruise {
	if cruise == nil {
		return Cruise{}
	}
	presetPoints := append([]Point(nil), cruise.PresetPoint...)
	route := append([]int(nil), cruise.Route...)
	return Cruise{
		Status:      cruise.Status,
		PresetPoint: presetPoints,
		Route:       route,
	}
}

func cloneSiren(siren *Siren) Siren {
	if siren == nil {
		return Siren{}
	}
	return Siren{
		Duration: siren.Duration,
		Volume:   siren.Volume,
	}
}

func nextDesiredVersion(currentEpoch, currentID int64) (int64, int64) {
	if currentID == int64(^uint64(0)>>1) {
		return currentEpoch + 1, 0
	}
	return currentEpoch, currentID + 1
}

func NewDesiredState(startID int64, origin, update *State) (int64, DesiredState) {
	nextID, _, state := NewDesiredStateWithEpoch(0, startID, update)
	return nextID, state
}

func NewDesiredStateWithEpoch(startEpoch, startID int64, update *State) (int64, int64, DesiredState) {
	state := DesiredState{}
	epoch := startEpoch
	id := startID
	now := time.Now().UnixMilli()

	if update.Video != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.Video = &DesiredVideo{
			Video:     cloneVideo(update.Video),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.Storage != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.Storage = &DesiredStorage{
			Storage:   cloneStorage(update.Storage),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.MotionDetection != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.MotionDetection = &DesiredVMD{
			VMD:       cloneVMD(update.MotionDetection),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.DecibelDetection != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.DecibelDetection = &DesiredDecibelDetection{
			Detection: cloneDetection(update.DecibelDetection),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.Record != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.Record = &DesiredRecord{
			Record:    cloneRecord(update.Record),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.Cruise != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.Cruise = &DesiredCruise{
			Cruise:    cloneCruise(update.Cruise),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.Siren != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.Siren = &DesiredSiren{
			Siren:     cloneSiren(update.Siren),
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.Volume != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.Volume = &DesiredVolume{
			Value:     update.Volume.Value,
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.MotionTracking != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.MotionTracking = &DesiredBoolean{
			Value:     update.MotionTracking.Value,
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.NightVision != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.NightVision = &DesiredBoolean{
			Value:     update.NightVision.Value,
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	if update.PrivacyMode != nil {
		epoch, id = nextDesiredVersion(epoch, id)
		state.PrivacyMode = &DesiredBoolean{
			Value:     update.PrivacyMode.Value,
			Epoch:     epoch,
			ID:        id,
			Timestamp: now,
		}
	}

	return id, epoch, state
}

func (state *DesiredState) MaxEpoch() int64 {
	maxEpoch := int64(0)
	if state.Video != nil && state.Video.Epoch > maxEpoch {
		maxEpoch = state.Video.Epoch
	}
	if state.Storage != nil && state.Storage.Epoch > maxEpoch {
		maxEpoch = state.Storage.Epoch
	}
	if state.Record != nil && state.Record.Epoch > maxEpoch {
		maxEpoch = state.Record.Epoch
	}
	if state.MotionDetection != nil && state.MotionDetection.Epoch > maxEpoch {
		maxEpoch = state.MotionDetection.Epoch
	}
	if state.DecibelDetection != nil && state.DecibelDetection.Epoch > maxEpoch {
		maxEpoch = state.DecibelDetection.Epoch
	}
	if state.Cruise != nil && state.Cruise.Epoch > maxEpoch {
		maxEpoch = state.Cruise.Epoch
	}
	if state.Siren != nil && state.Siren.Epoch > maxEpoch {
		maxEpoch = state.Siren.Epoch
	}
	if state.Volume != nil && state.Volume.Epoch > maxEpoch {
		maxEpoch = state.Volume.Epoch
	}
	if state.PrivacyMode != nil && state.PrivacyMode.Epoch > maxEpoch {
		maxEpoch = state.PrivacyMode.Epoch
	}
	if state.NightVision != nil && state.NightVision.Epoch > maxEpoch {
		maxEpoch = state.NightVision.Epoch
	}
	if state.MotionTracking != nil && state.MotionTracking.Epoch > maxEpoch {
		maxEpoch = state.MotionTracking.Epoch
	}
	return maxEpoch
}

func shouldDeleteDesired(reportedEpoch, reportedID, desiredEpoch, desiredID int64) bool {
	if reportedID > desiredID {
		return true
	}
	if reportedID == desiredID {
		return true
	}
	return reportedEpoch != desiredEpoch
}

func (state *DesiredState) Marshal() ([]byte, error) {
	return json.Marshal(state)
}
