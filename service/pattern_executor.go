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

		log.Printf("[INFO] Executing pattern: %s", p.Name())
		err := p.Apply(ctx, result)
		log.Printf("[INFO] Executed pattern: %s", p.Name())

		if err != nil {
			log.Printf("[WARN] Pattern %s failed: %v", p.Name(), err)
		}
	}

}
