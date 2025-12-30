package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikiasgoitom/RevProx/internal/contract"
	"github.com/mikiasgoitom/RevProx/internal/domain/entity"
)

type ProxyHandler struct {
	proxyUsecase contract.IProxyUseCase
	logger       contract.ILogger
}

func NewProxyHandler(proxyUC contract.IProxyUseCase, logger contract.ILogger) *ProxyHandler {
	return &ProxyHandler{proxyUsecase: proxyUC, logger: logger}
}

func (h *ProxyHandler) HandleProxy(c *gin.Context) {
    // Translate gin.Context to your domain's RequestModel
    body, err := io.ReadAll(c.Request.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read request body"})
        return
    }

    // Create a new URL object and only populate it with the path and query
    // that should be sent to the origin server.
    // We get the path from the wildcard parameter, which strips the prefix.
    originPath := c.Param("path")
    // Preserve the original query string
    rawQuery := c.Request.URL.RawQuery

    // Rebuild a clean URL for the use case
    originURL := &url.URL{
        Path:     originPath,
        RawQuery: rawQuery,
    }

    reqModel := entity.RequestModel{
        Method:   c.Request.Method,
        URL:      originURL, 
        Headers:  c.Request.Header,
        Body:     body,
        ClientIP: c.ClientIP(),
    }

    // Call the proxy use case
    respModel, err := h.proxyUsecase.ServeProxyRequest(c.Request.Context(), reqModel)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service error", "details": err.Error()})
        return
    }

    // Write the ResponseModel back to the client
    // Copy headers from the response model to the Gin response
    for key, values := range respModel.Headers {
        // Do not copy the "Content-Encoding" header if it's "gzip",
        // as Gin handles compression automatically.
        if key == "Content-Encoding" && strings.Contains(values[0], "gzip") {
            continue
        }
        for _, value := range values {
            c.Writer.Header().Add(key, value)
        }
    }
    c.Data(respModel.Status, respModel.Headers.Get("Content-Type"), respModel.Body)
}
