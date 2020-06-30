package handlers

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
)

func health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return web.NewRequestError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return web.Respond(ctx, w, status, http.StatusOK)
}
