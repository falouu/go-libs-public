package api

type Specification struct {
	Requirements []Requirement
	// empty means <script dir>
	RootDir string
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
}

type RequirementEnvironment interface {
	BuildDir() string
}
