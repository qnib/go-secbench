# secbench
Tool to run and evaluate the Docker SecurityBenchmark.

## Description

This tool should help evaluate the output of the Docker Security Benchmark.
Idealy this script runs periodically and if it founds something, it exists with a non-zero errorcode.

A configuration should be provided, as to allow to configure what rule numbers should be considered.

## Motivation

The idea is to run this script to verify that everything is up-to-date after an deployment or periodically.<br>
By using `--fail-on-warn` the script will return a non-zero error-code if warnings occur.

The filter options can be used to trim down to the rules, which are applicable in a given environment.

### Usage

The script can be used to pipe the output of the docker benchmark to.

```bash
$ docker run -it --net host --pid host --cap-add audit_control \
      -e DOCKER_CONTENT_TRUST=$DOCKER_CONTENT_TRUST \
      -v /var/lib:/var/lib \
      -v /var/run/docker.sock:/var/run/docker.sock \
      -v /etc:/etc --label docker_bench_security \
      docker/docker-bench-security |./go-secbench_Darwin --piped --fail-on-warn
1    | INFO  || Host Configuration
1.1  | WARN  || Ensure a separate partition for containers has been created
*snip*
$ echo $?
25
```

Or it can spin up the container by itself.

```bash
$ ./go-secbench_Darwin --fail-on-warn
  2017-08-14 13:17:28.459777568 +0200 CEST > Pulling image 'docker/docker-bench-security'
  2017-08-14 13:17:30.044573308 +0200 CEST > Start container 'docker/docker-bench-security' as 'secbench-1502709448333315733'
  2017-08-14 13:17:30.266071968 +0200 CEST > Attaching to log-stream and parsing log
  2017-08-14 13:17:54.41579784 +0200 CEST > Removing container secbench-1502709448333315733
1    | INFO  || Host Configuration
1.1  | WARN  || Ensure a separate partition for containers has been created
*snip*
$ echo $?
25
```

## Only show WARN

```bash
$ cat resources/moby.out|./go-secbench_Darwin --piped --fail-on-warn --rule-numbers-only-regex "5.*"
  5    | INFO  || Container Runtime
  5.1  | WARN  || Ensure AppArmor Profile is Enabled
       | WARN  || No AppArmorProfile Found: http.1.9y98tks4opf10pgs4d1i24iyf
  5.10 | WARN  || Ensure memory usage for container is limited
       | WARN  || Container running without memory restrictions: http.1.9y98tks4opf10pgs4d1i24iyf
  5.11 | WARN  || Ensure CPU priority is set appropriately on the container
       | WARN  || Container running without CPU restrictions: http.1.9y98tks4opf10pgs4d1i24iyf
  5.12 | WARN  || Ensure the container's root filesystem is mounted as read only
       | WARN  || Container running with root FS mounted R/W: http.1.9y98tks4opf10pgs4d1i24iyf
  5.13 | PASS  || Ensure incoming container traffic is binded to a specific host interface
  5.14 | WARN  || Ensure 'on-failure' container restart policy is set to '5'
       | WARN  || MaximumRetryCount is not set to 5: http.1.9y98tks4opf10pgs4d1i24iyf
  5.15 | PASS  || Ensure the host's process namespace is not shared
  5.16 | PASS  || Ensure the host's IPC namespace is not shared
  5.17 | PASS  || Ensure host devices are not directly exposed to containers
  5.18 | INFO  || Ensure the default ulimit is overwritten at runtime, only if needed
       | INFO  || Container no default ulimit override: http.1.9y98tks4opf10pgs4d1i24iyf
  5.19 | PASS  || Ensure mount propagation mode is not set to shared
  5.2  | WARN  || Ensure SELinux security options are set, if applicable
       | WARN  || No SecurityOptions Found: http.1.9y98tks4opf10pgs4d1i24iyf
  5.20 | PASS  || Ensure the host's UTS namespace is not shared
  5.21 | PASS  || Ensure the default seccomp profile is not Disabled
  5.22 | NOTE  || Ensure docker exec commands are not used with privileged option
  5.23 | NOTE  || Ensure docker exec commands are not used with user option
  5.24 | PASS  || Ensure cgroup usage is confirmed
  5.25 | WARN  || Ensure the container is restricted from acquiring additional privileges
       | WARN  || Privileges not restricted: http.1.9y98tks4opf10pgs4d1i24iyf
  5.26 | PASS  || Ensure container health is checked at runtime
  5.27 | INFO  || Ensure docker commands always get the latest version of the image
  5.28 | WARN  || Ensure PIDs cgroup limit is used
       | WARN  || PIDs limit not set: http.1.9y98tks4opf10pgs4d1i24iyf
  5.29 | INFO  || Ensure Docker's default bridge docker0 is not used
       | INFO  || Container in docker0 network:
  5.3  | PASS  || Ensure Linux Kernel Capabilities are restricted within containers
  5.30 | PASS  || Ensure the host's user namespaces is not shared
  5.31 | PASS  || Ensure the Docker socket is not mounted inside any containers
  5.4  | PASS  || Ensure privileged containers are not used
  5.5  | PASS  || Ensure sensitive host system directories are not mounted on containers
  5.6  | PASS  || Ensure ssh is not run within containers
  5.7  | PASS  || Ensure privileged ports are not mapped within containers
  5.8  | NOTE  || Ensure only needed ports are open on the container
  5.9  | PASS  || Ensure the host's network namespace is not shared
$ echo $?
  8
```


### Only show warning rules

In order to ignore rules with a given status, use `--modes-ignore`.

```bash
$ cat resources/moby.out |./go-secbench_Darwin --piped --modes-ignore "INFO,NOTE,PASS"
  1.1  | WARN  || Ensure a separate partition for containers has been created
  1.5  | WARN  || Ensure auditing is configured for the Docker daemon
  1.6  | WARN  || Ensure auditing is configured for Docker files and directories - /var/lib/docker
  2.1  | WARN  || Ensure network traffic is restricted between containers on the default bridge
  2.11 | WARN  || Ensure that authorization for Docker client commands is enabled
  2.12 | WARN  || Ensure centralized and remote logging is configured
  2.13 | WARN  || Ensure operations on legacy registry (v1) are Disabled
  2.15 | WARN  || Ensure Userland Proxy is Disabled
  2.17 | WARN  || Ensure experimental features are avoided in production
  2.8  | WARN  || Enable user namespace support
  4.1  | WARN  || Ensure a user for the container has been created
       | WARN  || Running as root: http.1.9y98tks4opf10pgs4d1i24iyf
  4.5  | WARN  || Ensure Content trust for Docker is Enabled
  4.6  | WARN  || Ensure HEALTHCHECK instructions have been added to the container image
       | WARN  || No Healthcheck found: [dock.gaikai.org/gaikai/docker-visualizer-client:2017-08-03]
       | WARN  || No Healthcheck found: [gocd/docker-gocd-agent-alpine-3.5:latest]
       | WARN  || No Healthcheck found: [dock.gaikai.org/snaik/docker-visualizer-client:20170801]
       | WARN  || No Healthcheck found: [alpine:edge]
       | WARN  || No Healthcheck found: [qnib/docker-burrow:latest]
       | WARN  || No Healthcheck found: [test:latest]
       | WARN  || No Healthcheck found: [qnib/docker-plugin-metrics-opentsdb:latest]
       | WARN  || No Healthcheck found: [dock.gaikai.org/wrex/big-fat-whale-game:latest]
       | WARN  || No Healthcheck found: [dock.gaikai.org/htonnies/kafka-websocket:0.0.1]
       | WARN  || No Healthcheck found: [qnibp/dfile:latest]
       | WARN  || No Healthcheck found: [qnib/dfile:latest]
       | WARN  || No Healthcheck found: [dockervisualizerclient_client:latest]
       | WARN  || No Healthcheck found: [dockervisualizerbackend_backend:latest]
       | WARN  || No Healthcheck found: [dock.gaikai.org/peifler/docker-visualizer-backend:latest]
       | WARN  || No Healthcheck found: [node:boron]
       | WARN  || No Healthcheck found: [dock.gaikai.org/wrex/docker-plugin-metrics-opentsdb:latest]
       | WARN  || No Healthcheck found: [qnib/prom2json:latest]
       | WARN  || No Healthcheck found: [qnib/docker-metrics-plugin-opentsdb:latest]
       | WARN  || No Healthcheck found: [minio/doctor:latest]
       | WARN  || No Healthcheck found: [gocd/gocd-agent-alpine-3.5:v17.7.0]
       | WARN  || No Healthcheck found: [gocd/gocd-server:v17.7.0]
       | WARN  || No Healthcheck found: [golang:alpine]
       | WARN  || No Healthcheck found: [qnib/alpn-consul:latest]
       | WARN  || No Healthcheck found: [qnib/alpn-base:latest]
       | WARN  || No Healthcheck found: [alpine:latest]
       | WARN  || No Healthcheck found: [alpine:3.5]
       | WARN  || No Healthcheck found: [postgres:latest]
       | WARN  || No Healthcheck found: [golang:latest]
       | WARN  || No Healthcheck found: [ubuntu:latest]
       | WARN  || No Healthcheck found: [ubuntu:14.04]
       | WARN  || No Healthcheck found: [busybox:latest]
       | WARN  || No Healthcheck found: [dock.gaikai.org/mayflower/docker-ls:latest]
       | WARN  || No Healthcheck found: [dock.gaikai.org/gaikai/docker-visualizer-backend:latest]
       | WARN  || No Healthcheck found: [ilagnev/alpine-nginx-lua:latest]
  5.1  | WARN  || Ensure AppArmor Profile is Enabled
       | WARN  || No AppArmorProfile Found: http.1.9y98tks4opf10pgs4d1i24iyf
  5.10 | WARN  || Ensure memory usage for container is limited
       | WARN  || Container running without memory restrictions: http.1.9y98tks4opf10pgs4d1i24iyf
  5.11 | WARN  || Ensure CPU priority is set appropriately on the container
       | WARN  || Container running without CPU restrictions: http.1.9y98tks4opf10pgs4d1i24iyf
  5.12 | WARN  || Ensure the container's root filesystem is mounted as read only
       | WARN  || Container running with root FS mounted R/W: http.1.9y98tks4opf10pgs4d1i24iyf
  5.14 | WARN  || Ensure 'on-failure' container restart policy is set to '5'
       | WARN  || MaximumRetryCount is not set to 5: http.1.9y98tks4opf10pgs4d1i24iyf
  5.2  | WARN  || Ensure SELinux security options are set, if applicable
       | WARN  || No SecurityOptions Found: http.1.9y98tks4opf10pgs4d1i24iyf
  5.25 | WARN  || Ensure the container is restricted from acquiring additional privileges
       | WARN  || Privileges not restricted: http.1.9y98tks4opf10pgs4d1i24iyf
  5.28 | WARN  || Ensure PIDs cgroup limit is used
       | WARN  || PIDs limit not set: http.1.9y98tks4opf10pgs4d1i24iyf
  7.1  | WARN  || Ensure swarm mode is not Enabled, if not needed
  7.3  | WARN  || Ensure swarm services are binded to a specific host interface
  7.4  | WARN  || Ensure data exchanged between containers are encrypted on different nodes on the overlay network
       | WARN  || Unencrypted overlay network: ingress (swarm)
  7.6  | WARN  || Ensure swarm manager is run in auto-lock mode
```
