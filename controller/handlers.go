package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseContainer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Running bool   `json:"running"`
	Created int64  `json:"created"`
}

type ResponseImage struct {
	ID      string `json:"id"`
	Size    int64  `json:"size"`
	Created int64  `json:"created"`
	Tags    string `json:"tags"`
}

type ResponseDockerStats struct {
	ID          string    `json:"id"`
	MemoryUsage uint64    `json:"memory_usage"`
	MemoryLimit uint64    `json:"memory_limit"`
	CPUUsage    float64   `json:"cpu_percentage"`
	Network     [2]uint64 `json:"network_io"`
}

func GetImageHandler(host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		images := DockersPool.FetchImages(host)
		resp := make([]ResponseImage, 0)
		for _, v := range images {
			resp = append(resp, ResponseImage{
				ID:      v.ID,
				Created: v.Created,
				Size:    v.Size,
				Tags:    ImageName(v),
			})
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func GetContainerHandler(host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status := ctx.DefaultQuery("status", "running")
		containers := DockersPool.FetchContainers(host)
		resp := make([]ResponseContainer, 0)
		for _, v := range containers {
			if v.State == status || status == "all" {
				resp = append(resp, ResponseContainer{
					ID:      v.ID,
					Name:    v.Names[0],
					Running: v.State == "running",
					Created: v.Created,
				})
			}
		}
		ctx.JSON(http.StatusOK, resp)
	}
}

func GetLogsHandler(host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("container")
		if len(id) == 0 {
			ctx.Status(http.StatusNotFound)
			return
		}
		DockersPool.GetLogs(host, id, ctx.Writer)
	}
}

func GetStatsHandler(host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		stats := make([]ResponseDockerStats, 0)
		containers := DockersPool.FetchContainers(host)
		for _, container := range containers {
			stat := DockersPool.GetStats(host, container.ID)
			stats = append(stats, ResponseDockerStats{
				ID:          container.ID,
				MemoryUsage: ContainerMemoryUsage(stat.MemoryStats),
				MemoryLimit: stat.MemoryStats.Limit,
				CPUUsage:    ContainerCPUPercantage(stat.PreCPUStats, stat.CPUStats),
				Network:     ContainerNetwork(stat.Networks),
			})
		}
		ctx.JSON(http.StatusOK, stats)
	}
}

func GetCommandHandler(host string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Query("container")
		cmd := ctx.Query("command")
		if len(id) == 0 {
			ctx.Status(http.StatusBadRequest)
			return
		}
		if cmd == "stop" {
			DockersPool.StopContainer(host, id)
		} else if cmd == "start" {
			DockersPool.StartContainer(host, id)
		}
	}
}
