package delivery

import (
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/perlinleo/vision/internal/domain"
	"github.com/perlinleo/vision/internal/middleware"
	"github.com/valyala/fasthttp"
)

type passHandler struct {
	PassUsecase domain.PassUsecase
	UserUsecase domain.UserUsecase
}

func NewPassesHandler(router *router.Router, usecase domain.PassUsecase, user domain.UserUsecase, su domain.SessionUsecase) {
	handler := &passHandler{
		PassUsecase: usecase,
		UserUsecase: user,
	}

	router.GET("/api/v1/passes", middleware.Cors(
		middleware.ReponseMiddlwareAndLogger(
			middleware.Auth(
				middleware.ReponseMiddlwareAndLogger(handler.Passes), su))))

}

func (h *passHandler) Passes(ctx *fasthttp.RequestCtx) {
	aid := ctx.UserValue("AID").(*domain.UserSession)
	passes, err := h.PassUsecase.GetUserPasses(aid.UserID)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
	}

	ctxBody, err := json.Marshal(passes)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}

	ctx.SetBody(ctxBody)
}
