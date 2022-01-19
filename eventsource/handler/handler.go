package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render

var currentID int
var userMap map[int]*User

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(&user)
	// if check.CheckError(err, w, http.StatusBadRequest) {
	// 	return
	// }
	if err != nil {
		rd.Text(w, http.StatusBadRequest, err.Error())
	}

	currentID++
	user.ID = currentID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user

	// w.Header().Add("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)

	// data, _ := json.Marshal(&user)
	// fmt.Fprint(w, string(data))
	rd.JSON(w, http.StatusOK, &user)
}

func handleGetUsers(w http.ResponseWriter, _ *http.Request) {
	if len(userMap) == 0 {
		rd.JSON(w, http.StatusOK, "No user")
		return
	}

	users := []*User{}
	for _, value := range userMap {
		users = append(users, value)
	}
	// w.Header().Add("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	// data, _ := json.Marshal(&users)
	// fmt.Fprint(w, string(data))
	rd.JSON(w, http.StatusOK, users)
}

func initData() {
	rd = render.New()
	currentID = 0
	userMap = make(map[int]*User)
}

func NewHttpHandler() http.Handler {
	initData()

	es := eventsource.New(nil, nil)
	defer es.Close()

	mux := pat.New()
	mux.Handle("/stream", es) // es 오픈될때 매핑
	mux.Post("/user", handleCreateUser)
	mux.Get("/users", handleGetUsers)

	n := negroni.Classic()
	n.UseHandler(mux)

	return n
}
