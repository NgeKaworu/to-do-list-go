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

	"github.com/NgeKaworu/to-do-list-go/src/app"
	"github.com/NgeKaworu/to-do-list-go/src/cors"
	"github.com/NgeKaworu/to-do-list-go/src/db"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8040", "绑定Host地址")
		dbinit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "to-do-list", "database name")
		ucHost = flag.String("uc", "http://user-center-go", "user center host")
		r      = flag.String("r", "localhost:6379", "rdb addr")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	mongoClient := db.NewMongoClient()
	err := mongoClient.Open(*mongo, *mdb, *dbinit)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     *r,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	app := app.New(ucHost, mongoClient, rdb)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	//task ctrl
	router.POST("/v1/task/create", app.AddTask)
	router.PATCH("/v1/task/update", app.SetTask)
	router.GET("/v1/task/list", app.ListTask)
	router.DELETE("/v1/task/:id", app.RemoveTask)

	srv := &http.Server{Handler: app.IsLogin(cors.CORS(router)), ErrorLog: nil}
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
			mongoClient.Close()
			rdb.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
