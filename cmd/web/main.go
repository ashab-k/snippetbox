package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golangcollege/sessions"

	"github.com/ashab-k/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application  struct{
	errLog *log.Logger
	infoLog *log.Logger
	session *sessions.Session
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {


	addr := flag.String("addr" ,":4000" , "HTTP network port") 
	dsn := flag.String("dsn" , "web:anas0707@/snippetbox?parseTime=true" , "MySQL datasource name")
	 
	flag.Parse()
	infoLog := log.New(os.Stdout , "INFO\t" , log.Ldate | log.Ltime)

	errLog := log.New(os.Stderr , "ERROR\t" , log.Ldate|log.Ltime|log.Llongfile)

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
    flag.Parse()

	db , err := openDB(*dsn)
	
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()
	
	templateCache , err := newTemplateCache("./ui/html/")
	if err != nil {
		errLog.Panic(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour


	app := &application{
		errLog: errLog,
		infoLog: infoLog,
		session: session,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache ,
	}
	

  

    infoLog.Printf("Starting server on %s" , *addr)

	srv := &http.Server{ //we create our own instance of a server insead of relying on http.Server
		Addr: *addr, 
		ErrorLog: errLog,
		Handler: app.routes(),
	}
    err = srv.ListenAndServeTLS("./tls/cert.pem" , "./tls/key.pem")
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
