package server

import (
	"fmt"
	"time"
)

//func (t *T) MethodName(argType T1, replyType *T2) error

type Server struct {
	StartTime time.Time
}

func NewServer() *Server {
	return &Server{
		StartTime: time.Now(),
	}
}

func (s *Server) Alive(args struct{}, re *time.Duration) error {
	fmt.Println("alive called")
	now := time.Now()
	*re = now.Sub(s.StartTime)

	return nil
}
