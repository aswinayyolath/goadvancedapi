package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Course represents the model for courses.
type Course struct {
	CourseId   string  `json:"courseid"`
	CourseName string  `json:"coursename"`
	Price      int     `json:"price"`
	Author     *Author `json:"author"`
}

// Author represents the model for the course author.
type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

var courses []Course // fake DB

// IsEmpty checks if the Course object is empty.
func (c *Course) IsEmpty() bool {
	return c.CourseName == ""
}

func main() {
	fmt.Println("Golang Advanced API")

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome).Methods("GET")
	r.HandleFunc("/courses", getAllCourses).Methods("GET")
	r.HandleFunc("/course/{id}", getSingleCourse).Methods("GET")
	r.HandleFunc("/course", createSingleCourse).Methods("POST")
	r.HandleFunc("/course/{id}", updateOneCourseById).Methods("PUT")
	r.HandleFunc("/course/{id}", deleteOneCourseById).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":4000", r))
}

// serveHome handles the root route.
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Golang Api</h1>"))
}

// getAllCourses returns all courses.
func getAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	returnCourseDetails(w)
}

// returnCourseDetails writes course details to the response.
func returnCourseDetails(w http.ResponseWriter) {
	if isCoursePresent() {
		json.NewEncoder(w).Encode(courses)
	} else {
		json.NewEncoder(w).Encode("No courses found")
	}
}

// isCoursePresent checks if there are any courses in the fake DB.
func isCoursePresent() bool {
	return courses != nil
}

// getSingleCourse retrieves a single course by ID.
func getSingleCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Grab ID from request
	params := mux.Vars(r)

	if isCoursePresent() {
		// Loop through courses and find a matching ID
		for _, course := range courses {
			if course.CourseId == params["id"] {
				json.NewEncoder(w).Encode(course)
				return
			}
		}
		json.NewEncoder(w).Encode("Course with id " + params["id"] + " not found")

	} else {
		json.NewEncoder(w).Encode("No courses found")
	}
}

// createSingleCourse creates a new course.
func createSingleCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		json.NewEncoder(w).Encode("No Request Body")
	}

	var course Course
	json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data present in JSON")
	} else {
		// Generate a unique ID as a string
		rand.New(rand.NewSource(time.Now().UnixNano()))
		course.CourseId = strconv.Itoa(rand.Intn(100))

		// Append the new course to the list of courses
		courses = append(courses, course)
		json.NewEncoder(w).Encode(course)
	}
}

// updateOneCourseById updates a course by ID.
func updateOneCourseById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Grab ID from request
	params := mux.Vars(r)
	coursePresent := false

	for index, course := range courses {
		if course.CourseId == params["id"] {
			coursePresent = true
			// Remove the existing course
			courses = append(courses[:index], courses[index+1:]...)

			// Decode the request body to update the course
			var updatedCourse Course
			json.NewDecoder(r.Body).Decode(&updatedCourse)
			updatedCourse.CourseId = params["id"]

			// Append the updated course to the list of courses
			courses = append(courses, updatedCourse)
			json.NewEncoder(w).Encode(updatedCourse)
		}
	}
	if !coursePresent {
		json.NewEncoder(w).Encode("The Provided Id is not present")
	}
}

// deleteOneCourseById deletes a course by ID.
func deleteOneCourseById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for index, course := range courses {
		if course.CourseId == params["id"] {
			// Remove the course from the list
			courses = append(courses[:index], courses[index+1:]...)
			json.NewEncoder(w).Encode("Record has been deleted")
			return
		}
	}
	// If the course with the specified ID is not found
	json.NewEncoder(w).Encode("Course with id " + params["id"] + " not found")
}
