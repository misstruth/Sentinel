package workflow

// Status 工作流状态
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Step 工作流步骤
type Step struct {
	Name   string
	Action func() error
}

// Workflow 工作流
type Workflow struct {
	Name   string
	Steps  []*Step
	Status Status
}

// NewWorkflow 创建工作流
func NewWorkflow(name string) *Workflow {
	return &Workflow{
		Name:   name,
		Steps:  make([]*Step, 0),
		Status: StatusPending,
	}
}

// AddStep 添加步骤
func (w *Workflow) AddStep(s *Step) {
	w.Steps = append(w.Steps, s)
}

// Run 执行工作流
func (w *Workflow) Run() error {
	w.Status = StatusRunning
	for _, step := range w.Steps {
		if err := step.Action(); err != nil {
			w.Status = StatusFailed
			return err
		}
	}
	w.Status = StatusCompleted
	return nil
}
