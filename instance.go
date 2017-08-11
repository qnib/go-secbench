package secbench

import "fmt"

type Instance struct {
	// [WARN]      * Unencrypted overlay network: ingress (swarm)
	Mode 		string
	Description string
	Source      string
}


func NewInstance(mode, desc string) Instance {
	return Instance{
		Mode: mode,
		Description: desc,
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("     | %-5s || %s", i.Mode, i.Description)
}
