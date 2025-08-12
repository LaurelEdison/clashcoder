package handlers

import (
	"net/http"
)

func (h *Handlers) HandlerHealth(w http.ResponseWriter,
	r *http.Request) {

	h.RespondWithJSON(w, 200, struct{}{})

}
