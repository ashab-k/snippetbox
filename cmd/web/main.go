package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/ashab-k/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application  struct{
	errLog *log.Logger
	infoLog *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	addr := flag.String("addr" ,":4000" , "HTTP network port")
	dsn := flag.String("dsn" , "web:anas0707@/snippetbox?parseTime=true" , "MySQL datasource name")
	 
	flag.Parse()
	infoLog := log.New(os.Stdout , "INFO\t" , log.Ldate | log.Ltime)

	errLog := log.New(os.Stderr , "ERROR\t" , log.Ldate|log.Ltime|log.Llongfile)


	db , err := openDB(*dsn)
	
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()
	app := &application{
		errLog: errLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

  

    infoLog.Printf("Starting server on %s" , *addr)

	srv := &http.Server{ //we create our own instance of a server insead of relying on http.Server
		Addr: *addr, 
		ErrorLog: errLog,
		Handler: app.routes(),
	}
    err = srv.ListenAndServe()
    errLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB , error){
	db , err := sql.Open("mysql" , dsn)
	if err != nil {
		return nil , err
	}

	if err = db.Ping() ; err != nil{
		return nil , err
	}
	return db , nil 
}
