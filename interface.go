package secbench

type Parser interface {
	ParseLine(str string) (map[string]string, error)
}
