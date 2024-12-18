package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	addr := flag.String("addr" ,":4000" , "HTTP network port")
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/snippet", showSnippet)
	
    mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/" , http.StripPrefix("/static" , fileServer))


	infoLog := log.New(os.Stdout , "INFO\t" , log.Ldate | log.Ltime)

	errLog := log.New(os.Stderr , "ERROR\t" , log.Ldate|log.Ltime|log.Llongfile)

    infoLog.Printf("Starting server on %s" , *addr)

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errLog,
		Handler: mux,
	}
    err := srv.ListenAndServe()
    errLog.Fatal(err)
}
