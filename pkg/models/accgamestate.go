package models

type AccGameState struct {
	SharedMemoryVersion string `json:"sharedMemoryVersion"`
	AssettoCorsaVersion string `json:"assettoCorsaVersion"`

	Status       string  `json:"status"`
	SessionType  string  `json:"sessionType"`
	Track        string  `json:"track"`
	CarModel     string  `json:"carModel"`
	SectorCount  int32   `json:"sectorCount"`
	NumberOfCars int32   `json:"numberOfCars"`
	Clock        float32 `json:"clock"`

	CompletedLaps   int32 `json:"completedLaps"`
	BestLapTime     int32 `json:"bestLapTime"`
	PreviousLapTime int32 `json:"previousLapTime"`

	CurrentLapTime     int32 `json:"currentLapTime"`
	CurrentSectorIndex int32 `json:"sectorIndex"`
	PreviousSectorTime int32 `json:"previousSectorTime"`
	IsValid            bool  `json:"isValid"`
	IsInPitLane        bool  `json:"isInPitLane"`
	IsInPit            bool  `json:"isInPit"`
}
