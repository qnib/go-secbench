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
)

const (
	imageName = "docker/docker-bench-security"
)

var (
	ctx = context.Background()
)

type SecBench struct {
	cli *client.Client
	grok Grok
}

func NewSecBenc() (sb *SecBench, err error) {
	c, err := client.NewEnvClient()
	if err != nil {
		return
	}
	g := NewGrok()
	sb = &SecBench{
		cli: c,
		grok: g,
	}
	return
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
	defer sb.removeContainer(container.ID)
	// Parse log
	sb.parseLog(container.ID)

}

func (sb *SecBench) removeContainer(cid string) {
	err := sb.cli.ContainerRemove(ctx, cid, types.ContainerRemoveOptions{})
	if err != nil {
		log.Panic(err.Error())
	}
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
		nline := vtclean.Clean(line, false)
		res, err := sb.grok.ParseLine(nline)
		if err != nil {
			//fmt.Printf("!! %s > %v\n", err.Error(), nline)
			continue
		}
		fmt.Printf("%v\n", res)

	}
	err = scanner.Err()
	if err != nil {
		msg := fmt.Sprintf("Something went wrong going through the log: %s", err.Error())
		log.Panic(msg)
		return
	}
}
