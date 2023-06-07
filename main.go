package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type zer struct {
	Id       int
	Username string
	Password string
	Mail     string
}

type post struct {
	Id       int
	Username string
	Title    string
	Post     string
	Topic    string
	Comments []post
	Likes    int
	Dislikes int
}

var id int
var username string
var mail string
var password string
var zerr zer
var User zer
var Users []zer
var BackupUsers []zer
var Post post
var Posts []post
var Filter []post
var backupPost = Posts
var backupFilter = Filter

func filter(topicFilter string) {
	Filter = backupFilter
	if topicFilter != "Choisissez un sujet" {
		for i := range Posts {
			if Posts[i].Topic == topicFilter {
				Filter = append(Filter, Posts[i])
			}
		}
	}
}

func userTest(user string) bool {
	for place := range Users {
		if Users[place].Username == user {
			return false
		}
	}
	return true
}

func addressTest(mail string) bool {
	for place := range Users {
		if Users[place].Mail == mail {
			return false
		}
	}
	return true
}

func main() {
	Post.Username = "Coranthin"
	Post.Title = "J'ai raté mes pates"
	Post.Post = "J'ai mis le feu à ma cuisine, quelqu'un connais un bon resto"
	Post.Topic = "Cuisine"
	Post.Comments = backupPost
	Post.Likes = 0
	Post.Dislikes = 0
	Post.Id = len(Posts) + 1

	Posts = append(Posts, Post)
	fmt.Println(Posts[0])

	Post.Title = "Le meilleur film du monde"
	Post.Username = "Daniel"
	Post.Post = "Regardez Batman vs ET, c'est de loin le meilleur film de cette decenie"
	Post.Topic = "Cinema"
	Post.Comments = backupPost
	Post.Likes = 0
	Post.Dislikes = 0
	Post.Id = len(Posts) + 1

	Posts = append(Posts, Post)
	fmt.Println(Posts[1])

	Post.Title = "Aidez moi"
	Post.Username = "Ryan"
	Post.Post = "J'ai fini Elden Ring"
	Post.Topic = "Jeux vidéos"
	Post.Comments = backupPost
	Post.Likes = 0
	Post.Dislikes = 0
	Post.Id = len(Posts) + 1

	Posts = append(Posts, Post)
	fmt.Println(Posts[2])

	//connecter = false
	database, _ := sql.Open("sqlite3", "./nraboy.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, username TEXT, mail TEXT, password TEXT)")
	statement.Exec()
	tmpl := template.Must(template.ParseGlob("html/*"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		mdp := r.FormValue("password")
		user := r.FormValue("username")
		adresse := r.FormValue("mail")

		bim, _ := database.Query("SELECT id, username, mail, password FROM people")
		for bim.Next() {
			bim.Scan(&id, &username, &mail, &password)
			User.Id = id
			User.Username = username
			User.Mail = mail
			User.Password = password
			Users = append(Users, User)
		}
		if mdp != "" && user != "" && adresse != "" {
			for x := range Users {
				if mdp == Users[x].Password && user == Users[x].Username && adresse == Users[x].Mail {
					zerr.Username = user
					http.Redirect(w, r, "http://localhost:5550/post", http.StatusSeeOther)
					fmt.Println("Bienvenue")
				} else {
					fmt.Printf("Casse toi d'ici")
				}
			}
		}

		tmpl.ExecuteTemplate(w, "login.html", zerr)
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, z *http.Request) {
		topicFilter := z.FormValue("topic")
		if topicFilter != "Choisissez un sujet" {
			filter(topicFilter)
		}
		switch z.Method {
		case "POST":
			http.Redirect(w, z, "http://localhost:5550/filter", http.StatusSeeOther)
		}
		tmpl.ExecuteTemplate(w, "vitrine.html", Posts)
	})

	http.HandleFunc("/filter", func(w http.ResponseWriter, z *http.Request) {
		tmpl.ExecuteTemplate(w, "vitrine.html", Filter)
		topicFilter := z.FormValue("topic")
		if topicFilter != "Choisissez un sujet" {
			filter(topicFilter)
		}
	})

	http.HandleFunc("/new_post", func(w http.ResponseWriter, z *http.Request) {
		newTopic := z.FormValue("topic")
		newTitle := z.FormValue("title")
		newPost := z.FormValue("post")
		fmt.Println(newTopic, newTitle, newPost)
		if newPost != "" && newTitle != "" && newTopic != "Choisissez un sujet" {
			Post.Username = zerr.Username
			Post.Title = newTitle
			Post.Post = newPost
			Post.Topic = newTopic
			Post.Comments = backupPost
			Post.Likes = 0
			Post.Dislikes = 0
			Post.Id = len(Posts) + 1

			Posts = append(Posts, Post)
			idLink := strconv.Itoa(Post.Id)
			http.Redirect(w, z, "http://localhost:5550/post/"+idLink, http.StatusSeeOther)
		}
		tmpl.ExecuteTemplate(w, "topic.html", Post)
	})

	http.HandleFunc("/post/", func(w http.ResponseWriter, z *http.Request) {

		url := z.URL.RequestURI()
		var num string
		for _, rn := range url {
			if rn >= '0' && rn <= '6' {
				num += string(rn)
			}
		}
		fmt.Println(num)
		number, _ := strconv.Atoi(num)
		fmt.Println(number)

		newComment := z.FormValue("comment")
		if newComment != "" {
			Post.Username = zerr.Username
			Post.Title = Posts[number-1].Title
			Post.Post = newComment
			Post.Topic = Posts[number-1].Topic
			Post.Comments = backupPost
			Post.Likes = 0
			Post.Dislikes = 0
			Post.Id = len(Posts) + 1

			Posts[number-1].Comments = append(Posts[number-1].Comments, Post)
		}
		newComment = ""

		tmpl.ExecuteTemplate(w, "page.html", Posts[number-1])
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, z *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", zerr)
		mdp1 := z.FormValue("newpassword")
		adresse1 := z.FormValue("newmail")
		user1 := z.FormValue("newusername")
		if mdp1 != "" && user1 != "" && adresse1 != "" {
			userTest := userTest(user1)
			addressTest := addressTest(adresse1)
			if addressTest && userTest {
				statement, _ = database.Prepare("INSERT INTO people (username, mail, password) VALUES (?, ?, ?)")
				statement.Exec(user1, adresse1, mdp1)
				rows, _ := database.Query("SELECT id, username, mail, password FROM people")
				tmpl.ExecuteTemplate(w, "pop2.html", zerr)
				for rows.Next() {
					rows.Scan(&id, &username, &mail, &password)
					fmt.Println(strconv.Itoa(id) + ": " + username + " " + mail + " " + password)
				}
			} else {
				tmpl.ExecuteTemplate(w, "pop3.html", zerr)
			}
		}
	})

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe("localhost:5550", nil)
}
