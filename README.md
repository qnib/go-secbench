# secbench
Tool to run and evaluate the Docker SecurityBenchmark.

## Description

This tool should help evaluate the output of the Docker Security Benchmark.
Idealy this script runs periodically and if it founds something, it exists with a non-zero errorcode.

A configuration should be provided, as to allow to configure what rule numbers should be considered.

## Only show WARN

```
$ go run main.go --modes-show WARN
Rule  : 1.1  | Cur:WARN  || Ensure a separate partition for containers has been created
Rule  : 1.5  | Cur:WARN  || Ensure auditing is configured for the Docker daemon
Rule  : 1.6  | Cur:WARN  || Ensure auditing is configured for Docker files and directories - /var/lib/docker
Rule  : 2.1  | Cur:WARN  || Ensure network traffic is restricted between containers on the default bridge
Rule  : 2.8  | Cur:WARN  || Enable user namespace support
Rule  : 2.11 | Cur:WARN  || Ensure that authorization for Docker client commands is enabled
Rule  : 2.12 | Cur:WARN  || Ensure centralized and remote logging is configured
Rule  : 2.13 | Cur:WARN  || Ensure operations on legacy registry (v1) are Disabled
Rule  : 2.15 | Cur:WARN  || Ensure Userland Proxy is Disabled
Rule  : 2.17 | Cur:WARN  || Ensure experimental features are avoided in production
Rule  : 4.1  | Cur:WARN  || Ensure a user for the container has been created
Result:     |  WARN | Running as root: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 4.5  | Cur:WARN  || Ensure Content trust for Docker is Enabled
Rule  : 5.1  | Cur:WARN  || Ensure AppArmor Profile is Enabled
Result:     |  WARN | No AppArmorProfile Found: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.2  | Cur:WARN  || Ensure SELinux security options are set, if applicable
Result:     |  WARN | No SecurityOptions Found: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.10 | Cur:WARN  || Ensure memory usage for container is limited
Result:     |  WARN | Container running without memory restrictions: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.11 | Cur:WARN  || Ensure CPU priority is set appropriately on the container
Result:     |  WARN | Container running without CPU restrictions: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.12 | Cur:WARN  || Ensure the container's root filesystem is mounted as read only
Result:     |  WARN | Container running with root FS mounted R/W: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.14 | Cur:WARN  || Ensure 'on-failure' container restart policy is set to '5'
Result:     |  WARN | MaximumRetryCount is not set to 5: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.25 | Cur:WARN  || Ensure the container is restricted from acquiring additional privileges
Result:     |  WARN | Privileges not restricted: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 5.28 | Cur:WARN  || Ensure PIDs cgroup limit is used
Result:     |  WARN | PIDs limit not set: http.1.i1ck6akamfe9dkpr86ey3cxwa
Rule  : 7.1  | Cur:WARN  || Ensure swarm mode is not Enabled, if not needed
Rule  : 7.3  | Cur:WARN  || Ensure swarm services are binded to a specific host interface
Rule  : 7.4  | Cur:WARN  || Ensure data exchanged between containers are encrypted on different nodes on the overlay network
Result:     |  WARN | Unencrypted overlay network: ingress (swarm)
Rule  : 7.6  | Cur:WARN  || Ensure swarm manager is run in auto-lock mode
```

### Narrow down the Rules

In order to narrow down certain rules, use `--rule-numbers-only`.

```bash
$ go run main.go --modes-ignore "NOTE,PASS,INFO" --quiet --rule-numbers-only 5.28,7.4
  5.28 | WARN  || Ensure PIDs cgroup limit is used
       | WARN  || PIDs limit not set: http.1.dk05f9d6nl81dnongce4msvja
  7.4  | WARN  || Ensure data exchanged between containers are encrypted on different nodes on the overlay network
       | WARN  || Unencrypted overlay network: ingress (swarm)
```
