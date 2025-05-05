package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kardianos/service"
	"github.com/kat-lego/acc-laptime-tracker/pkg/watchers"
	"github.com/shirou/gopsutil/v3/process"
)

var logger service.Logger

type program struct {
	quit chan struct{}
}

func (p *program) Start(s service.Service) error {
	logger.Info("Service starting...")
	p.quit = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	watcher := watchers.NewAccWatcher(os.Getenv("ACCLTRCR_POSTGRES_CONNECTION_STRING"), &logger)

	for {
		select {
		case <-p.quit:
			logger.Info("Service stopping...")
			return
		default:
			if isACCRunning() {
				watcher.Peep()
			} else {
				logger.Info("ACC is not running")
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.quit)
	return nil
}

func isACCRunning() bool {
	procs, err := process.Processes()
	if err != nil {
		return false
	}
	for _, p := range procs {
		name, err := p.Name()
		if err == nil && name == "acc.exe" {
			return true
		}
	}
	return false
}

func main() {
	svcConfig := &service.Config{
		Name:        "ACCLaptimeTracker",
		DisplayName: "ACC Laptime Service",
		Description: "Tracks ACC laptimes",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
