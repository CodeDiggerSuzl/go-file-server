package handler

import (
	"fmt"
	"go-file-server/db"
	"go-file-server/utils"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	// User password encryption salt
	passwordSalt = "*#890"
	tokenSalt    = "_tokensalt"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	// return register page
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
		return
	}
	// else is post and it will do the login stuff
	_ = r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	// simple check the form name
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

// UserSignInHandler first check user info
// if passed the return a certificate: either by cookie and session or by token(user will carry the token every time)
func UserSignInHandler(w http.ResponseWriter, r *http.Request) {
	getwd, _ := os.Getwd()
	fmt.Println(getwd)
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	_ = r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	passwdToCheck := utils.Sha1([]byte(password + passwordSalt))
	// check user info
	ok := db.UserSignin(username, passwdToCheck)
	if !ok {
		_, _ = w.Write([]byte("FAILED"))
		return
	}
	// if passed return token
	token := GenerateToken(username)
	// put in into redisOps
	_, _ = db.SetUserTokenIntoRedis(username, token)
	tokenOK := db.UpdateUserToken(username, token)
	if !tokenOK {
		_, _ = w.Write([]byte("FAILED"))
		return
	}
	// return login user info or redirect to home page
	// temp return a temp url
	// http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
	resp := utils.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			UserName string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			UserName: username,
			Token:    token,
		},
	}
	_, _ = w.Write(resp.JSONBytes())
	// w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

// GenerateToken get a random token
func GenerateToken(username string) string {
	// 40 bit len tokenPrefix(use_name+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := utils.MD5([]byte(username + ts + tokenSalt))
	return tokenPrefix + ts[:8]
}

// UserInfoHandler query user info
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")
	// validate user token
	if !validUserToken(username, token) {
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

}
func validUserToken(username, token string) bool {
	// valid length
	if len(token) != 40 {
		return false
	}
	// dbToken, err := db.GetUserToken(username)
	// if err != nil {
	// 	fmt.Printf("err getting user token:%s", err.Error())
	// 	return false
	// }

	// valid user token is expired or not should put in redisOps
	redisToken, err := db.GetUserTokenFromRedis(username)
	// get nothing
	if err != nil && redisToken == "" {
		return false
	}
	// valid is the right token or not,query from db
	if redisToken != token {
		return false
	}
	// if dbToken != token {
	// 	return false
	// }
	return true
}
