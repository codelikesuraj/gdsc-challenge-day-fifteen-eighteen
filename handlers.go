package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type BookHandler struct {
	DB *gorm.DB
}

func (b *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/index.html"))
	var books []Book
	b.DB.Order("id DESC").Find(&books)
	tmpl.Execute(w, map[string][]Book{"books": books})
}

func (b *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/create.html"))
	tmpl.Execute(w, nil)
}

func (b *BookHandler) StoreBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "error parsing form: %s", err.Error())
		return
	}

	txn := b.DB.Create(&Book{
		Author: r.FormValue("author"),
		Title:  r.FormValue("title"),
	})
	if err := txn.Error; err != nil {
		fmt.Fprintf(w, "error saving form: %s", err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (b *BookHandler) EditBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimPrefix(strings.Trim(r.URL.Path, "/"), "edit/"))
	var book Book
	res := b.DB.First(&book, id)
	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, "error fetching book: %s", err.Error())
		}
		return
	}
	tmpl := template.Must(template.ParseFiles("template/edit.html"))
	tmpl.Execute(w, map[string]Book{"book": book})
}

func (b *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "error parsing form: %s", err.Error())
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	var book Book
	res := b.DB.First(&book, id)
	if err := res.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.NotFound(w, r)
		} else {
			fmt.Fprintf(w, "error fetching book: %s", err.Error())
		}
		return
	}

	txn := b.DB.Model(&book).Updates(map[string]interface{}{
		"author": r.FormValue("author"),
		"title":  r.FormValue("title"),
	})
	if err := txn.Error; err != nil {
		fmt.Fprintf(w, "error udpating form: %s", err.Error())
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/edit/%d", book.ID), http.StatusMovedPermanently)
}

func (b *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "error parsing form: %s", err.Error())
		return
	}

	id, _ := strconv.Atoi(r.FormValue("id"))

	res := b.DB.Delete(&Book{}, id)
	if err := res.Error; err != nil {
		fmt.Fprintf(w, "error deleting book: %s", err.Error())
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
