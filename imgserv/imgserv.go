package imgserv

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/sdbx/hang-man-bot/display"

	"github.com/go-chi/chi"
)

func Start() {
	go func() {
		r := chi.NewRouter()
		r.Get("/man/{maxhp}/{id}/{hp}", getMan)
		if err := http.ListenAndServe(":8053", r); err != nil {
			panic(err)
		}
	}()
}

func getMan(w http.ResponseWriter, r *http.Request) {
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

	var hp int
	fmt.Sscanf(hp_, "%d", &hp)

	p := display.GetManImage(maxhp, id, hp)
	http.ServeFile(w, r, p)
}
