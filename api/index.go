package handler

import (
	"encoding/json"
	"net/http"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pocketbase/dbx"
)



var dbconn *dbx.DB
var once *sync.Once
func Handler(w http.ResponseWriter, r * http.Request) {
	var err error
	once.Do(func() {
		dbconn, err = dbx.Open("sqlite3", "file:test.db")
	})

	if err != nil || dbconn == nil {
		http.Error(w, "error when connecting db", http.StatusInternalServerError)
	}

	NewServer(dbconn).
		ServeHTTP(w, r)
}


type Server struct {
	dal *dbx.DB
}

func NewServer(dal *dbx.DB) *Server {
	return &Server{
		dal: dal,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.dal = s.dal.WithContext(r.Context())
	router := chi.NewRouter()
	
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(router chi.Router) {
		router.Route("/users", func(router chi.Router){
			router.Post("/login", s.HandleUserLogin)
			router.Post("/register", s.HandleUserRegister)
			router.Get("/me", s.HandleUserMe)
			router.Route("/{{user_id}}", func(router chi.Router) {
				router.Use(s.MiddlewareUserCtx)
				router.Route("/requests", func(router chi.Router) {
					router.Get("/", s.HandleUserRequestIndex)
					router.Post("/", s.HandleUserRequestStore)
					router.Route("/{{request_id}}", func(router chi.Router) {
						router.Use(s.MiddlewareUserRequestCtx)
						router.Delete("/", s.HandleUserRequestDestroy)
					})
				})
			})
		})

		router.Route("/posts", func(router chi.Router){
			router.Get("/", s.HandlePostIndex)
			router.Post("/", s.HandlePostStore)
			
			router.Route("/{{post_id}}", func(router chi.Router) {
				router.Use(s.MiddlewarePostCtx)
				router.Put("/", s.HandlePostUpdate)
				router.Delete("/", s.HandlePostDestroy)
				
				router.Route("/comments", func(router chi.Router){
					router.Get("/", s.HandlePostCommentIndex)
					router.Post("/", s.HandlePostCommentStore)
					
					router.Route("/{{comment_id}}", func(router chi.Router){
						router.Use(s.MiddlewarePostCommentCtx)
						router.Put("/", s.HandlePostCommentUpdate)
						router.Delete("/", s.HandlePostCommentDestroy)
					})
				})
			})
		})

		router.Route("/friends", func(router chi.Router){
			router.Get("/", s.HandleFriendIndex)
			router.Route("/{{friend_id}}", func(router chi.Router){
				router.Use(s.MiddlewareFriendCtx)
				router.Delete("/", s.HandleFriendDestroy)
			})
		})
	})

	router.ServeHTTP(w, r)
}

func (s *Server) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) HandleUserRegister(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleUserMe(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) MiddlewareUserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HandleUserRequestIndex(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) HandleUserRequestStore(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) MiddlewareUserRequestCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
	})
}

func (s *Server) HandleUserRequestDestroy(w http.ResponseWriter, r *http.Request) {

}


func (s *Server) HandlePostStore(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandlePostIndex(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) MiddlewarePostCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HandlePostUpdate(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandlePostDestroy(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandlePostCommentIndex(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandlePostCommentStore(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) MiddlewarePostCommentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) HandlePostCommentUpdate(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandlePostCommentDestroy(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandleFriendIndex(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) HandleFriendDestroy(w http.ResponseWriter, r *http.Request) {
	
}

func (s *Server) MiddlewareFriendCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}


func toJSON(w http.ResponseWriter, data interface{}, httpStatusCode... int) {
	statusCode := http.StatusOK

	if len(httpStatusCode) > 0 {
		statusCode = httpStatusCode[0]
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"code": 500, "message": "error while encoding json"}`, http.StatusInternalServerError)
	}
}