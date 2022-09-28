package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	forums "github.com/kuzkuss/VK_DB_Project/internal/app/models/forumsRepository"
	threads "github.com/kuzkuss/VK_DB_Project/internal/app/models/threadsRepository"
	users "github.com/kuzkuss/VK_DB_Project/internal/app/models/usersRepository"
)

type ForumRouter struct {
	*mux.Router
	fr *forums.ForumsRep
}

type ThreadRouter struct {
	*mux.Router
	tr *threads.ThreadsRep
}

type UserRouter struct {
	*mux.Router
	ur *users.UsersRep
}

type MessageError struct {
	Message string
}

func NewForumRouter(fr *forums.ForumsRep) *ForumRouter {
	r := mux.NewRouter()

	frR := &ForumRouter {
		Router: r,
		fr: fr,
	}

	frR.HandleFunc("/forum/create", frR.CreateForum)
	frR.HandleFunc("/forum/{slug}/details", frR.GetForum)
	return frR
}

func NewThreadRouter(tr *threads.ThreadsRep) *ThreadRouter {
	r := mux.NewRouter()

	tdR := &ThreadRouter {
		Router: r,
		tr: tr,
	}

	tdR.HandleFunc("/forum/{slug}/create", tdR.CreateThread)
	return tdR
}

func NewUserRouter(tr *threads.ThreadsRep) *ThreadRouter {
	r := mux.NewRouter()

	tdR := &ThreadRouter {
		Router: r,
		tr: tr,
	}

	tdR.HandleFunc("/forum/{slug}/create", tdR.CreateThread)
	return tdR
}

func (router *ForumRouter) CreateForum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	f := forums.Forum{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	createdForum, err := router.fr.CreateForum(f)
	if err.Error() == "Can't find user with name " + f.User {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{Message: err.Error()})
		http.Error(w, "Can't find user", http.StatusBadRequest)
		return
	} else if err.Error() == "Forum existed" {
		http.Error(w, "Forum existed", http.StatusBadRequest)
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		http.Error(w, "Not created",http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(createdForum)
}

func (router *ForumRouter) GetForum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	forum, err := router.fr.GetForum(slug)
	if err.Error() == "Can't find forum with slug " + slug {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{Message: err.Error()})
		http.Error(w, "Can't find forum", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Not found",http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(forum)
}

func (router *ThreadRouter) CreateThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	t := threads.Thread{}

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	createdThread, err := router.tr.CreateThread(t, slug)
	if err.Error() == "Can't find user with name " + t.Author {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{Message: err.Error()})
		http.Error(w, "Can't find user", http.StatusNotFound)
		return
	} else if err.Error() == "Can't find forum with slug " + slug {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{Message: err.Error()})
		http.Error(w, "Can't find forum", http.StatusNotFound)
		return
	} else if err.Error() == "Thread existed" {
		http.Error(w, "Thread existed", http.StatusConflict)
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		http.Error(w, "Not created",http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(createdThread)
}

func (router *UserRouter) GetUsersFromForum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	slug := vars["slug"]

	users, err := router.ur.SelectUsersFromForum(slug)
	if err.Error() == "Can't find forum with slug " + slug {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{Message: err.Error()})
		http.Error(w, "Can't find forum", http.StatusBadRequest)
		return
	} else if err != nil {
		http.Error(w, "Not found",http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(forum)
}