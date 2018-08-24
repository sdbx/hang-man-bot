package imgserv

import (
	"net/http"
	"strconv"

	"github.com/sdbx/hang-man-bot/display"

	"github.com/go-chi/chi"
)

func Start() {
	go func() {
		r := chi.NewRouter()
		r.Get("/img/{maxhp}/{id}/{hp}", getImage)
		if err := http.ListenAndServe(":8053", r); err != nil {
			panic(err)
		}
	}()
}

func getImage(w http.ResponseWriter, r *http.Request) {
	maxhp_ := chi.URLParam(r, "maxhp")
	id_ := chi.URLParam(r, "id")
	hp_ := chi.URLParam(r, "hp")

	maxhp, err := strconv.Atoi(maxhp_)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	id, err := strconv.Atoi(id_)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	hp, err := strconv.Atoi(hp_)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	p := display.GetPicture(maxhp, id, hp)
	http.ServeFile(w, r, p)
}
