package api

type Specification struct {
	Requirements []Requirement
	// empty means <script dir>
	RootDir string
	// Console prompt. If empty, default will be used
	PromptText string
}

type Requirement interface {
	Info() *RequirementInfo
	Ensure(env RequirementEnvironment) (*RequirementResult, error)
}
type RequirementInfo struct {
	Description string
}

type RequirementResult struct {
	AddToPath []string
	// template to append to init file, Args will be available under {{.Extra.<key>}}
	BashRcAppend *Template
}

// Template passed to text/template
type Template struct {
	Template string
	Args     map[string]any
}

type RequirementEnvironment interface {
	BuildDir() string
	RootDir() string
}
