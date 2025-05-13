package sessionreader

import (
	"time"

	"github.com/kardianos/service"
	"github.com/kat-lego/acc-laptime-tracker/pkg/accstatereader"
	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
)

type SessionReader struct {
	session     *models.Session
	statereader *accstatereader.AccStateReader
	logger      *service.Logger
}

func New(logger *service.Logger) *SessionReader {
	r := SessionReader{}
	r.session = nil
	r.logger = logger
	r.statereader = accstatereader.New()

	return &r
}

func (r *SessionReader) GetSessionUpdates() []*models.Session {
	logger := *r.logger
	updates := make([]*models.Session, 0)

	state, err := r.statereader.GetState()
	if err != nil {
		logger.Error("failed to get the game state")
		r.session = nil
		return updates
	}

	if state.Status == "ACC_PAUSE" || state.Status == "ACC_REPLAY" {
		logger.Infof("game state is paused. status: %s", state.Status)
		return updates
	}

	s, isNewSession := r.getSession(state)
	if s == nil {
		if r.session != nil {
			r.completeSession(state)
			updates = append(updates, r.session)
		}
		r.session = nil
		return updates
	}

	if isNewSession {
		if r.session != nil {
			r.completeSession(state)
			updates = append(updates, r.session)
		}
		updates = append(updates, s)
		r.session = s
		return updates
	}

	l, isNewLap := r.getLap(state)
	if isNewLap {
		r.completeLastLap(state)
		s.Laps = append(s.Laps, l)
		return updates
	}

	ls, isNewLapSector := r.getLapSector(state)
	if isNewLapSector {
		r.completeLastLapSector(state)
		l.LapSectors = append(l.LapSectors, ls)
	}

	updates = append(updates, s)
	return updates
}

func (r *SessionReader) completeSession(state *models.AccGameState) {
	r.session.IsActive = false
	r.completeLastLap(state)
}

func (r *SessionReader) completeLastLap(state *models.AccGameState) {
	r.session.CompletedLaps = state.CompletedLaps

	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]
	l.IsActive = false
	l.LapTime = state.PreviousLapTime

	if nlaps > 1 {
		l.LapDelta = l.LapTime - r.session.Laps[nlaps-2].LapTime
	}

	if r.session.BestLap == 0 || l.LapTime < r.session.Laps[r.session.BestLap-1].LapTime {
		r.session.BestLap = l.LapNumber
	}
	r.completeLastLapSector(state)
}

func (r *SessionReader) completeLastLapSector(state *models.AccGameState) {
	l := r.session.Laps[len(r.session.Laps)-1]
	ls := l.LapSectors[len(l.LapSectors)-1]
	ls.IsActive = false
	ls.SectorTime = state.PreviousSectorTime

	if ls.SectorNumber == r.session.NumberOfSectors {
		ls.SectorTime = l.LapTime
	}
}

func (r *SessionReader) getSession(state *models.AccGameState) (*models.Session, bool) {
	logger := *r.logger

	if state.SessionType == "ACC_UNKNOWN" || state.Status == "ACC_OFF" {
		logger.Infof("there is a break in session continuity [session type: %s][status: %s]",
			state.SessionType, state.Status)
		return nil, true
	}

	if r.session != nil && r.session.CompletedLaps <= state.CompletedLaps {
		r.session.CompletedLaps = state.CompletedLaps
		return r.session, false
	}

	logger.Info("starting new session")

	ts := time.Now().UTC().Unix()
	t := time.Unix(ts, 0).UTC()
	id := t.Format("20060102T150405Z")

	session := models.Session{
		Id:              id,
		StartTime:       ts,
		SessionType:     state.SessionType,
		Track:           state.Track,
		NumberOfSectors: state.SectorCount,
		CarModel:        state.CarModel,
		IsActive:        true,
		CompletedLaps:   0,
		BestLap:         0,
		Player:          "anonymous",

		Laps: []*models.Lap{
			{
				LapNumber: 1,
				IsValid:   true,
				IsActive:  true,
				LapSectors: []*models.LapSector{
					{
						SectorNumber: 1,
						IsActive:     true,
					},
				},
			},
		},
	}

	return &session, true
}

func (r *SessionReader) getLap(state *models.AccGameState) (*models.Lap, bool) {
	logger := *r.logger

	trackedLaps := len(r.session.Laps)
	if int32(trackedLaps-1) == state.CompletedLaps {
		l := r.session.Laps[trackedLaps-1]
		l.IsValid = state.IsValid
		l.LapTime = state.CurrentLapTime
		return l, false
	}

	logger.Info("new lap started")
	lap := models.Lap{
		LapNumber: state.CompletedLaps + 1,
		IsValid:   true,
		IsActive:  true,
		LapSectors: []*models.LapSector{
			{
				SectorNumber: 1,
				IsActive:     true,
			},
		},
	}

	return &lap, true
}

func (r *SessionReader) getLapSector(state *models.AccGameState) (*models.LapSector, bool) {
	logger := *r.logger
	trackedLaps := len(r.session.Laps)
	latestLap := r.session.Laps[trackedLaps-1]
	trackedSectors := len(latestLap.LapSectors)

	if int32(trackedSectors-1) == state.CurrentSectorIndex {
		ls := latestLap.LapSectors[trackedSectors-1]
		return ls, false
	}

	logger.Info("new lap sector started")
	lapSector := models.LapSector{
		SectorNumber: state.CurrentSectorIndex + 1,
		IsActive:     true,
	}

	return &lapSector, true
}
