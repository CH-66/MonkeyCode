package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/GoYoko/web"

	"github.com/chaitin/MonkeyCode/backend/config"
	"github.com/chaitin/MonkeyCode/backend/domain"
	"github.com/chaitin/MonkeyCode/backend/internal/middleware"
	"github.com/chaitin/MonkeyCode/backend/internal/proxy"
)

type V1Handler struct {
	logger   *slog.Logger
	proxy    *proxy.LLMProxy
	proxyUse domain.ProxyUsecase
	usecase  domain.OpenAIUsecase
	euse     domain.ExtensionUsecase
	uuse     domain.UserUsecase
	config   *config.Config
}

func NewV1Handler(
	logger *slog.Logger,
	w *web.Web,
	proxy *proxy.LLMProxy,
	proxyUse domain.ProxyUsecase,
	usecase domain.OpenAIUsecase,
	euse domain.ExtensionUsecase,
	uuse domain.UserUsecase,
	middleware *middleware.ProxyMiddleware,
	active *middleware.ActiveMiddleware,
	config *config.Config,
) *V1Handler {
	h := &V1Handler{
		logger:   logger.With(slog.String("handler", "openai")),
		proxy:    proxy,
		proxyUse: proxyUse,
		usecase:  usecase,
		euse:     euse,
		uuse:     uuse,
		config:   config,
	}

	w.GET("/api/config", web.BindHandler(h.GetConfig), middleware.Auth())
	w.GET("/v1/version", web.BaseHandler(h.Version), middleware.Auth())
	w.GET("/v1/health", web.BaseHandler(h.HealthCheck))

	g := w.Group("/v1", middleware.Auth())
	g.GET("/models", web.BaseHandler(h.ModelList))
	g.POST("/completion/accept", web.BindHandler(h.AcceptCompletion), active.Active("apikey"))
	g.POST("/report", web.BindHandler(h.Report), active.Active("apikey"))
	g.POST("/chat/completions", web.BaseHandler(h.ChatCompletion), active.Active("apikey"))
	g.POST("/completions", web.BaseHandler(h.Completions), active.Active("apikey"))
	g.POST("/embeddings", web.BaseHandler(h.Embeddings), active.Active("apikey"))
	return h
}

func BadRequest(c *web.Context, msg string) error {
	c.JSON(http.StatusBadRequest, echo.Map{
		"error": echo.Map{
			"message": msg,
			"type":    "invalid_request_error",
		},
	})
	return nil
}

func (h *V1Handler) Version(c *web.Context) error {
	v, err := h.euse.Latest(c.Request().Context())
	if err != nil {
		return err
	}

	s, err := h.uuse.GetSetting(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, domain.VersionInfo{
		Version: v.Version,
		URL:     fmt.Sprintf("%s/api/v1/static/vsix/%s", h.config.GetBaseURL(c.Request(), s), v.Version),
	})
}

// AcceptCompletion 接受补全请求
//
//	@Tags			OpenAIV1
//	@Summary		接受补全请求
//	@Description	接受补全请求
//	@ID				accept-completion
//	@Accept			json
//	@Produce		json
//	@Param			param	body		domain.AcceptCompletionReq	true	"补全请求"
//	@Success		200		{object}	web.Resp{}
//	@Router			/v1/completion/accept [post]
func (h *V1Handler) AcceptCompletion(c *web.Context, req domain.AcceptCompletionReq) error {
	if err := h.proxyUse.AcceptCompletion(c.Request().Context(), &req); err != nil {
		return BadRequest(c, err.Error())
	}
	return nil
}

// Report 报告
//
//	@Tags			OpenAIV1
//	@Summary		报告
//	@Description	报告，支持多种操作：accept（接受补全）、suggest（建议）、reject（拒绝补全并记录用户输入）、file_written（文件写入）
//	@ID				report
//	@Accept			json
//	@Produce		json
//	@Param			param	body		domain.ReportReq	true	"报告请求"
//	@Success		200		{object}	web.Resp{}
//	@Router			/v1/report [post]
func (h *V1Handler) Report(c *web.Context, req domain.ReportReq) error {
	h.logger.DebugContext(c.Request().Context(), "Report", slog.Any("req", req))
	req.UserID = middleware.GetApiKey(c).UserID
	if err := h.proxyUse.Report(c.Request().Context(), &req); err != nil {
		return err
	}
	return c.Success(nil)
}

// ModelList 模型列表
//
//	@Tags			OpenAIV1
//	@Summary		模型列表
//	@Description	模型列表
//	@ID				model-list
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{data=domain.ModelListResp}
//	@Router			/v1/models [get]
func (h *V1Handler) ModelList(c *web.Context) error {
	resp, err := h.usecase.ModelList(c.Request().Context())
	if err != nil {
		return err
	}
	return c.Success(resp)
}

// ChatCompletion 处理聊天补全请求
//
//	@Tags			OpenAIV1
//	@Summary		处理聊天补全请求
//	@Description	处理聊天补全请求
//	@ID				chat-completion
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/chat/completions [post]
func (h *V1Handler) ChatCompletion(c *web.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Completions 处理文本补全请求
//
//	@Tags			OpenAIV1
//	@Summary		处理文本补全请求
//	@Description	处理文本补全请求
//	@ID				completions
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/completions [post]
func (h *V1Handler) Completions(c *web.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

// Embeddings 处理嵌入请求
//
//	@Tags			OpenAIV1
//	@Summary		处理嵌入请求
//	@Description	处理嵌入请求
//	@ID				embeddings
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/embeddings [post]
func (h *V1Handler) Embeddings(c *web.Context) error {
	h.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}

func (h *V1Handler) GetConfig(c *web.Context, req domain.ConfigReq) error {
	key := middleware.GetApiKey(c)
	s, err := h.uuse.GetSetting(c.Request().Context())
	if err != nil {
		return err
	}
	req.Key = key.Key
	req.BaseURL = h.config.GetBaseURL(c.Request(), s)
	resp, err := h.usecase.GetConfig(c.Request().Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, resp)
}

// HealthCheck 健康检查
//
//	@Tags			OpenAIV1
//	@Summary		健康检查
//	@Description	固定回包 `{"code": 0, "data": "MonkeyCode"}`
//	@ID				health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	web.Resp{}
//	@Router			/v1/health [get]
func (h *V1Handler) HealthCheck(c *web.Context) error {
	return c.Success("MonkeyCode")
}
