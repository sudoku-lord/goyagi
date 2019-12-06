// pkg/server/server_test.go
package server

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
    srv := New()

    t.Run("serves registered endpoint", func(tt *testing.T) {
        w := httptest.NewRecorder()
        req, err := http.NewRequest("GET", "/", nil)
        require.Nil(tt, err, "unexpected error when making new request")

        srv.Handler.ServeHTTP(w, req)

        assert.Equal(tt, http.StatusNotFound, w.Code, "incorrect status code")
        assert.Contains(tt, w.Body.String(), "Not Found", "incorrect response")
    })
}
