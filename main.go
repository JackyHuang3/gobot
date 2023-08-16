package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gobot/global"
	"gobot/webhook"

	"github.com/julienschmidt/httprouter"
)

var (
	port = flag.String("port", "8080", "--port=8080")
)

func main() {
	flag.Parse()

	// init
	webhook.LogPrintf = global.LogPrintf
	webhook.LogPrintln = global.LogPrintln
	config, err := global.InitConfig("")
	FatalCheck(err)
	store, err := global.InitStore(config)
	FatalCheck(err)
	webhook.Init(config, store)

	// router
	handler := initRouter(webhook.RouterInit)

	// server
	log.Printf("start server for web at http://0.0.0.0:%s\n", *port)
	server := &http.Server{Addr: fmt.Sprintf(":%v", *port), Handler: handler}
	if err := server.ListenAndServe(); err != nil {
		FatalCheck(fmt.Errorf("start webServer failed, detail: %s", err.Error()))
	}
}

func initRouter(routerFuncList ...func(*httprouter.Router)) *httprouter.Router {
	router := httprouter.New()

	for _, routerFunc := range routerFuncList {
		if routerFunc != nil {
			routerFunc(router)
		}
	}
	return router
}

func FatalCheck(err error) {
	if err != nil {
		log.Println(err.Error() + "\n")
		os.Exit(0)
	}
}
