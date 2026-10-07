package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"skintrace/internal/deepseek"
	"skintrace/internal/record"
)

const HistoryToolName = "get_skin_history"

func HistoryToolDefinition() deepseek.Tool {
	return deepseek.Tool{
		Type: "function",
		Function: deepseek.FunctionDefinition{
			Name: HistoryToolName,
			Description: "查询最近的皮肤记录，按创建时间倒序返回。" +
				"包含记录信息和用户备注，不读取或分析照片。" +
				"limit 表示记录条数，不是天数。",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit": map[string]any{
						"type":        "integer",
						"description": "查询最近多少条记录，范围为 1 到 30",
						"minimum":     1,
						"maximum":     30,
					},
				},
				"required":             []string{"limit"},
				"additionalProperties": false,
			},
		},
	}
}

func GetSkinHistory(
	ctx context.Context,
	records *record.Service,
	arguments string,
) (string, error) {
	var args struct {
		Limit int `json:"limit"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("parse history tool arguments: %w", err)
	}

	result, err := records.Recent(ctx, args.Limit)
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("encode history tool result: %w", err)
	}

	return string(data), nil
}
