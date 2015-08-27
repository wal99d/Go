/*
 This code will do an HTTP middleware action as illustrated below thanks to Alex Edwards:
	mainHandler => dispatchHandler => serveDisptach
	mainHandler will do the following:
		1.Authorizing the request in order to serve the HTTP Request using JWT.
		2.Checking the Content-Length if empty request will discarded
		3.Checking the Content-Type if not application/json => which is "text/plain"
	dispatchHandler will do nothing as of this time until I made my decision :)
	serveDisptach will respond with the required HTTP Response also it modifies the header by adding 
	my own company name which is "InnovativeTech"
*/

package main

import (
	"bytes"
	"net/http"
	"log"
	"github.com/dgrijalva/jwt-go"
)

var serectkey="waleed"

func mainHandler(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
	
		token,err:=jwt.ParseFromRequest(r,func(token *jwt.Token)(interface{}, error){
			return []byte(serectkey),nil
		})
		if err ==nil && token.Valid {
			log.Println("Executing mainHandler")
			log.Println("Content-Length:",r.ContentLength)
			if r.ContentLength == 0 {
				w.Header().Set("xxx","xxx")
				log.Println(http.StatusText(400))
				http.Error(w, http.StatusText(400),400)
				return
			}
			buf :=new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			d:=buf.Bytes()
			log.Println("Content-Type:",http.DetectContentType(d))
			//log.Println(string(d[:33]))
			if http.DetectContentType(buf.Bytes()) != "text/plain; charset=utf-8"{
				w.Header().Set("xxx","xxx")
				http.Error(w, http.StatusText(415), 415)
				return
			}
			h.ServeHTTP(w,r)
			log.Println("Executing mainHandler again")
		} else {
			log.Println("Invalid Token!")
			return
		}
	})
}
func dispatchHandler(h http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		log.Println("Exceuting disptachHandler")
		if r.URL.Path !="/"{
			return
		}
		h.ServeHTTP(w,r)
		log.Println("Exceuting disptachHandler again")
	})
}
func serveDisptach(w http.ResponseWriter, r *http.Request){
	log.Println("Executing serveDisptach")
	w.Header().Set("Developed-By","InnovativeTech")
	w.Write([]byte("Served!!"))
}

func main(){
	token:=jwt.New(jwt.GetSigningMethod("HS256"))
	tokenString,_:=token.SignedString([]byte(serectkey))
	log.Println(tokenString)
	sd :=http.HandlerFunc(serveDisptach)
	
	http.Handle("/", mainHandler(dispatchHandler(sd)))
	http.ListenAndServe(":3000", nil)
}
