package main

import (
	"crypto/tls"
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
type contextKey string

var contextKeyUser = contextKey("user")
type application  struct{
	errLog *log.Logger
	infoLog *log.Logger
	session *sessions.Session
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
	users *mysql.UserModel
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
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		errLog: errLog,
		infoLog: infoLog,
		session: session,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: templateCache ,
		users : &mysql.UserModel{DB: db},
	}
	
	tlsConfig :=  &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}
  

    infoLog.Printf("Starting server on %s" , *addr)

	srv := &http.Server{
        Addr:         *addr,
        ErrorLog:     errLog,
        Handler:      app.routes(),
        TLSConfig:    tlsConfig,
        // Add Idle, Read and Write timeouts to the server.
        IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
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
