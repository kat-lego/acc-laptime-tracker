package streams

import (
	"time"

	"github.com/kat-lego/acc-laptime-tracker/pkg/accshm"
	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
	"github.com/kat-lego/acc-laptime-tracker/pkg/utils"
)

const (
	lapRecordBufSize = 2
)

type LapRecordStream struct {
	lapRecordBuf       [lapRecordBufSize]*models.LapRecord
	nextBufferPos      int32
	seenFirstSession   bool
	seenFirstLap       bool
	seenFirstLapSector bool
}

func NewLapRecordStream() *LapRecordStream {
	return &LapRecordStream{
		nextBufferPos:      0,
		seenFirstSession:   false,
		seenFirstLap:       false,
		seenFirstLapSector: false,
	}
}

func (m *LapRecordStream) GetNextSession() *models.Session {
	lr := m.getLatestLapRecord()
	sr := m.getShadowLapRecord()

	if m.seenFirstSession && ((lr.CompletedLaps >= sr.CompletedLaps) ||
		(lr.Clock >= sr.Clock)) {
		return nil
	}

	session := models.Session{
		StartTime:   time.Now().UTC().Unix(),
		SessionType: lr.SessionType,
		Track:       lr.Track,
		CarModel:    lr.CarModel,
	}

	m.seenFirstSession = true

	return &session
}

func (m *LapRecordStream) GetNextLap() *models.Lap {
	lr := m.getLatestLapRecord()
	sr := m.getShadowLapRecord()

	if m.seenFirstLap && lr.CompletedLaps <= sr.CompletedLaps {
		return nil
	}

	lap := models.Lap{
		LapNumber: lr.CompletedLaps + 1,
		IsValid:   lr.IsValid,
	}

	m.seenFirstLap = true
	return &lap
}

func (m *LapRecordStream) GetNextLapSector() *models.LapSector {
	lr := m.getLatestLapRecord()
	sr := m.getShadowLapRecord()

	if m.seenFirstLapSector && lr.CurrentSectorIndex == sr.CurrentSectorIndex {
		return nil
	}

	lapSector := models.LapSector{
		SectorNumber: lr.CurrentSectorIndex + 1,
	}

	m.seenFirstLapSector = true
	return &lapSector
}

func (m *LapRecordStream) RefreshBuffer() error {
	graphics, err := accshm.ReadSharedMemoryStruct[accshm.Graphics]("Local\\acpmf_graphics")
	if err != nil {
		// fmt.Printf("Failed to read graphics struct: %v\n", err)
		return err
	}

	info, err := accshm.ReadSharedMemoryStruct[accshm.StaticInfo]("Local\\acpmf_static")
	if err != nil {
		// fmt.Printf("Failed to read static info struct: %v\n", err)
		return err
	}

	m.lapRecordBuf[m.nextBufferPos] = &models.LapRecord{
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
		CurrentSectorIndex:  graphics.CurrentSectorIndex,
		PreviousSectorTime:  graphics.LastSectorTime,
		IsValid:             graphics.IsValidLap == 1,
		IsInPitLane:         graphics.IsInPitLane == 1,
		IsInPit:             graphics.IsInPit == 1,
	}

	m.nextBufferPos = (m.nextBufferPos + 1) % lapRecordBufSize

	if m.lapRecordBuf[m.nextBufferPos] == nil {
		m.RefreshBuffer()
	}
	return nil
}

func (m *LapRecordStream) CompleteSector(sector *models.LapSector) {
	sector.SectorTime = m.getLatestLapRecord().PreviousSectorTime
}

func (m *LapRecordStream) CompleteLap(lap *models.Lap) {
	lap.LapTime = m.getLatestLapRecord().PreviousLapTime
}

func (m *LapRecordStream) GetSessionStatus() string {
	return m.getLatestLapRecord().Status
}

func (m *LapRecordStream) GetSessionType() any {
	return m.getLatestLapRecord().SessionType
}

func (m *LapRecordStream) GetBestLapTime() int32 {
	return m.getLatestLapRecord().BestLapTime
}

func (m *LapRecordStream) UpdateLapState(lap *models.Lap) {
	lap.IsValid = !m.getLatestLapRecord().IsValid
}

func (m *LapRecordStream) getLatestLapRecord() *models.LapRecord {
	return m.lapRecordBuf[(m.nextBufferPos-1+lapRecordBufSize)%lapRecordBufSize]
}

func (m *LapRecordStream) getShadowLapRecord() *models.LapRecord {
	return m.lapRecordBuf[(m.nextBufferPos-2+lapRecordBufSize)%lapRecordBufSize]
}
