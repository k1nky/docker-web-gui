package controller

import (
	"strings"

	"github.com/docker/docker/api/types"
)

func ContainerMemoryUsage(mem types.MemoryStats) uint64 {
	// see https://github.com/docker/docker-ce/blob/ad2621c56a511ee668d1431a93266583d1567a9b/components/cli/cli/command/container/stats_helpers.go#L251
	if v, isCgroup1 := mem.Stats["total_inactive_file"]; isCgroup1 && v < mem.Usage {
		return mem.Usage - v
	}
	if v := mem.Stats["inactive_file"]; v < mem.Usage {
		return mem.Usage - v
	}
	return mem.Usage
}

func ContainerCPUPercantage(prev types.CPUStats, cur types.CPUStats) float64 {
	// https://github.com/docker/docker-ce/blob/ad2621c56a511ee668d1431a93266583d1567a9b/components/cli/cli/command/container/stats_helpers.go#L166
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(cur.CPUUsage.TotalUsage) - float64(prev.CPUUsage.TotalUsage)
		// calculate the change for the entire system between readings
		systemDelta = float64(cur.SystemUsage) - float64(prev.SystemUsage)
		onlineCPUs  = float64(cur.OnlineCPUs)
	)

	if onlineCPUs == 0.0 {
		onlineCPUs = float64(len(cur.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}

func ContainerNetwork(net map[string]types.NetworkStats) [2]uint64 {
	s := [2]uint64{0, 0}
	for _, v := range net {
		s[0] += v.RxBytes
		s[1] += v.TxBytes
	}
	return s
}

func ImageName(image types.ImageSummary) string {
	name := strings.Split(image.RepoDigests[0], "@")[0]
	if image.RepoTags == nil {
		return name + ":<none>"
	}
	return name + ":" + strings.Split(image.RepoTags[0], ":")[1]
}
