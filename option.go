package logger

type Option struct {
	IsEnable            bool
	EnableStackTrace    bool
	EnableMaskingFields bool
	MaskingFields       []string
}
