package types

type ConfigUtils struct {
	STORAGE_PATH string
	JWT          string
	IMAGE_TEST   string
	PROJECT_PATH string
	DATABASE     []any
}

type Bootsterap struct {
	Config   ConfigUtils
	Database []any
}
