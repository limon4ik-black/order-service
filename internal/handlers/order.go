package handlers

import (
	"fmt"
	"net/http"
)

type Order struct {
	n int
}

func (h *Order) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.n = 10
	fmt.Fprintf(w, "count is %d\n", h.n)
	// чтение из кэша
	// чтение из бд
	// маршал ддсж и отправка
}
