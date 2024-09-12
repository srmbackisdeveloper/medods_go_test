package tests

import (
	"fmt"
	"server/internal/server"
)

func main() {
	s := server.NewServer()
	err := s.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
