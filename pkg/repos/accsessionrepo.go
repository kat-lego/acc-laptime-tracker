package repos

import (
	"log"

	"github.com/kat-lego/acc-laptime-tracker/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresAccSessionRepo struct {
	db *gorm.DB
}

func NewPostgresAccSessionsRepo(dsn string) *PostgresAccSessionRepo {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Automigrate your models to keep the schema up to date
	err = db.AutoMigrate(&models.Session{}, &models.Lap{}, &models.LapSector{})
	if err != nil {
		log.Fatalf("failed to automigrate: %v", err)
	}

	return &PostgresAccSessionRepo{db: db}
}

func (p *PostgresAccSessionRepo) CreateSession(session *models.Session) (uint, error) {
	result := p.db.Create(session)
	return session.Id, result.Error
}

func (p *PostgresAccSessionRepo) CreateLap(lap *models.Lap) (uint, error) {
	result := p.db.Create(lap)
	return lap.Id, result.Error
}

func (p *PostgresAccSessionRepo) CreateLapSector(lapSector *models.LapSector) (uint, error) {
	result := p.db.Create(lapSector)
	return lapSector.Id, result.Error
}

func (p *PostgresAccSessionRepo) UpdateSessionStats(sessionId uint, completedLaps int32, bestLapTime int32, previousLapTime int32) error {
	result := p.db.Model(&models.Session{}).Where("id = ?", sessionId).Updates(models.Session{
		CompletedLaps:   completedLaps,
		BestLapTime:     bestLapTime,
		PreviousLapTime: previousLapTime,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (p *PostgresAccSessionRepo) GetSessionsWithCount(limit int, offset int) ([]models.Session, int64, error) {
	var sessions []models.Session
	var count int64

	if err := p.db.Model(&models.Session{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	query := p.db.
		Preload("Laps").
		Preload("Laps.LapSectors").
		Order("start_time desc")

	if limit != -1 {
		query = query.Limit(limit).Offset(offset)
	}

	result := query.Find(&sessions)
	return sessions, count, result.Error
}
