package models

type Session struct {
	Id              string `json:"id"`
	StartTime       int64  `json:"startTime"`
	SessionType     string `json:"sessionType"`
	Track           string `json:"track"`
	CarModel        string `json:"carModel"`
	NumberOfSectors int32  `json:"numberOfSectors"`
	CompletedLaps   int32  `json:"lapNumber"`
	BestLapTime     int32  `json:"bestLapTime"`
	PreviousLapTime int32  `json:"previousLapTime"`
	IsActive        bool   `json:"isActive"`
	Player          string `json:"player"`

	Laps []*Lap `json:"laps"`
}

type Lap struct {
	LapNumber int32 `json:"lapNumber"`
	LapTime   int32 `json:"lapTime"`
	IsValid   bool  `json:"isValid"`
	IsActive  bool  `json:"isActive"`

	LapSectors []*LapSector `json:"lapSectors"`
}

type LapSector struct {
	SectorNumber int32 `json:"sectorNumber"`
	SectorTime   int32 `json:"sectorTime"`
	IsActive     bool  `json:"isActive"`
}
