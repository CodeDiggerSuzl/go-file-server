package handler

import (
	"go-file-server/db"
	"go-file-server/utils"
	"io/ioutil"
	"net/http"
)

const passwordSalt = "*#890"

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// return register page
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	// else is post and it will do the login stuff
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		_, _ = w.Write([]byte("invalid parameter"))
		return
	}
	// encrypt the password
	encryptedPassword := utils.Sha1([]byte(password + passwordSalt))
	ok := db.UserSign(username, encryptedPassword)
	if ok {
		_, _ = w.Write([]byte("SUCCESS"))
	} else {
		_, _ = w.Write([]byte("FAILED"))
	}
}
