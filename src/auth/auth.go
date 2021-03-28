package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/NgeKaworu/to-do-list-go/src/resultor"
	"github.com/julienschmidt/httprouter"
)

// Auth 加解密相关
type Auth struct {
	UCHost *string
}

// NewAuth 工厂方法
func NewAuth(usHost *string) *Auth {
	return &Auth{
		UCHost: usHost,
	}
}

// IsLogin isLogin middleware
func (a *Auth) IsLogin(next httprouter.Handle) httprouter.Handle {
	//权限验证
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, errors.New("token is empty"))
			return
		}
		url := *a.UCHost + "/isLogin"
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, errors.New("invalidate request"))
			return
		}

		req.Header.Set("Authorization", auth)

		client := &http.Client{}

		res, err := client.Do(req)

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, errors.New("invalidate uc reqest"))
			return
		}

		body, err := ioutil.ReadAll(res.Body)
		if body != nil {
			defer r.Body.Close()
		}

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, errors.New("invalidate uc response"))
			return
		}

		p := make(map[string]interface{})
		err = json.Unmarshal(body, &p)

		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, errors.New("json params fail"))
			return
		}

		if p["ok"] == true {
			r.Header.Set("uid", p["data"].(string))
			next(w, r, ps)
		} else {
			w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
			w.WriteHeader(http.StatusUnauthorized)
			resultor.RetFail(w, errors.New(p["errMsg"].(string)))
			return
		}

	}
}
