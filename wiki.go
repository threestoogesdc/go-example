package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
)

type Page struct {
  Title string
  Body []byte
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
      return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  t, _ := template.ParseFiles("views/" + tmpl + ".html")
  t.Execute(w, p)
}

const lenPath = len("/view/")

func viewHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[lenPath:]
  p, error := loadPage(title)
  if error != nil {
    http.Redirect(w, r, "/edit/" + title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[lenPath:]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }

  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[lenPath:]
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  p.save()
  http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func main() {
  http.HandleFunc("/view/", viewHandler)
  http.HandleFunc("/edit/", editHandler)
  http.HandleFunc("/save/", saveHandler)
  http.ListenAndServe(":8080", nil)
}
