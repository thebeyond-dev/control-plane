package unleash

import (
	"net/http"
	"time"

	"github.com/Unleash/unleash-go-sdk/v6"
)

const defaultEnvironment = "development"

func Init(appName, environment, url, token string) error {
	if environment == "" {
		environment = defaultEnvironment
	}

	if err := unleash.Initialize(
		unleash.WithAppName(appName),
		unleash.WithEnvironment(environment),
		unleash.WithUrl(url),
		unleash.WithCustomHeaders(http.Header{"Authorization": {token}}),
		unleash.WithRefreshInterval(30*time.Second),
	); err != nil {
		return err
	}

	ready := make(chan struct{})
	go func() {
		unleash.WaitForReady()
		close(ready)
	}()

	select {
	case <-ready:
	case <-time.After(3 * time.Second):
	}

	return nil
}

func Close() {
	unleash.Close()
}
