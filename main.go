package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
)

// Config is the runner app config structure.
type Config struct {
	Name, DisplayName, Description string
	Dir, Exec                      string
	Args                           []string
	Env                            []string
	Stderr, Stdout                 string
	Interval                       int // Interval in minutes
}

var logger service.Logger

type program struct {
	exit    chan struct{}
	service service.Service
	*Config
}

func (p *program) Start(s service.Service) error {
	go p.runBatchScriptAtInterval()
	return nil
}

func (p *program) runBatchScriptAtInterval() {
	// Sofortiges Ausführen des Batch-Skripts beim Start
	p.runBatchScript()

	interval := time.Duration(p.Interval) * time.Minute
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.runBatchScript()
		case <-p.exit:
			return
		}
	}
}

func (p *program) runBatchScript() {
	cmd := exec.Command(p.Exec, p.Args...)
	cmd.Dir = p.Dir
	cmd.Env = append(os.Environ(), p.Env...)

	if p.Stderr != "" {
		f, err := os.OpenFile(p.Stderr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logger.Warningf("Failed to open stderr %q: %v", p.Stderr, err)
			return
		}
		defer f.Close()
		cmd.Stderr = f
	}
	if p.Stdout != "" {
		f, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			logger.Warningf("Failed to open stdout %q: %v", p.Stdout, err)
			return
		}
		defer f.Close()
		cmd.Stdout = f
	}

	logger.Infof("Running command: %s %v", p.Exec, p.Args)
	err := cmd.Run()
	if err != nil {
		logger.Errorf("Error running command: %v", err)
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	logger.Info("Stopping ", p.DisplayName)
	return nil
}

func getConfigPath() (string, error) {
	fullexecpath, err := os.Executable()
	if err != nil {
		return "", err
	}

	dir, execname := filepath.Split(fullexecpath)
	ext := filepath.Ext(execname)
	name := execname[:len(execname)-len(ext)]

	return filepath.Join(dir, name+".json"), nil
}

func getConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Config{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	configPath, err := getConfigPath()
	if err != nil {
		log.Fatal(err)
	}
	config, err := getConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	svcConfig := &service.Config{
		Name:        config.Name,
		DisplayName: config.DisplayName,
		Description: config.Description,
	}

	prg := &program{
		exit:   make(chan struct{}),
		Config: config,
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	prg.service = s

	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
