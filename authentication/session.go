package authentication

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func SetSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := CookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/admin/",
		}
		http.SetCookie(response, cookie)
	}
}

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func GetUserNameString(value string) string {
	cookieValue := make(map[string]string)
	if err := CookieHandler.Decode("session", value, &cookieValue); err == nil {
		return cookieValue["name"]
	}

	return ""
}

func ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/admin/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
