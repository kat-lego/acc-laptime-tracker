package models

type Session struct {
	Id              uint   `gorm:"primaryKey" json:"id"`
	StartTime       int64  `json:"startTime"`
	SessionType     string `json:"sessionType"`
	Track           string `json:"track"`
	CarModel        string `json:"carModel"`
	CompletedLaps   int32  `json:"lapNumber"`
	BestLapTime     int32  `json:"bestLapTime"`
	PreviousLapTime int32  `json:"previousLapTime"`

	Laps []Lap `gorm:"foreignKey:SessionId"`
}

type Lap struct {
	Id        uint `gorm:"primaryKey" json:"id"`
	SessionId uint `json:"sessionId"`

	LapNumber int32 `json:"lapNumber"`
	LapTime   int32 `json:"lapTime"`
	IsValid   bool  `json:"isValid"`

	LapSectors []LapSector `gorm:"foreignKey:LapId"`
}

type LapSector struct {
	Id    uint `gorm:"primaryKey" json:"id"`
	LapId uint `json:"lapId"`

	SectorNumber int32 `json:"sectorNumber"`
	SectorTime   int32 `json:"sectorTime"`
}
