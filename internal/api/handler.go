package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/clevertechware/demo-renovate-basic-service-plan/internal/plan"
)

const svgSuffix = ".svg"

type PlanFinder interface {
	FindByCode(ctx context.Context, code string) (plan.Plan, error)
}

func NewRouter(plans PlanFinder) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/healthz", healthz)
	router.GET("/plans/:file", planSVG(plans))
	return router
}

func healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func planSVG(plans PlanFinder) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, ok := strings.CutSuffix(c.Param("file"), svgSuffix)
		if !ok || code == "" {
			c.Status(http.StatusNotFound)
			return
		}

		p, err := plans.FindByCode(c.Request.Context(), code)
		switch {
		case errors.Is(err, plan.ErrNotFound):
			c.Status(http.StatusNotFound)
		case err != nil:
			slog.ErrorContext(c.Request.Context(), "find plan", "code", code, "error", err)
			c.Status(http.StatusInternalServerError)
		default:
			c.Data(http.StatusOK, "image/svg+xml", []byte(p.SVG()))
		}
	}
}
