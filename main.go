package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-co-op/gocron/v2"
)

const (
	constConfigPath = "config"
	constShell      = "/bin/sh"
)

var (
	shell             string
	configPath        string
	defaultShell      string
	defaultConfigPath string
	crons             []string
	s                 gocron.Scheduler
	job         gocron.Job
)

func main() {

	cancelChan := make(chan os.Signal, 1)

	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	envConfigPath, ok := os.LookupEnv("CONFIG")
	if ok {
		defaultConfigPath = envConfigPath
	} else {
		defaultConfigPath = constConfigPath
	}
	envShell, ok := os.LookupEnv("RUN_SHELL")
	if ok {
		defaultShell = envShell
	} else {
		defaultShell = constShell
	}

	flag.StringVar(&shell, "shell", defaultShell, "define shell for run cron command")
	flag.StringVar(&configPath, "config", defaultConfigPath, "crontab file")
	flag.Parse()

	MainLog := log.New(os.Stdout, "Cron: ", log.LstdFlags|log.Lmsgprefix)

	file, err := os.ReadFile(configPath)
	if err != nil {
		MainLog.Printf("Can't read file %s\n", err)
	}
	crons = strings.Split(string(file), "\n")

	if s == nil {
		s, err = gocron.NewScheduler()

		if err != nil {
			MainLog.Printf("Can't create scheduler %s\n", err)
		}
	}
	JobLog := log.New(os.Stdout, "Cron: ", log.LstdFlags|log.Lmsgprefix)
	for _, line := range crons {
		if len(line) < 1 {
			continue
		}
		crontab, name, command := ParseLine(line)

		job, err = s.NewJob(
			gocron.CronJob(crontab, false),
			gocron.NewTask(
				func(shell, name, command string) {
					JobLog.SetPrefix(name + ": ")

					pwd, err := os.Executable()
					if err != nil {
						JobLog.Printf("Can't get pwd: %s\n", err)
					}
					cmd := exec.Command(shell, "-c", command)
					cmd.Dir = filepath.Dir(pwd)

					output, err := cmd.CombinedOutput()
					if err != nil {
						JobLog.Printf("Can't run command: %s. Reason: %s", command, err)
					}
					JobLog.Println(string(output))

				},

				shell,
				name,
				command,
			),
		)
		if err != nil {
			MainLog.Printf("Can't create job: %s\n", err)
		}

	}
	// start the scheduler
	MainLog.Println("Run crontabs")

	s.Start()
	nextRun, err := job.NextRun()
	if err != nil {
		MainLog.Printf("Can't calculate nextRun. %s", err)
	}
	MainLog.Printf("Next run: %s", &nextRun)

	x := <-cancelChan
	MainLog.Printf("Stop crontab. Receive signal: %v", x)
	err = s.Shutdown()
	if err != nil {
		MainLog.Printf("Can't shutdown crontab. Reason: %s", err)
	}

}

func ParseLine(line string) (string, string, string) {
	// 1 * * * *

	slice := strings.Split(line, " ")
	crontab := strings.Join(slice[0:5], " ")
	name := slice[5]
	command := strings.Join(slice[6:], " ")
	return crontab, name, command
}
