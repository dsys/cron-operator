package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	apiServerURL  = flag.String("api-server", "https://kubernetes.default.svc.cluster.local", "URL of the Kubernetes API Server")
	fetchSchedule = flag.String("fetch-schedule", "@every 10s", "the frequency to request the list of CronJob resources")
)

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func main() {
	flag.Parse()
	s := newCronJobScheduler(*apiServerURL, *fetchSchedule)

	s.Start()
	waitForSignal()
	s.Stop()

	logrus.Info("exiting")
}
