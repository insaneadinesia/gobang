package logger

type Option struct {
	IsEnable         bool
	EnableStackTrace bool
	MaskingFields    []string
}
