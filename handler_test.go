package knaudit_proxy_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	knaudit_proxy "github.com/navikt/knaudit-proxy"
)

type mockSender struct {
	Err error
}

func (m mockSender) Send(_ string) error {
	return m.Err
}

func (m mockSender) Ping() error {
	return m.Err
}

func (m mockSender) Close() error {
	return m.Err
}

func (m mockSender) Open() error {
	return m.Err
}

func TestHandlers(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		handler    http.HandlerFunc
		expect     string
		expectCode int
	}{
		{
			name:       "ReportHandler",
			handler:    knaudit_proxy.NewServer(&mockSender{}).ReportHandler,
			expect:     `{"status":"ok","message":"audit data stored","code":200}`,
			expectCode: http.StatusOK,
		},
		{
			name:       "ReportHandler with error",
			handler:    knaudit_proxy.NewServer(mockSender{Err: fmt.Errorf("oops")}).ReportHandler,
			expect:     `{"status":"internal server error","message":"opening audit backend","code":500}`,
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			tc.handler.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{}")))

			if w.Code != tc.expectCode {
				t.Errorf("expected status code %d, got %d", tc.expectCode, w.Code)
			}

			if diff := cmp.Diff(strings.TrimSpace(w.Body.String()), strings.TrimSpace(tc.expect)); diff != "" {
				t.Errorf("unexpected response (-want +got):\n%s", diff)
			}
		})
	}
}
