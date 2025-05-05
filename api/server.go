package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Api struct {
	http.Server
}

func InitServer() *Api {
	return &Api{
		http.Server{
			Addr:         os.Getenv("IP_ADDR"),
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      routes(),
		},
	}
}

func (a *Api) Run() {

	go func() {
		fmt.Printf("Listening on %v\n", a.Addr)
		if err := a.ListenAndServe(); err != nil {
			fmt.Printf("Stoped Listening: %v\n", err)
		}
	}()

	shutdown, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	<-shutdown.Done()

	fmt.Println("Shutting down server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Shutdown(ctx); err != nil {
		fmt.Printf("Shutting Down with error: %v", err)
	}

	fmt.Println("Shutdown Complete")
	fmt.Println("Waiting for Background Jobs to finish ....")
	wg.Wait()
	fmt.Println("Background Jobs ended")
}
