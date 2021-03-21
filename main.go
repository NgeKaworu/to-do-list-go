package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/NgeKaworu/to-do-list-go/src/auth"
	"github.com/NgeKaworu/to-do-list-go/src/cors"
	"github.com/NgeKaworu/to-do-list-go/src/engine"
	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8031", "绑定Host地址")
		dbinit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		db     = flag.String("db", "time-mgt", "database name")
		ucHost = flag.String("uc", "http://localhost:8011", "user center host")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	auth := auth.NewAuth(ucHost)
	eng := engine.NewDbEngine()
	err := eng.Open(*mongo, *db, *dbinit)

	if err != nil {
		log.Println(err.Error())
	}

	router := httprouter.New()
	// tag ctrl
	router.POST("/v1/tag/create", auth.IsLogin(eng.AddTag))
	router.PUT("/v1/tag/update", auth.IsLogin(eng.SetTag))
	router.GET("/v1/tag/list", auth.IsLogin(eng.ListTag))
	router.DELETE("/v1/tag/:id", auth.IsLogin(eng.RemoveTag))
	//record ctrl
	router.POST("/v1/record/create", auth.IsLogin(eng.AddRecord))
	router.PUT("/v1/record/update", auth.IsLogin(eng.SetRecord))
	router.GET("/v1/record/list", auth.IsLogin(eng.ListRecord))
	router.DELETE("/v1/record/:id", auth.IsLogin(eng.RemoveRecord))
	router.POST("/v1/record/statistic", auth.IsLogin(eng.StatisticRecord))

	srv := &http.Server{Handler: cors.CORS(router), ErrorLog: nil}
	srv.Addr = *addr

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server on http port", srv.Addr)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			go func() {
				_ = srv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			eng.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
