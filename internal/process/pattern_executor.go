package process

/*
	Executes the registered patterns
*/

import (
	"time"
	"web-page-analyzer/internal/process/patterns"

	"github.com/sirupsen/logrus"
)

type PatternExecutor struct {
}

func NewPatternExecutor() *PatternExecutor {
	return &PatternExecutor{}
}

func (pe *PatternExecutor) Execute(ctx *patterns.Context, result map[string]any) {

	for _, p := range patterns.All() {

		start := time.Now()

		logrus.WithField("pattern", p.Name()).Info("Executing pattern")
		err := p.Apply(ctx, result)

		duration := time.Since(start)

		entry := logrus.WithFields(logrus.Fields{
			"pattern": p.Name(),
			"elapsed": duration,
		})

		if err != nil {
			entry.WithError(err).Warn("Pattern execution failed")
		} else {
			entry.Info("Pattern executed successfully")
		}
	}

}
