package anthropic

import "fmt"

type Response struct {
	Role         string  `json:"role"`
	Content      string  `json:"content"`
	Type         string  `json:"type"`
	InputTokens  float64 `json:"input_tokens"`
	OutputTokens float64 `json:"output_tokens"`
	TotalTokens  float64 `json:"total_tokens"`
}

func ParseResponse(respData map[string]interface{}) (Response, error) {
	responseType, ok := respData["content"].([]interface{})[0].(map[string]interface{})["type"].(string)
	if !ok {
		return Response{}, fmt.Errorf("failed to parse response type")
	}

	if responseType != "text" {
		return Response{}, fmt.Errorf("response is not of type text")
	}

	response, ok := respData["content"].([]interface{})[0].(map[string]interface{})["text"].(string)
	if !ok {
		return Response{}, fmt.Errorf("failed to parse response content")
	}

	role, ok := respData["role"].(string)
	if !ok {
		return Response{}, fmt.Errorf("failed to parse response role")
	}

	inputTokens, ok := respData["usage"].(map[string]interface{})["input_tokens"].(float64)
	if !ok {
		return Response{}, fmt.Errorf("failed to parse input tokens")
	}

	outputTokens, ok := respData["usage"].(map[string]interface{})["output_tokens"].(float64)
	if !ok {
		return Response{}, fmt.Errorf("failed to parse output tokens")
	}

	totalTokens := inputTokens + outputTokens

	return Response{
		Role:         role,
		Content:      response,
		Type:         responseType,
		InputTokens:  inputTokens,
		OutputTokens: outputTokens,
		TotalTokens:  totalTokens,
	}, nil
}

func (r *Response) GetMessageRequest() map[string]interface{} {
	return map[string]interface{}{
		"role":    r.Role,
		"content": r.Content,
	}
}
