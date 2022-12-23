package graph

import (
	"github.com/sergiovenicio/graphql-courses/repositories/categories"
	"github.com/sergiovenicio/graphql-courses/repositories/courses"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
type Resolver struct {
	CategoriesRepository *categories.CategoriesRepository
	CoursesRepository    *courses.CoursesRepository
}
