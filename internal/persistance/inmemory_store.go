package persistance

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"web-page-analyzer/model"
)

/*
maintain jobs and results in memory
*/
type InMemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]*model.Job
}

func NewMemoryStore() *InMemoryStore {
	logrus.Info("Initializing InMemoryStore")
	return &InMemoryStore{
		jobs: make(map[string]*model.Job),
	}
}

func (s *InMemoryStore) CreateJob(url string) model.Job {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.New().String()
	job := model.Job{
		ID:        id,
		Status:    "PROCESSING",
		CraetedAt: time.Now(),
		Result:    nil,
	}
	s.jobs[id] = &job

	logrus.WithFields(logrus.Fields{
		"job_id": job.ID,
		"url":    url,
	}).Info("Created new job")

	return job
}

func (s *InMemoryStore) CompleteJob(id string, result model.AnalysisResult, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[id]
	if !ok {
		logrus.WithField("job_id", id).Warn("Job not found while completing")
		return
	}

	if err != nil {
		job.Error = err.Error()
		job.Status = "FAILED"
		job.Result = nil

		logrus.WithFields(logrus.Fields{
			"job_id": job.ID,
			"error":  err.Error(),
		}).Error("Job marked as FAILED")
	} else {
		job.Status = "COMPLETED"
		job.Result = &result

		logrus.WithFields(logrus.Fields{
			"job_id": job.ID,
		}).Info("Job marked as COMPLETED")
	}
}

func (s *InMemoryStore) GetJob(id string) (*model.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	if ok {
		logrus.WithField("job_id", id).Info("Job retrieved from store")
	} else {
		logrus.WithField("job_id", id).Warn("Job not found in store")
	}

	return job, ok
}

func (s *InMemoryStore) ListJobs() []*model.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		list = append(list, job)
	}

	logrus.WithField("count", len(list)).Info("Listing all jobs")

	return list
}
