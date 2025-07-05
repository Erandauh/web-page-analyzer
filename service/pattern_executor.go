package analyze_service

/*
	Executes the registered patterns
*/

import (
	"log"

	"web-page-analyzer/process/patterns"
)

func Execute(ctx *patterns.Context, result map[string]any) {

	for _, p := range patterns.All() {
		err := p.Apply(ctx, result)
		if err != nil {
			log.Printf("[WARN] Pattern %s failed: %v", p.Name(), err)
		}
	}

}
