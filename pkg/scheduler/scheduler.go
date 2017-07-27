package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

const kubeCertsPath = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
const kubeTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
const listCronJobsPath = "/apis/kubeheads.pavlov.ai/v1/cronjobs"
const createJobPath = "/apis/batch/v1/namespaces/{{.Metadata.Namespace}}/jobs"

var (
	createJobTmpl *template.Template
)

func init() {
	var err error
	createJobTmpl, err = template.New("createJobPath").Parse(createJobPath)
	if err != nil {
		panic(err)
	}
}

type cronJob struct {
	Metadata struct {
		Namespace       string `json:"namespace"`
		Name            string `json:"name"`
		UID             string `json:"uid"`
		ResourceVersion string `json:"resourceVersion"`
	} `json:"metadata"`
	Spec struct {
		Schedule    string `json:"schedule"`
		JobTemplate struct {
			Spec interface{} `json:"spec"`
		} `json:"jobTemplate"`
	} `json:"spec"`
}

type jobMetadata struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type job struct {
	APIVersion string      `json:"apiVersion"`
	Kind       string      `json:"kind"`
	Metadata   jobMetadata `json:"metadata"`
	Spec       interface{} `json:"spec"`
}

type cronJobList struct {
	Items []cronJob `json:"items"`
}

type cronJobSchedulerEntry struct {
	EntryID cron.EntryID
	CronJob cronJob
}

type cronJobScheduler struct {
	client         *http.Client
	apiServerURL   string
	apiServerToken string
	cron           *cron.Cron
	cronJobEntries map[string]cronJobSchedulerEntry
}

func logCronJobStatus(cj cronJob, message string) {
	logrus.WithFields(logrus.Fields{
		"namespace":       cj.Metadata.Namespace,
		"name":            cj.Metadata.Name,
		"uid":             cj.Metadata.UID,
		"resourceVersion": cj.Metadata.ResourceVersion,
		"schedule":        cj.Spec.Schedule,
	}).Info(message)
}

func newCronJobScheduler(apiServerURL string, fetchSchedule string) *cronJobScheduler {
	c := cron.New()

	caCertPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(kubeCertsPath)
	if err == nil {
		caCertPool.AppendCertsFromPEM(caCert)
	}

	apiServerToken, _ := ioutil.ReadFile(kubeTokenPath)

	s := cronJobScheduler{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		},
		apiServerURL:   apiServerURL,
		apiServerToken: string(apiServerToken),
		cron:           c,
		cronJobEntries: make(map[string]cronJobSchedulerEntry),
	}

	c.AddFunc(fetchSchedule, func() {
		err := s.SyncCronJobs()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"err": err,
			}).Fatal("error fetching CronJobs")
		}
	})

	return &s
}

func (s *cronJobScheduler) Start() {
	err := s.SyncCronJobs()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Fatal("error fetching CronJobs")
	}

	s.cron.Start()
}

func (s *cronJobScheduler) ScheduleCronJob(cj cronJob) error {
	cjID := cj.Metadata.UID
	schedule := cj.Spec.Schedule

	entryID, err := s.cron.AddFunc(schedule, func() {
		t := time.Now()
		s.TriggerCronJob(cj, t)
	})

	if err != nil {
		return err
	}

	s.cronJobEntries[cjID] = cronJobSchedulerEntry{entryID, cj}
	return nil
}

func (s *cronJobScheduler) Stop() {
	s.cron.Stop()
}

func (s *cronJobScheduler) SyncCronJobs() error {
	cronJobs, err := s.fetchCronJobs()
	if err != nil {
		return err
	}

	cjIDSet := make(map[string]bool)

	for _, cj := range cronJobs {
		cjID := cj.Metadata.UID

		if entry, ok := s.cronJobEntries[cjID]; ok {
			if differentResourceVersions(entry.CronJob, cj) {
				s.cron.Remove(entry.EntryID)

				if err = s.ScheduleCronJob(cj); err != nil {
					return err
				}

				logCronJobStatus(cj, "updated CronJob")
			}
		} else {
			if err = s.ScheduleCronJob(cj); err != nil {
				return err
			}

			logCronJobStatus(cj, "added CronJob")
		}

		cjIDSet[cjID] = true
	}

	for cjID, entry := range s.cronJobEntries {
		if _, ok := cjIDSet[cjID]; !ok {
			delete(s.cronJobEntries, cjID)
			s.cron.Remove(entry.EntryID)
			logCronJobStatus(entry.CronJob, "removed CronJob")
		}
	}

	return nil
}

func (s *cronJobScheduler) TriggerCronJob(cj cronJob, t time.Time) {
	logCronJobStatus(cj, "triggered CronJob")
	newJob := newJobFromCronJob(cj, t)
	if err := s.createJob(newJob); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Error("error creating Job")
	}
	logrus.WithFields(logrus.Fields{
		"namespace": newJob.Metadata.Namespace,
		"name":      newJob.Metadata.Name,
	}).Info("creating Job")
}

func (s *cronJobScheduler) request(method string, path string, body io.Reader) (*http.Response, error) {
	url := s.apiServerURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+s.apiServerToken)
	return s.client.Do(req)
}

func (s *cronJobScheduler) fetchCronJobs() ([]cronJob, error) {
	res, err := s.request("GET", listCronJobsPath, nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cjL cronJobList
	err = json.Unmarshal(body, &cjL)
	if err != nil {
		return nil, err
	}

	return cjL.Items, nil
}

func differentResourceVersions(a cronJob, b cronJob) bool {
	return a.Metadata.UID == b.Metadata.UID && a.Metadata.ResourceVersion != b.Metadata.ResourceVersion
}

func newJobFromCronJob(cj cronJob, t time.Time) job {
	jobName := cj.Metadata.Name + "-" + fmt.Sprintf("%d", t.Unix())

	return job{
		APIVersion: "batch/v1",
		Kind:       "Job",
		Metadata: jobMetadata{
			Namespace: cj.Metadata.Namespace,
			Name:      jobName,
		},
		Spec: cj.Spec.JobTemplate.Spec,
	}
}

func (s *cronJobScheduler) createJob(j job) error {
	buf := bytes.NewBufferString("")
	if err := createJobTmpl.Execute(buf, j); err != nil {
		return err
	}

	b, err := json.Marshal(j)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)

	_, err = s.request("POST", buf.String(), reader)
	if err != nil {
		return err
	}

	return nil
}
