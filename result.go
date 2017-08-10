package secbench

type Result struct {
	// [WARN]      * Unencrypted overlay network: ingress (swarm)
	Mode 		string
	Description string
	Source      string
}


func NewResult(mode, desc string) Result {
	return Result{
		Mode: mode,
		Description: desc,
	}
}
