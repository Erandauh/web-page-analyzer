package persistance

import (
	"sync"
	"time"

	"github.com/google/uuid"

	"web-page-analyzer/model"
)

type MemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]*model.Job
}

var DefaultStore = NewMemoryStore()

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]*model.Job),
	}
}

func (s *MemoryStore) CreateJob(url string) model.Job {
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
	return job
}

func (s *MemoryStore) CompleteJob(id string, result model.AnalysisResult, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if job, ok := s.jobs[id]; ok {
		if err != nil {
			job.Error = err
			job.Status = "FAILED"
			job.Result = nil
		} else {
			job.Status = "COMPLETED"
			job.Result = &result
		}
	}
}

func (s *MemoryStore) GetJob(id string) (*model.Job, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	job, ok := s.jobs[id]
	return job, ok
}

func (s *MemoryStore) ListJobs() []*model.Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		list = append(list, job)
	}
	return list
}
