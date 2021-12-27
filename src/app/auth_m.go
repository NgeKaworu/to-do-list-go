package app

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/NgeKaworu/to-do-list-go/src/resultor"
)

func (app *App) IsLogin(next http.Handler) http.Handler {
	//权限验证
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := app.checkUser(r)

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, err)
			return
		}

		r.Header.Set("uid", *s)
		next.ServeHTTP(w, r)

	})
}

// checkUser
func (app *App) checkUser(r *http.Request) (*string, error) {
	bear, err := app.getBearer(r)
	if err != nil {
		return nil, err
	}

	s, err := app.rdb.Get(context.Background(), *bear).Result()

	if err != nil {
		return nil, err
	}

	return &s, nil
}

// getBearer
func (app *App) getBearer(r *http.Request) (*string, error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("unknown authorization type")
	}
	auth = auth[7:]
	return &auth, nil
}
