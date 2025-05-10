package accstatereader

import (
	"github.com/kat-lego/acc-laptime-tracker/pkg/accshm"
	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
	"github.com/kat-lego/acc-laptime-tracker/pkg/utils"
)

type AccStateReader struct{}

func New() *AccStateReader {
	return &AccStateReader{}
}

func (m *AccStateReader) GetState() (*models.AccGameState, error) {
	graphics, err := accshm.ReadSharedMemoryStruct[accshm.Graphics]("Local\\acpmf_graphics")
	if err != nil {
		return nil, err
	}

	info, err := accshm.ReadSharedMemoryStruct[accshm.StaticInfo]("Local\\acpmf_static")
	if err != nil {
		return nil, err
	}

	state := &models.AccGameState{
		Status:              utils.Int32ToAccStatus(graphics.Status),
		SessionType:         utils.Int32ToAccSession(graphics.Session),
		SharedMemoryVersion: utils.Utf16ToString(info.SMVersion[:]),
		AssettoCorsaVersion: utils.Utf16ToString(info.ACVersion[:]),
		Track:               utils.Utf16ToString(info.Track[:]),
		CarModel:            utils.Utf16ToString(info.CarModel[:]),
		SectorCount:         info.SectorCount,
		NumberOfCars:        info.NumCars,
		Clock:               graphics.Clock,
		CompletedLaps:       graphics.CompletedLaps,
		BestLapTime:         graphics.IBestTime,
		PreviousLapTime:     graphics.ILastTime,
		CurrentLapTime:      graphics.ICurrentTime,
		CurrentSectorIndex:  graphics.CurrentSectorIndex,
		PreviousSectorTime:  graphics.LastSectorTime,
		IsValid:             graphics.IsValidLap == 1,
		IsInPitLane:         graphics.IsInPitLane == 1,
		IsInPit:             graphics.IsInPit == 1,
	}

	// jsonBytes, _ := json.MarshalIndent(state, "", "  ")
	// fmt.Println(string(jsonBytes))

	return state, nil
}
