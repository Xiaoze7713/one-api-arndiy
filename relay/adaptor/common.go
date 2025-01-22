package adaptor

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/songquanpeng/one-api/common/client"
	"github.com/songquanpeng/one-api/common/ctxkey"
	"github.com/songquanpeng/one-api/relay/billing/ratio"
	"github.com/songquanpeng/one-api/relay/meta"
)

func SetupCommonRequestHeader(c *gin.Context, req *http.Request, meta *meta.Meta) {
	req.Header.Set("Content-Type", c.Request.Header.Get("Content-Type"))
	req.Header.Set("Accept", c.Request.Header.Get("Accept"))
	if meta.IsStream && c.Request.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "text/event-stream")
	}
}

func DoRequestHelper(a Adaptor, c *gin.Context, meta *meta.Meta, requestBody io.Reader) (*http.Response, error) {
	fullRequestURL, err := a.GetRequestURL(meta)
	if err != nil {
		return nil, errors.Wrap(err, "get request url failed")
	}

	req, err := http.NewRequestWithContext(c.Request.Context(),
		c.Request.Method, fullRequestURL, requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "new request failed")
	}

	req.Header.Set("Content-Type", c.GetString(ctxkey.ContentType))

	err = a.SetupRequestHeader(c, req, meta)
	if err != nil {
		return nil, errors.Wrap(err, "setup request header failed")
	}
	resp, err := DoRequest(c, req)
	if err != nil {
		return nil, errors.Wrap(err, "do request failed")
	}
	return resp, nil
}

func DoRequest(c *gin.Context, req *http.Request) (*http.Response, error) {
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.New("resp is nil")
	}
	_ = req.Body.Close()
	_ = c.Request.Body.Close()

	return resp, nil
}

func GetRatioHelper(meta *meta.Meta, ratioMap map[string]ratio.Ratio) *ratio.Ratio {
	var result ratio.Ratio
	if ratio, ok := ratioMap[meta.OriginModelName]; ok {
		result = ratio
		return &result
	}
	if ratio, ok := ratioMap[meta.ActualModelName]; ok {
		result = ratio
		return &result
	}
	return nil
}

func GetModelListHelper(ratioMap map[string]ratio.Ratio) []string {
	var modelList []string
	for model := range ratioMap {
		modelList = append(modelList, model)
	}
	return modelList
}
