package secbench

import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"

	"github.com/docker/docker/api/types/mount"
	"fmt"
	"github.com/ahmetalpbalkan/dlog"
	"time"
	"bufio"
	"github.com/lunixbochs/vtclean"
	"github.com/zpatrick/go-config"
)

const (
	imageName = "docker/docker-bench-security"
)

var (
	ctx = context.Background()
)

type SecBench struct {
	cli 	*client.Client
	parser 	Parser
	Results ResultSets
	cfg 	map[string]string
	skip 	bool
}

func NewSecBenc(p Parser) (sb *SecBench, err error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return
	}
	sb = &SecBench{
		cli: c,
		parser: p,
		cfg: map[string]string{},
	}
	return
}

func (sb *SecBench) Run(cfg *config.Config) {
	sb.cfg, _ = cfg.Settings()
	sb.RunBench()
}

func (sb *SecBench) RunBench() {
	info, err := sb.cli.Info(ctx)
	if err != nil {
		log.Printf("Failed to fetch engine info: %s", err.Error())
		return
	}
	cntCfg := &container.Config{
		AttachStdin: false,
		AttachStdout: false,
		AttachStderr: false,
		Image: imageName,
	}
	/*
	--cap-add audit_control
	-v /srv/docker:/var/lib/docker
	-v /var/run/docker.sock:/var/run/docker.sock
	-v /etc:/etc
	--label docker_bench_security
	*/
	mSock := mount.Mount{
		Type: mount.TypeBind,
		Source: "/var/run/docker.sock",
		Target: "/var/run/docker.sock",
	}
	mEtc := mount.Mount{
		Type: mount.TypeBind,
		Source: "/etc",
		Target: "/etc",
	}
	mLib := mount.Mount{
		Type: mount.TypeBind,
		Source: info.DockerRootDir,
		Target: "/var/lib/docker",
	}
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode("host"),
		PidMode: container.PidMode("host"),
		CapAdd: strslice.StrSlice{"audit_control"},
		Mounts: []mount.Mount{mSock,mEtc, mLib},
	}
	networkingConfig := &network.NetworkingConfig{}
	container, err := sb.cli.ContainerCreate(ctx, cntCfg, hostConfig,
		networkingConfig, fmt.Sprintf("secbench-%d", time.Now().UnixNano()))
	if err != nil {
		log.Printf("Failed to create container: %s", err.Error())
		return
	}
	err = sb.cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		log.Printf("Failed to start container: %s", err.Error())
		return
	}
	// Parse log
	sb.parseLog(container.ID)
	sb.removeContainer(container.ID)
}

func (sb *SecBench) removeContainer(cid string) {
	err := sb.cli.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{})
	if err != nil {
		log.Panic(err.Error())
	}
}
func (sb *SecBench) pullImnage() error {
	cl, err := sb.cli.ImagePull(ctx, "docker/docker-bench-security", types.ImagePullOptions{})
	defer cl.Close()
	return err
}

func (sb *SecBench) parseLine(line string) {
	nline := vtclean.Clean(line, false)
	res, err := sb.parser.ParseLine(nline)
	if err != nil {
		return
	}
	num := ""
	val, isNum := res["num"]

	if isNum {
		num = val
		rule := NewRule(num,res["rule"], res["mode"])
		sb.skip = rule.Skip(sb.cfg)
		if !sb.skip {
			fmt.Printf("Rule  : %s\n", rule.String())
		}
		//sb.Results.AddRule(num, res["rule"], "level")
	} else {
		if !sb.skip {
			fmt.Printf("Result: %3s | %5s | %s\n", num, res["mode"], res["msg"])
		}
		// Add Instance
		//sb.Results.AddResult(num, res["mode"], res["msg"])
	}
	_ = num

}
func (sb *SecBench) parseLog(CntID string) {
	logOpts := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Tail: "any", Timestamps: false}
	reader, err := sb.cli.ContainerLogs(ctx, CntID, logOpts)
	if err != nil {
		msg := fmt.Sprintf("Error when connecting to log: %s", err.Error())
		log.Panic(msg)
		return
	}
	r := dlog.NewReader(reader)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		sb.parseLine(line)
	}
	err = scanner.Err()
	if err != nil {
		msg := fmt.Sprintf("Something went wrong going through the log: %s", err.Error())
		log.Panic(msg)
		return
	}
}

