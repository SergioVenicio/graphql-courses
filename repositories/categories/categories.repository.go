package categories

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/sergiovenicio/graphql-courses/graph/model"
	sqlReader "github.com/sergiovenicio/graphql-courses/sql"
)

type CategoriesRepository struct {
	db *sql.DB
}

func (r *CategoriesRepository) Create(input model.NewCategory) (*model.Category, error) {
	log.Printf("[categories_repository] create: (%v)...", input)
	id, err := uuid.NewUUID()
	if err != nil {
		log.Printf("[categories_repository] create: error(%v)!", err)
		return nil, err
	}
	category := &model.Category{
		ID:          id.String(),
		Name:        input.Name,
		Description: input.Description,
	}

	query, err := sqlReader.ReadQuery("categories/insert_category")
	if err != nil {
		log.Printf("[categories_repository] create: error(%v)!", err)
		return nil, err
	}
	_, err = r.db.Exec(query, category.ID, category.Name, category.Description)
	if err != nil {
		log.Printf("[categories_repository] create: error(%v)!", err)
		return nil, err
	}

	log.Printf("[categories_repository] create: created(%v)!", category)
	return category, nil
}

func (r *CategoriesRepository) Find() ([]*model.Category, error) {
	log.Println("[categories_repository] find...")
	query, err := sqlReader.ReadQuery("categories/list_categories")
	if err != nil {
		log.Printf("[categories_repository] find: error(%v)...", err)
		return nil, err
	}
	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("[categories_repository] find: error(%v)...", err)
		return nil, err
	}
	defer rows.Close()

	var categories []*model.Category
	for rows.Next() {
		var category model.Category
		if err = rows.Err(); err != nil {
			log.Printf("[categories_repository] find: error(%v)...", err)
			return nil, err
		}
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err == nil {
			categories = append(categories, &category)
		} else {
			log.Printf("[categories_repository] find: error(%v)...", err)
			return nil, err
		}
	}

	log.Printf("[categories_repository] find: founded[%v]!", categories)
	return categories, nil
}

func (r *CategoriesRepository) FindById(id string) (*model.Category, error) {
	log.Printf("[categories_repository] find_by_id: (%s)...", id)
	var category model.Category
	query, err := sqlReader.ReadQuery("categories/get_category_by_id")
	if err != nil {
		log.Printf("[categories_repository] find_by_id: error(%v)...", err)
		return nil, err
	}

	row := r.db.QueryRow(query, id)
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		log.Printf("[categories_repository] find_by_id: error(%v)...", err)
		return nil, fmt.Errorf("category not found (id: %s)", id)
	}

	log.Printf("[categories_repository] find_by_id: founded(%v)...", category)
	return &category, nil
}

func (r *CategoriesRepository) FindByCourse(course *model.Course) (*model.Category, error) {
	log.Printf("[categories_repository] find_by_course: (%v)...", course)
	query, err := sqlReader.ReadQuery("categories/get_category_by_course_id")
	if err != nil {
		log.Printf("[categories_repository] find_by_course: error(%v)...", err)
		return nil, err
	}
	var category model.Category
	row := r.db.QueryRow(query, course.ID)
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		log.Printf("[categories_repository] find_by_course: error(%v)...", err)
		return nil, fmt.Errorf("cant find category from course (course_id: %s)", course.ID)
	}

	log.Printf("[categories_repository] find_by_id: founded(%v)...", category)
	return &category, nil
}

func NewCategoriesRepository(db *sql.DB) *CategoriesRepository {
	return &CategoriesRepository{db: db}
}
