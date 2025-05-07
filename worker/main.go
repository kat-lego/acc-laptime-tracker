package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/kardianos/service"
	"github.com/kat-lego/acc-laptime-tracker/pkg/watchers"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	logger      service.Logger
	startACCOpt bool
)

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
	// Start ACC if flag passed and ACC not running
	if startACCOpt && !isACCRunning() {
		err := startACC()
		if err != nil {
			logger.Errorf("Failed to start ACC: %v", err)
		} else {
			logger.Info("ACC started successfully")
			time.Sleep(10 * time.Second) // allow time for ACC to launch
		}
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
				if startACCOpt {
					logger.Info("ACC has stopped. Exiting because --start-acc was passed.")
					p.Stop(nil)
					os.Exit(0)
				} else {
					logger.Info("ACC is not running")
				}
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

func startACC() error {
	accPath := os.Getenv("ACC_PATH")
	if accPath == "" {
		return &os.PathError{Op: "startACC", Path: "ACC_PATH", Err: os.ErrNotExist}
	}
	cmd := exec.Command(accPath)
	return cmd.Start()
}

func parseFlags() {
	for _, arg := range os.Args[1:] {
		if arg == "--start-acc" {
			startACCOpt = true
			break
		}
	}
}

func main() {
	parseFlags()

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
