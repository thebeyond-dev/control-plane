package launcher

import (
	"context"
	"net/http"
	"time"

	"github.com/go-telegram/bot"
	"go.uber.org/zap"
)

type WebHook struct {
	bot       *bot.Bot
	appLogger *zap.Logger
	port      string
}

func NewWebHook(
	bot *bot.Bot,
	appLogger *zap.Logger,
	port string,
) *WebHook {
	return &WebHook{bot, appLogger, port}
}

func (wh *WebHook) Launch(ctx context.Context) {
	srv := &http.Server{
		Addr:         ":" + wh.port,
		Handler:      wh.bot.WebhookHandler(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		<-ctx.Done()
		wh.appLogger.Info("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			wh.appLogger.Error("Server shutdown error", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	go wh.bot.StartWebhook(ctx)

	wh.appLogger.Info("Listening on :" + wh.port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		wh.appLogger.Fatal("Failed to start server", zap.Error(err))
	}

	<-idleConnsClosed
	wh.appLogger.Info("Server exited properly")
}
