package vertexai

import (
	"github.com/gin-gonic/gin"
	claude "github.com/songquanpeng/one-api/relay/adaptor/vertexai/claude"
	gemini "github.com/songquanpeng/one-api/relay/adaptor/vertexai/gemini"
	"github.com/songquanpeng/one-api/relay/adaptor/vertexai/imagen"
	"github.com/songquanpeng/one-api/relay/billing/ratio"
	"github.com/songquanpeng/one-api/relay/meta"
	"github.com/songquanpeng/one-api/relay/model"
	"net/http"
)

type VertexAIModelType int

const (
	VertexAIClaude VertexAIModelType = iota + 1
	VertexAIGemini
	VertexAIImagen
)

var modelMapping = map[string]VertexAIModelType{}

func init() {
	for model := range claude.RatioMap {
		modelMapping[model] = VertexAIClaude
	}

	for model := range gemini.RatioMap {
		modelMapping[model] = VertexAIGemini
	}

	for model := range imagen.RatioMap {
		modelMapping[model] = VertexAIImagen
	}
}

type innerAIAdapter interface {
	ConvertRequest(c *gin.Context, relayMode int, request *model.GeneralOpenAIRequest) (any, error)
	ConvertImageRequest(c *gin.Context, request *model.ImageRequest) (any, error)
	DoResponse(c *gin.Context, resp *http.Response, meta *meta.Meta) (usage *model.Usage, err *model.ErrorWithStatusCode)
	GetRatio(meta *meta.Meta) *ratio.Ratio
}

func GetAdaptor(model string) innerAIAdapter {
	adaptorType := modelMapping[model]
	switch adaptorType {
	case VertexAIClaude:
		return &claude.Adaptor{}
	case VertexAIGemini:
		return &gemini.Adaptor{}
	case VertexAIImagen:
		return &imagen.Adaptor{}
	default:
		return nil
	}
}
