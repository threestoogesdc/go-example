// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
  "html/template"
  "io/ioutil"
  "net/http"
  "strings"
  "regexp"
  "log"
)

type Page struct {
  Title string
  Body  []byte
}

const (
  VIEW_PATH = "view/"
  DATA_PATH = "data/"
)

func (p *Page) save() error {
  filename := DATA_PATH + p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
  filename := DATA_PATH + title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title: title, Body: body}, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
  list, err := ioutil.ReadDir(DATA_PATH)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }


  _len := len(list)
  var _pages []string
  _pages = make([]string, _len)

  for i := 0; i < _len; i++ {
    log.Println(list[i].Name())
    _pages[i] = strings.Replace(list[i].Name(), ".txt", "", -1)
  }

  error := templates.ExecuteTemplate(w, "home.html", _pages)
  if error != nil {
    http.Error(w, error.Error(), http.StatusInternalServerError)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    http.Redirect(w, r, "/edit/"+title, http.StatusFound)
    return
  }
  renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title: title}
  }
  renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
  body := r.FormValue("body")
  p := &Page{Title: title, Body: []byte(body)}
  err := p.save()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  http.Redirect(w, r, "/view/"+title, http.StatusFound)
}


var templates = template.Must(template.ParseFiles(VIEW_PATH + "home.html", VIEW_PATH + "edit.html", VIEW_PATH + "view.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
  err := templates.ExecuteTemplate(w, tmpl + ".html", p)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

const lenPath = len("/view/")

var titleValidator = regexp.MustCompile("^[a-zA-Z0-9]+$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    title := r.URL.Path[lenPath:]
    if !titleValidator.MatchString(title) {
      http.NotFound(w, r)
      return
    }
    fn(w, r, title)
  }
}

func main() {
  http.HandleFunc("/", rootHandler)
  http.HandleFunc("/view/", makeHandler(viewHandler))
  http.HandleFunc("/edit/", makeHandler(editHandler))
  http.HandleFunc("/save/", makeHandler(saveHandler))
  http.ListenAndServe(":8080", nil)
}