package db

import (
	"context"
	"eventy/pkg/models"
	"fmt"

	"github.com/rs/zerolog/log"
)

// GetAllCategories retrieves all categories from the database
func GetAllCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := Db_GlobalVar.NewSelect().Model(&categories).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting all categories: %w", err)
	}
	return categories, nil
}

// GetCategoryByID retrieves a single category by its ID
func GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	category := new(models.Category)
	err := Db_GlobalVar.NewSelect().Model(category).Where("category_id = ?", id).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting category by ID %d: %w", id, err)
	}
	return category, nil
}

// AddCategory creates a new category in the database
func AddCategory(ctx context.Context, category *models.Category) error {
	_, err := Db_GlobalVar.NewInsert().Model(category).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error creating category: %w", err)
	}
	log.Debug().Msgf("New category added with ID: %d", category.CategoryID)
	return nil
}

// UpdateCategory updates an existing category in the database
func UpdateCategory(ctx context.Context, id int, updates *models.Category) (int64, error) {
	res, err := Db_GlobalVar.NewUpdate().
		Model(updates).
		Where("category_id = ?", id).
		OmitZero().
		Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error updating category with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Updated category with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}

// DeleteCategory removes a category from the database by its ID
func DeleteCategory(ctx context.Context, id int) (int64, error) {
	res, err := Db_GlobalVar.NewDelete().Model(&models.Category{}).Where("category_id = ?", id).Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error deleting category with ID %d: %w", id, err)
	}

	rowsAffected, _ := res.RowsAffected()
	log.Debug().Msgf("Deleted category with ID: %d, rows affected: %d", id, rowsAffected)
	return rowsAffected, nil
}
