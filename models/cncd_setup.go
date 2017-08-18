package models

import (
	cc "context"
	"strings"

	"github.com/cncd/logging"
	"github.com/cncd/pubsub"
	"github.com/cncd/queue"
	"github.com/urfave/cli"
)

func LoadCNCDGlobalConfig(c *cli.Context) {

	// limits
	CNCDConfig.Pipeline.Limits.MemSwapLimit = c.Int64("limit-mem-swap")
	CNCDConfig.Pipeline.Limits.MemLimit = c.Int64("limit-mem")
	CNCDConfig.Pipeline.Limits.ShmSize = c.Int64("limit-shm-size")
	CNCDConfig.Pipeline.Limits.CPUQuota = c.Int64("limit-cpu-quota")
	CNCDConfig.Pipeline.Limits.CPUShares = c.Int64("limit-cpu-shares")
	CNCDConfig.Pipeline.Limits.CPUSet = c.String("limit-cpu-set")

	// server configuration
	CNCDConfig.Server.Cert = c.String("server-cert")
	CNCDConfig.Server.Key = c.String("server-key")
	CNCDConfig.Server.Pass = c.String("agent-secret")
	CNCDConfig.Server.Host = strings.TrimRight(c.String("server-host"), "/")
	CNCDConfig.Server.Port = c.String("server-addr")
	CNCDConfig.Pipeline.Networks = c.StringSlice("network")
	CNCDConfig.Pipeline.Volumes = c.StringSlice("volume")
	CNCDConfig.Pipeline.Privileged = c.StringSlice("escalate")
	// CNCDConfig.Server.Open = cli.Bool("open")
	// CNCDConfig.Server.Orgs = sliceToMap(cli.StringSlice("orgs"))
	// CNCDConfig.Server.Admins = sliceToMap(cli.StringSlice("admin"))
}

func InitCNCDGlobalService() {

	// services
	CNCDConfig.Services.Queue = setupQueue()
	CNCDConfig.Services.Logs = logging.New()
	CNCDConfig.Services.Pubsub = pubsub.New()
	CNCDConfig.Services.Pubsub.Create(cc.Background(), "topic/events")

}

func setupQueue() queue.Queue {
	return WithTaskStore(queue.New())
}
