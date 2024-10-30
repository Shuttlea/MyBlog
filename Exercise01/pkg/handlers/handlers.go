package handlers

import (
	"fmt"
	"log"
	"myblog/db"
	"myblog/structs"
	"net/http"
	"slices"
	"strconv"
	// "strings"

	"text/template"
	"time"

	"bytes"

	// "github.com/golang-jwt/jwt"
	// "github.com/golang-jwt/jwt/v4"
	"github.com/yuin/goldmark"
	"gorm.io/gorm"

	// "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

// var mySignKey = []byte("very sicret key")

type DBHandler struct {
	DB *gorm.DB
}

func (h DBHandler) HandlerMain(w http.ResponseWriter, r *http.Request) {
	// log.Println("main page ", r.URL)
	page := 1
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		page, _ = strconv.Atoi(pageStr)
	}
	// tmpl, err := template.ParseFiles("../templates/template1.html")
	// if err != nil {
	// 	log.Fatalln("error tempalate.ParseFiles: ", err)
	// }
	d := db.ReadDB(page, h.DB)
	slices.Reverse(d.Nodes)
	d.Previous = fmt.Sprintf("/?page=%d", page-1)
	d.Next = fmt.Sprintf("/?page=%d", page+1)
	tmpl, err := template.ParseFiles("../templates/template1.html")
	if err != nil {
		log.Fatalln("error tempalate.ParseFiles: ", err)
	}
	tmpl.Execute(w, d)
}

func (h DBHandler) HandlerArticle(w http.ResponseWriter, r *http.Request) {
	// log.Println("article page ", r.URL)
	strID := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(strID)
	data := db.ReadArticle(id, h.DB)
	tmpl, err := template.ParseFiles("../templates/article.html")
	if err != nil {
		log.Fatalln("error tempalate.ParseFiles: ", err)
	}
	tmpl.Execute(w, data)
}

func (h DBHandler) HandlerCreate(w http.ResponseWriter, r *http.Request) {
	// log.Println("create page", r.URL)
	// cookie,_:=r.Cookie("token")
	// token := cookie.Value
	// var claims =jwt.MapClaims{}
	// _,err := jwt.ParseWithClaims(token,claims, func(token *jwt.Token)(interface{},error){return []byte(mySignKey),nil})
	// if err!= nil {
	// 	log.Println("Can't parse JWT token:",err)
	// }
	// log.Println(claims.VerifyExpiresAt(time.Now().Unix(),true))
	// log.Printf("Exp from token claims: %s; Type: %T",claims["exp"],claims["exp"])

	// log.Println("Time compare:",time.Now().Unix() - claims["exp"].int64())
	// token.Parse
	// log.Println("Token from cookies: ",token)
	switch r.Method {
	// SHOW PAGE FOR CREATE NEW ARTICLE
	case "GET":
		// if isTokenValid(r) {
		// 	log.Println("Token is valid")
		// 	tmpl, err := template.ParseFiles("../templates/create.html")
		// 	if err != nil {
		// 		log.Fatalln("error tempalate.ParseFiles: ", err)
		// 	}
		// 	tmpl.Execute(w, nil)
		// } else {
		// 	log.Println("Invalid token")
		HandlerAuth(w, r)
		// }
	// CREATE NEW ARTICLE
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Println("ParseForm() err: ", err)
			h.HandlerMain(w, r)
			return
		}
		// USE PACKAGE GOLDMARK TO CONVERT MARKDOWN TO HTML
		markdown := goldmark.New(
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		)
		var buf bytes.Buffer
		err := markdown.Convert([]byte(r.PostFormValue("Article")), &buf)
		if err != nil {
			log.Fatal(err)
		}
		// CREATE NEW ARTICLE STRUCT AND WRITE IT TO DATABASE
		a := structs.Article{Data: time.Now().Format(time.DateTime), Node: buf.String()}
		db.WriteDB(a, h.DB)
		// REDIRECT TO MAIN PAGE
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}

func HandlerFavicon(w http.ResponseWriter, r *http.Request) {
	// EMPTY HANDLER FOR BROWSER REQUEST FOR STYLE
}

func HandlerAuth(w http.ResponseWriter, r *http.Request) {
	// log.Println("Auth page")
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "../templates/auth.html")
	case "POST":
		if isValidLogin(r.PostFormValue("Username"), r.PostFormValue("Password")) {
			// token, _ := generateJWT()
			// http.SetCookie(w, &http.Cookie{
			// 	Name:  "token",
			// 	Value: token,
			// })
			log.Println("User Authorized")
			tmpl, err := template.ParseFiles("../templates/create.html")
			if err != nil {
				log.Fatalln("error tempalate.ParseFiles: ", err)
			}
			tmpl.Execute(w, nil)
			// r.Method = "GET"
			// http.Redirect(w,r,"http://localhost:8888/create",200)
			// w.WriteHeader(http.StatusUnauthorized)
		} else {
			log.Println("Username or Password is wrong")
			r.Method = "GET"
			// w.WriteHeader(http.StatusUnauthorized)
			HandlerAuth(w, r)
		}
	}
}

// func isTokenValid(r *http.Request) bool {
// 	cookie, _ := r.Cookie("token")
// 	token := cookie.Value
// 	var claims = jwt.MapClaims{}
// 	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) { return []byte(mySignKey), nil })
// 	if err != nil {
// 		log.Println("Can't parse JWT token:", err)
// 	}
// 	return claims.VerifyExpiresAt(time.Now().Unix(), true)
// }

func isValidLogin(user, pass string) bool {
	storeUser, storePass := db.LoadCfg(0)
	// fmt.Println(storeUser," ",storePass)
	if storeUser == user && storePass == pass {
		return true
	}
	return false
}

// func generateJWT() (string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["exp"] = time.Now().Add(time.Minute).Unix()

// 	tokenString, err := token.SignedString(mySignKey)
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }
