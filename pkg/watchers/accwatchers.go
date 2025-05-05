package watchers

import (
	"time"

	"github.com/kardianos/service"
	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
	"github.com/kat-lego/acc-laptime-tracker/pkg/repos"
	"github.com/kat-lego/acc-laptime-tracker/pkg/streams"
)

type AccWatcher struct {
	session             *models.Session
	sessionInvalidNoted bool
	lap                 *models.Lap
	lapSectors          []*models.LapSector
	lapSector           *models.LapSector
	repo                *repos.PostgresAccSessionRepo
	lapRecordSteam      *streams.LapRecordStream
	logger              *service.Logger
}

func NewAccWatcher(posgresConn string, logger *service.Logger) *AccWatcher {
	w := AccWatcher{}
	w.session = nil
	w.sessionInvalidNoted = true
	w.lap = nil
	w.lapSectors = nil
	w.lapSector = nil
	w.repo = repos.NewPostgresAccSessionsRepo(posgresConn)
	w.lapRecordSteam = streams.NewLapRecordStream()
	w.logger = logger

	return &w
}

func (w *AccWatcher) Peep() {
	logger := *w.logger

	if err := w.lapRecordSteam.RefreshBuffer(); err != nil {
		logger.Error("Error with refreshing buffer")
		w.lap = nil
		w.session = nil
		w.lapSector = nil
		for i, _ := range w.lapSectors {
			w.lapSectors[i] = nil
		}
		return
	}

	s := w.lapRecordSteam.GetNextSession()
	if s != nil {
		logger.Infof("Session Started | [%s][%s][%s][%s]",
			time.Unix(s.StartTime, 0).UTC().Format(time.RFC3339), s.Track, s.CarModel, s.SessionType)

		w.session = s
		w.repo.CreateSession(w.session)
	}

	var sessionInvalid = w.session == nil || w.session.Id == 0 || w.lapRecordSteam.GetSessionStatus() != "ACC_LIVE" || w.lapRecordSteam.GetSessionStatus() == "ACC_UNKNOWN"
	if sessionInvalid {
		if !w.sessionInvalidNoted {
			logger.Infof("session tracking paused | session status %s | session type %s",
				w.lapRecordSteam.GetSessionStatus(), w.lapRecordSteam.GetSessionType())
			w.sessionInvalidNoted = true
		}

		return
	}

	if w.sessionInvalidNoted {
		logger.Infof("session tracking resumed")
		w.sessionInvalidNoted = false
	}

	ls := w.lapRecordSteam.GetNextLapSector()
	if ls != nil {
		logger.Infof("Sector %d Started", ls.SectorNumber)

		if w.lapSector != nil {
			w.lapRecordSteam.CompleteSector(w.lapSector)
			w.lapSectors = append(w.lapSectors, w.lapSector)
			logger.Infof("Sector %d Completed | time %s", w.lapSector.SectorNumber, time.Duration(w.lapSector.SectorTime)*time.Millisecond)
		}
		w.lapSector = ls
	}

	if w.lap != nil {
		w.lapRecordSteam.UpdateLapState(w.lap)
	}

	l := w.lapRecordSteam.GetNextLap()
	if l != nil {
		logger.Infof("Lap %d Started", l.LapNumber)
		if w.lap != nil {
			w.lap.SessionId = w.session.Id
			w.lapRecordSteam.CompleteLap(w.lap)

			w.lap.LapSectors = make([]models.LapSector, len(w.lapSectors))
			for i, sector := range w.lapSectors {
				w.lap.LapSectors[i] = *sector
				w.lapSectors[i] = nil
			}
			w.lap.LapSectors[len(w.lap.LapSectors)-1].SectorTime = w.lap.LapTime

			w.repo.CreateLap(w.lap)
			w.repo.UpdateSessionStats(w.lap.SessionId, w.lap.LapNumber, w.lapRecordSteam.GetBestLapTime(), w.lap.LapTime)
			logger.Infof("Lap %d Completed | time %s", w.lap.LapNumber, time.Duration(w.lap.LapTime)*time.Millisecond)
			for _, sector := range w.lap.LapSectors {
				logger.Infof("Sector %d | time %s", sector.SectorNumber, time.Duration(sector.SectorTime)*time.Millisecond)
			}
		}
		w.lap = l
	}

}
