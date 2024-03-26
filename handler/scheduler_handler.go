package handler

import (
	"dq_scheduler_v2/config"
	"dq_scheduler_v2/model"
	"dq_scheduler_v2/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SchedulerHandler struct {
	schedulerService *service.SchedulerService
	config           *config.Config
}

func NewSchedulerHandler(schedulerService *service.SchedulerService, config *config.Config) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerService: schedulerService,
		config:           config,
	}
}

func (h *SchedulerHandler) CreateScheduler(c *gin.Context) {
	var scheduler model.Scheduler
	if err := c.ShouldBindJSON(&scheduler); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.config.SaveScheduler(scheduler)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save scheduler"})
		return
	}

	if scheduler.Enabled {
		h.schedulerService.StartScheduler()
	}

	c.JSON(http.StatusCreated, scheduler)
}

func (h *SchedulerHandler) ListSchedulers(c *gin.Context) {
	schedulers, err := h.config.LoadSchedulers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load schedulers"})
		return
	}
	c.JSON(http.StatusOK, schedulers)
}

func (h *SchedulerHandler) GetScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	schedulers, err := h.config.LoadSchedulers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load schedulers"})
		return
	}

	for _, scheduler := range schedulers {
		if scheduler.ID == schedulerID {
			c.JSON(http.StatusOK, scheduler)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
}

func (h *SchedulerHandler) UpdateScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	var updatedScheduler model.Scheduler
	if err := c.ShouldBindJSON(&updatedScheduler); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedScheduler.ID = schedulerID
	err := h.config.UpdateScheduler(updatedScheduler)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update scheduler"})
		return
	}
	c.JSON(http.StatusOK, updatedScheduler)
}

func (h *SchedulerHandler) DeleteScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	err := h.config.DeleteScheduler(schedulerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete scheduler"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Scheduler deleted successfully"})
}

func (h *SchedulerHandler) StartScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	schedulers, err := h.config.LoadSchedulers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load schedulers"})
		return
	}

	for i, scheduler := range schedulers {
		if scheduler.ID == schedulerID {
			schedulers[i].Enabled = true
			err := h.config.UpdateScheduler(schedulers[i])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start scheduler"})
				return
			}
			h.schedulerService.StartScheduler()
			c.JSON(http.StatusOK, gin.H{"message": "Scheduler started successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
}

func (h *SchedulerHandler) StopScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	schedulers, err := h.config.LoadSchedulers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load schedulers"})
		return
	}

	for i, scheduler := range schedulers {
		if scheduler.ID == schedulerID {
			schedulers[i].Enabled = false
			err := h.config.UpdateScheduler(schedulers[i])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop scheduler"})
				return
			}
			h.schedulerService.StopScheduler()
			c.JSON(http.StatusOK, gin.H{"message": "Scheduler stopped successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
}

func (h *SchedulerHandler) ResumeScheduler(c *gin.Context) {
	schedulerID := c.Param("id")
	schedulers, err := h.config.LoadSchedulers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load schedulers"})
		return
	}

	for i, scheduler := range schedulers {
		if scheduler.ID == schedulerID {
			schedulers[i].Enabled = true
			err := h.config.UpdateScheduler(schedulers[i])
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resume scheduler"})
				return
			}
			h.schedulerService.StartScheduler()
			c.JSON(http.StatusOK, gin.H{"message": "Scheduler resumed successfully"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Scheduler not found"})
}
