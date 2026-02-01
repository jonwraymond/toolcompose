package skill

import (
	"context"
	"fmt"
	"testing"
)

type benchRunner struct{}

func (r *benchRunner) Run(ctx context.Context, step Step) (any, error) {
	return step.ID, nil
}

func BenchmarkPlanner_Plan(b *testing.B) {
	for _, size := range []int{10, 50, 100} {
		steps := make([]Step, size)
		for i := range steps {
			steps[i] = Step{
				ID:     fmt.Sprintf("step-%03d", i),
				ToolID: fmt.Sprintf("ns:tool%d", i),
				Inputs: map[string]any{"key": "value"},
			}
		}
		skill := Skill{Name: "bench", Steps: steps}
		planner := NewPlanner()

		b.Run(fmt.Sprintf("steps_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = planner.Plan(skill)
			}
		})
	}
}

func BenchmarkExecute(b *testing.B) {
	for _, size := range []int{10, 50, 100} {
		steps := make([]Step, size)
		for i := range steps {
			steps[i] = Step{
				ID:     fmt.Sprintf("step-%03d", i),
				ToolID: fmt.Sprintf("ns:tool%d", i),
			}
		}
		plan := Plan{Name: "bench", Steps: steps}
		runner := &benchRunner{}

		b.Run(fmt.Sprintf("steps_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = Execute(context.Background(), plan, runner)
			}
		})
	}
}

func BenchmarkSkill_Validate(b *testing.B) {
	steps := make([]Step, 50)
	for i := range steps {
		steps[i] = Step{
			ID:     fmt.Sprintf("step-%03d", i),
			ToolID: fmt.Sprintf("ns:tool%d", i),
		}
	}
	skill := Skill{Name: "bench", Steps: steps}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = skill.Validate()
	}
}

func BenchmarkMaxStepsGuard(b *testing.B) {
	steps := make([]Step, 50)
	for i := range steps {
		steps[i] = Step{ID: fmt.Sprintf("%d", i), ToolID: "ns:tool"}
	}
	skill := Skill{Name: "bench", Steps: steps}
	guard := MaxStepsGuard(100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = guard.Validate(skill)
	}
}

func BenchmarkAllowedToolIDsGuard(b *testing.B) {
	allowed := make([]string, 100)
	for i := range allowed {
		allowed[i] = fmt.Sprintf("ns:tool%d", i)
	}

	steps := make([]Step, 50)
	for i := range steps {
		steps[i] = Step{ID: fmt.Sprintf("%d", i), ToolID: fmt.Sprintf("ns:tool%d", i)}
	}
	skill := Skill{Name: "bench", Steps: steps}
	guard := AllowedToolIDsGuard(allowed)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = guard.Validate(skill)
	}
}
