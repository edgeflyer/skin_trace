package agent

import (
	"context"
	"fmt"
	"log"

	"skintrace/internal/agent/tools"
	"skintrace/internal/deepseek"
)

func (a *Agent) executeTool(
	ctx context.Context,
	call deepseek.ToolCall,
) (string, error) {
	log.Println("执行工具")
	switch call.Function.Name {
	case tools.HistoryToolName:
		return tools.GetSkinHistory(
			ctx,
			a.records,
			call.Function.Arguments,
		)
	case tools.PhotoToolName:
		return tools.AnalyzeSkinPhoto(
			ctx,
			a.records,
			a.client,
			call.Function.Arguments,
		)
	default:
		return "", fmt.Errorf("unknown tool: %s", call.Function.Name)
	}
}
