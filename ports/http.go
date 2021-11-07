package ports

import (
	"errors"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
	"net/http"
	"nononsensecode.com/session-management/domain/model"
)

type HttpServer struct {
	repo model.UserRepository
	sessionManager *scs.SessionManager
}

func NewHttpServer(repo model.UserRepository, sessionManager *scs.SessionManager) HttpServer {
	if repo == nil {
		panic("repo is nil")
	}

	if sessionManager == nil {
		panic("session manager is nil")
	}

	return HttpServer{
		repo: repo,
		sessionManager: sessionManager,
	}
}

func (s HttpServer) LoginUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("request came to login")
	var b LoginUserJSONRequestBody
	err := render.Bind(r, &b)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m := Message{Msg: "error occurred while extracting body"}
		render.Render(w, r, &m)
		return
	}

	u, err := s.repo.Find(b.Username, b.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		m := Message{Msg: "unauthorized user"}
		render.Render(w, r, &m)
		return
	}

	user := User{
		Name:     u.Name,
		Username: u.UserName,
	}
	s.sessionManager.Put(r.Context(), "user", user.Username)

	w.WriteHeader(http.StatusOK)
	render.Render(w, r, &user)
}

func (s HttpServer) LogoutUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("request came to logout the user")
	s.sessionManager.Destroy(r.Context())
	w.WriteHeader(http.StatusOK)
}

func (s HttpServer) ListUserDetails(w http.ResponseWriter, r *http.Request) {
	logrus.Info("request came to list user details")
	username := s.sessionManager.GetString(r.Context(), "user")
	if username == "" {
		logrus.Warn("session empty")
		s.LogoutUser(w, r)
		return
	}
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		m := Message{Msg: "there is no such user"}
		render.Render(w, r, &m)
		return
	}

	w.WriteHeader(http.StatusOK)
	user := User{
		Name:     u.Name,
		Username: u.UserName,
	}
	render.Render(w, r, &user)
}

func (b *LoginUserJSONRequestBody) Bind(r *http.Request) error {
	if b == nil {
		return errors.New("request body is nil")
	}

	return nil
}

func (m *Message) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (u *User) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
