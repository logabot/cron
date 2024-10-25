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
	constEntrypoint = ""
)

var (
	defaultShell      string
	defaultConfigPath string
	defaultEntrypoint string

	shell      string
	configPath string
	entrypoint string
	crons      []string
	s          gocron.Scheduler
	job        gocron.Job
)

func main() {

	envConfigPath, ok := os.LookupEnv("CONFIG")
	if ok {
		defaultConfigPath = envConfigPath
	} else {
		defaultConfigPath = constConfigPath
	}

	envEntrypoint, ok := os.LookupEnv("ENTRYPOINT")
	if ok {
		defaultEntrypoint = envEntrypoint
	}

	envShell, ok := os.LookupEnv("SHELL")
	if ok {
		defaultShell = envShell
	} else {
		defaultShell = constShell
	}

	flag.StringVar(&shell, "shell", defaultShell, "define shell for run cron command")
	flag.StringVar(&configPath, "config", defaultConfigPath, "crontab file")
	flag.StringVar(&entrypoint, "entrypoint", defaultEntrypoint, "shell script filename that runs before cron start")
	flag.Parse()

	MainLog := log.New(os.Stdout, "Cron: ", log.LstdFlags|log.Lmsgprefix)

	pwd, _ := os.Executable()
	if !strings.HasPrefix(entrypoint, "/") {
		entrypoint = filepath.Join(filepath.Dir(pwd), entrypoint)
	}
	_, err := os.Stat(entrypoint)
	if len(entrypoint) > 0 {
		if err != nil {
			MainLog.Printf("Can't open entrypoint file: %s\n", err)
			os.Exit(1)
		}
		MainLog.Printf("Run entrypoint: %s\n", entrypoint)
		startup := exec.Command(shell, entrypoint)

		if err != nil {
			MainLog.Printf("Can't run entrypoint: %s\n", err)

		}
		startup.Dir = filepath.Dir(pwd)
		output, err := startup.CombinedOutput()
		if err != nil {
			MainLog.Printf("Can't run entrypoint. Reason: %s\n", err)
		}
		MainLog.Println(string(output))

	}
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
	cancelChan := make(chan os.Signal, 1)

	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
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
