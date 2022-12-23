package courses

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sergiovenicio/graphql-courses/graph/model"
	"github.com/sergiovenicio/graphql-courses/repositories/categories"
	sqlReader "github.com/sergiovenicio/graphql-courses/sql"
)

type CoursesRepository struct {
	db         *sql.DB
	categories *categories.CategoriesRepository
}

func (r *CoursesRepository) Create(input model.NewCourse) (*model.Course, error) {
	log.Printf("[courses_repository] create (%v)...", input)
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("[courses_repository] create: error(%v)!", err)
		return nil, err
	}

	category, err := r.categories.FindById(input.CategoryID)
	if err != nil {
		log.Printf("[courses_repository] create: error(%v)!", err)
		return nil, fmt.Errorf("cant find category (category_id: %s)", input.CategoryID)
	}

	log.Printf("[courses_repository] create: created(%v)!", category)

	course := &model.Course{
		ID:          id.String(),
		Name:        input.Name,
		Description: input.Description,
	}

	query, err := sqlReader.ReadQuery("courses/insert_course")
	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(query, course.ID, course.Name, course.Description, category.ID)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (r *CoursesRepository) Find() ([]*model.Course, error) {
	log.Println("[courses_repository] find...")
	var courses []*model.Course
	query, err := sqlReader.ReadQuery("courses/list_courses")
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("[courses_repository] find: error(%v)!", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var course model.Course
		var categoryID string
		if err = rows.Err(); err != nil {
			log.Printf("[courses_repository] find: error(%v)!", err)
			return nil, err
		}
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &categoryID); err == nil {
			courses = append(courses, &course)
		} else {
			log.Printf("[courses_repository] find: error(%v)!", err)
		}
	}

	log.Printf("[courses_repository] courses founded: [%v]...", courses)
	return courses, nil
}

func (r *CoursesRepository) ListByCategory(category *model.Category) ([]*model.Course, error) {
	log.Printf("[courses_repository] list_by_category: (%s)...", category.ID)
	query, err := sqlReader.ReadQuery("courses/list_courses_by_category_id")
	if err != nil {
		log.Printf("[courses_repository] list_by_category: error(%v)!", err)
		return nil, err
	}
	rows, err := r.db.Query(query, category.ID)
	if err != nil {
		log.Printf("[courses_repository] list_by_category: error(%v)!", err)
		return nil, err
	}

	var courses []*model.Course
	for rows.Next() {
		var course model.Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description); err == nil {
			courses = append(courses, &course)
		}
		if err = rows.Err(); err != nil {
			log.Printf("[courses_repository] list_by_category: error(%v)!", err)
			return nil, fmt.Errorf("courses with category (category_id: %s)", category.ID)
		}
	}

	log.Printf("[courses_repository] courses founded: [%v]...", courses)
	return courses, nil
}

func NewCoursesRepository(db *sql.DB, categoriesRepository *categories.CategoriesRepository) *CoursesRepository {
	return &CoursesRepository{
		db:         db,
		categories: categoriesRepository,
	}
}
