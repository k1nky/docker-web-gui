package controller

import (
	"context"
	"encoding/json"
	"io"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ClientsPool map[string]*client.Client

var (
	DockersPool ClientsPool
)

func init() {
	DockersPool = NewClientsPool()
}

func NewClientsPool() ClientsPool {
	return make(map[string]*client.Client, 0)
}

func (cp ClientsPool) GetClient(name string) *client.Client {
	_, exists := cp[name]
	if !exists {
		cli, err := client.NewClientWithOpts(client.WithHost(name),
			client.WithAPIVersionNegotiation())
		if err != nil {
			log.Println(err)
			return nil
		}
		cp[name] = cli
	}
	return cp[name]
}

func (cp ClientsPool) FetchImages(host string) []types.ImageSummary {
	ctx := context.Background()
	if cli := cp.GetClient(host); cli != nil {
		images, err := cli.ImageList(ctx, types.ImageListOptions{})
		if err == nil {
			return images
		}
		log.Println(err)
	}

	return []types.ImageSummary{}
}

func (cp ClientsPool) FetchContainers(host string) []types.Container {
	ctx := context.Background()
	if cli := cp.GetClient(host); cli != nil {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
			All: true,
		})
		if err == nil {
			return containers
		}

		log.Println(err)
	}

	return []types.Container{}
}

func (cp ClientsPool) GetLogs(host string, containerId string, w io.Writer) {
	ctx := context.Background()
	if cli := cp.GetClient(host); cli != nil {
		logs, err := cli.ContainerLogs(ctx, containerId, types.ContainerLogsOptions{
			ShowStdout: true,
			ShowStderr: true,
		})
		if err != nil {
			log.Println(err)
			return
		}
		io.Copy(w, logs)
	}
}

func (cp ClientsPool) GetStats(host string, containerId string) types.StatsJSON {
	var cs types.StatsJSON

	ctx := context.Background()
	if cli := cp.GetClient(host); cli != nil {
		stats, err := cli.ContainerStatsOneShot(ctx, containerId)
		if err != nil {
			log.Println(err)
		}
		if err := json.NewDecoder(stats.Body).Decode(&cs); err != nil {
			log.Println(err)
		}
		stats.Body.Close()
	}
	return cs
}

func (cp ClientsPool) StopContainer(host string, containerId string) {
	ctx := context.Background()
	if cli := cp.GetClient(host); cli != nil {
		if err := cli.ContainerStop(ctx, containerId, nil); err != nil {
			log.Println(err)
		}
	}
}

func (cp ClientsPool) StartContainer(host string, containerId string) {
	ctx := context.Background()
	if cli := cp.GetClient(host); cli != nil {
		if err := cli.ContainerStart(ctx, containerId, types.ContainerStartOptions{}); err != nil {
			log.Println(err)
		}
	}
}
