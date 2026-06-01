package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"pansou/config"
)

func TestSnakeGameRoutesArePublic(t *testing.T) {
	if config.AppConfig == nil {
		config.Init()
	}

	originalAuthEnabled := config.AppConfig.AuthEnabled
	defer func() { config.AppConfig.AuthEnabled = originalAuthEnabled }()
	config.AppConfig.AuthEnabled = true
	config.AppConfig.AsyncPluginEnabled = false

	router := SetupRouter(nil)

	for _, path := range []string{"/snake", "/game/snake"} {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)
			router.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
			}
			if contentType := recorder.Header().Get("Content-Type"); !strings.Contains(contentType, "text/html") {
				t.Fatalf("expected html content type, got %q", contentType)
			}
			if body := recorder.Body.String(); !strings.Contains(body, "贪吃蛇") || !strings.Contains(body, "createFood") {
				t.Fatalf("snake game html was not rendered")
			}
		})
	}
}
