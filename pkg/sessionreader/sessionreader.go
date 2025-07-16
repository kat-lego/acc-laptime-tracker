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
		logger.Error("Failed to get the game state.")
		r.stopTrackingSession()
		return updates
	}

	if state.IsSessionPaused() {
		logger.Infof("Session is paused. Status: %s.", state.Status)
		return updates
	}

	if state.IsSessionStopped() {
		if r.hasTrackedSession() {
			updates = append(updates, r.completeSession(state))
		}
		r.stopTrackingSession()
		logger.Infof("Session is stopped. Status: %s.", state.Status)
		return updates
	}

	if !r.hasTrackedSession() {
		updates = append(updates, r.initializeSession(state))
		logger.Infof("No previously tracked session. New session started.")
		return updates
	}

	if r.sessionChanged(state) {
		updates = append(updates, r.completeSession(state))
		updates = append(updates, r.initializeSession(state))
		logger.Infof("New session started")
		return updates
	}

	r.updateSession(state)

	if r.lapChanged(state) {
		r.completeLap(state)
		updates = append(updates, r.initializeLap(state))
		logger.Infof("New lap started")
		return updates
	}

	if r.lapRestarted(state) {
		r.stopTrackingLap()
		updates = append(updates, r.initializeLap(state))
		logger.Infof("Lap restarted")
		return updates
	}

	r.updateLap(state)

	if r.sectorChanged(state) {
		r.completeSector(state)
		updates = append(updates, r.initializeSector(state))
		logger.Infof("New sector started")
		return updates
	}

	return updates
}

// session
func (r *SessionReader) hasTrackedSession() bool {
	return r.session != nil
}

func (r *SessionReader) sessionChanged(state *models.AccGameState) bool {
	return r.session.LapsCompleted > state.CompletedLaps
}

func (r *SessionReader) initializeSession(state *models.AccGameState) *models.Session {
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
		LapsCompleted:   0,
		BestLap:         1,
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

	r.session = &session

	return &session
}

func (r *SessionReader) updateSession(state *models.AccGameState) *models.Session {
	r.session.LapsCompleted = state.CompletedLaps

	return r.session
}

func (r *SessionReader) completeSession(state *models.AccGameState) *models.Session {
	session := r.session
	session.IsActive = false

	nlaps := len(session.Laps)
	l := session.Laps[nlaps-1]
	l.IsActive = false

	nSectors := len(l.LapSectors)
	sc := l.LapSectors[nSectors-1]
	sc.IsActive = false

	// assume the last lap will be incomplete
	session.Laps = session.Laps[:nlaps-1]

	logger := *r.logger
	logger.Infof("Completing session with %s completed laps and %s laps",
		session.LapsCompleted, len(session.Laps))

	return session
}

func (r *SessionReader) stopTrackingSession() {
	r.session = nil
	r.session.IsActive = false
}

// lap
func (r *SessionReader) lapChanged(state *models.AccGameState) bool {
	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]

	return l.LapNumber < state.CompletedLaps+1
}

func (r *SessionReader) lapRestarted(state *models.AccGameState) bool {
	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]

	return l.LapNumber == state.CompletedLaps+1 && l.LapTime > state.CurrentLapTime
}

func (r *SessionReader) initializeLap(state *models.AccGameState) *models.Session {
	lap := models.Lap{
		LapTime:   0,
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

	r.session.Laps = append(r.session.Laps, &lap)

	return r.session
}

func (r *SessionReader) updateLap(state *models.AccGameState) *models.Session {
	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]

	l.LapTime = state.CurrentLapTime
	l.IsValid = state.IsValid

	return r.session
}

func (r *SessionReader) completeLap(state *models.AccGameState) *models.Session {
	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]

	l.IsActive = false
	l.LapTime = state.PreviousLapTime

	if nlaps > 1 {
		l.LapDelta = l.LapTime - r.session.Laps[nlaps-2].LapTime
	}

	if l.LapNumber > 1 && l.LapTime < r.session.Laps[r.session.BestLap-1].LapTime {
		r.session.BestLap = l.LapNumber
	}

	r.completeSector(state)

	return r.session
}

func (r *SessionReader) stopTrackingLap() {
	nlaps := len(r.session.Laps)
	r.session.Laps = r.session.Laps[:nlaps-1]
}

// sector
func (r *SessionReader) sectorChanged(state *models.AccGameState) bool {
	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]

	nSectors := len(l.LapSectors)
	s := l.LapSectors[nSectors-1]

	return s.SectorNumber-1 != state.CurrentSectorIndex
}

func (r *SessionReader) initializeSector(state *models.AccGameState) *models.Session {
	nlaps := len(r.session.Laps)
	l := r.session.Laps[nlaps-1]

	lapSector := models.LapSector{
		SectorNumber: state.CurrentSectorIndex + 1,
		IsActive:     true,
	}

	l.LapSectors = append(l.LapSectors, &lapSector)
	return r.session
}

func (r *SessionReader) completeSector(state *models.AccGameState) {
	l := r.session.Laps[len(r.session.Laps)-1]
	ls := l.LapSectors[len(l.LapSectors)-1]
	ls.IsActive = false
	ls.SectorTime = state.PreviousSectorTime

	if ls.SectorNumber == r.session.NumberOfSectors {
		ls.SectorTime = l.LapTime
	}
}
