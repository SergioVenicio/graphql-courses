package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sergiovenicio/graphql-courses/graph"
	"github.com/sergiovenicio/graphql-courses/graph/generated"
	"github.com/sergiovenicio/graphql-courses/repositories/categories"
	"github.com/sergiovenicio/graphql-courses/repositories/courses"
)

const defaultPort = "8080"

func main() {
	db, err := sql.Open("sqlite3", "./courses.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	categoriesRepository := categories.NewCategoriesRepository(db)
	coursesRepository := courses.NewCoursesRepository(db, categoriesRepository)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		CategoriesRepository: categoriesRepository,
		CoursesRepository:    coursesRepository,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
