package service

import (
	"dq_scheduler_v2/config"
	"dq_scheduler_v2/model"
	"log"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	cron    *cron.Cron
	config  *config.Config
	jobFunc func(string)
}

func NewSchedulerService(config *config.Config, jobFunc func(string)) *SchedulerService {
	return &SchedulerService{
		cron:    cron.New(),
		config:  config,
		jobFunc: jobFunc,
	}
}

func (s *SchedulerService) StartScheduler() {
	schedulers, err := s.config.LoadSchedulers()
	if err != nil {
		log.Printf("Failed to load schedulers: %v", err)
		return
	}

	for _, scheduler := range schedulers {
		if scheduler.Enabled {
			s.cron.AddFunc(scheduler.Schedule, func() {
				s.jobFunc(scheduler.Endpoint)
			})
		}
	}
	s.cron.Start()
}

func (s *SchedulerService) StopScheduler() {
	s.cron.Stop()
}

func (s *SchedulerService) GetScheduler(schedulerID string) (*model.Scheduler, error) {
	schedulers, err := s.config.LoadSchedulers()
	if err != nil {
		return nil, err
	}

	for _, scheduler := range schedulers {
		if scheduler.ID == schedulerID {
			return &scheduler, nil
		}
	}
	return nil, nil
}

func (s *SchedulerService) SaveScheduler(scheduler *model.Scheduler) error {
	return s.config.SaveScheduler(*scheduler)
}

func (s *SchedulerService) UpdateScheduler(scheduler *model.Scheduler) error {
	return s.config.UpdateScheduler(*scheduler)
}

func (s *SchedulerService) DeleteScheduler(schedulerID string) error {
	return s.config.DeleteScheduler(schedulerID)
}
