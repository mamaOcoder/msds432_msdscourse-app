package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const PORT = ":8080"

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	Body := "Thanks for visiting!\n"
	fmt.Fprintf(w, "%s", Body)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get CID
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)
	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)

	cid := paramStr[2]
	err := DeleteCourse(cid)
	if err != nil {
		fmt.Println(err)
		Body := err.Error() + "\n"
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", Body)
		return
	}

	Body := cid + " deleted!\n"
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", Body)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	courseList, err := ListCourses()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Body := "Failed to list courses\n"
		fmt.Fprintf(w, "%s", Body)
		return
	}
	var Body string
	for _, c := range courseList {
		Body = Body + c.CID + " " + c.CNAME + " " + c.CPREREQ + "\n"
	}
	fmt.Fprintf(w, "%s", Body)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving:", r.URL.Path, "from", r.Host)
	w.WriteHeader(http.StatusOK)
	courseList, err := ListCourses()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Body := "Failed to get courses\n"
		fmt.Fprintf(w, "%s", Body)
		return
	}
	Body := fmt.Sprintf("Total entries: %d\n", len(courseList))
	fmt.Fprintf(w, "%s", Body)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	// Split URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 5 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not enough arguments: "+r.URL.Path)
		return
	}

	cid := paramStr[2]
	cname := paramStr[3]
	cp := paramStr[4]

	temp := MSDSCourse{CID: cid, CNAME: cname, CPREREQ: cp}
	err := AddCourse(temp)

	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		Body := "Failed to add record\n"
		fmt.Fprintf(w, "%s", Body)
	} else {
		log.Println("Serving:", r.URL.Path, "from", r.Host)
		Body := "New record added successfully\n"
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", Body)
	}

	log.Println("Serving:", r.URL.Path, "from", r.Host)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	// Get Search value from URL
	paramStr := strings.Split(r.URL.Path, "/")
	fmt.Println("Path:", paramStr)

	if len(paramStr) < 3 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "Not found: "+r.URL.Path)
		return
	}

	var Body string
	cid := paramStr[2]
	t, err := SearchCourse(cid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		Body = "Could not be found: " + cid + "\n"
	} else {
		w.WriteHeader(http.StatusOK)
		Body = t.CID + " " + t.CNAME + " " + t.CPREREQ + "\n"
	}

	fmt.Println("Serving:", r.URL.Path, "from", r.Host)
	fmt.Fprintf(w, "%s", Body)
}
