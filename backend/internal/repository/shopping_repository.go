package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainyroot/bitewise/backend/internal/domain"
)

type ShoppingListRepository interface {
	Create(ctx context.Context, list *domain.ShoppingList) error
	GetByID(ctx context.Context, id int64) (*domain.ShoppingList, error)
	GetCurrentByUser(ctx context.Context, userID int64) (*domain.ShoppingList, error)
	CreateItems(ctx context.Context, listID int64, items []domain.ShoppingListItem) error
	GetItems(ctx context.Context, listID int64) ([]domain.ShoppingListItem, error)
	ToggleItem(ctx context.Context, itemID, userID int64) error
	DeleteByID(ctx context.Context, id, userID int64) error
}

type SQLiteShoppingListRepository struct {
	db *sql.DB
}

func NewSQLiteShoppingListRepository(db *sql.DB) *SQLiteShoppingListRepository {
	return &SQLiteShoppingListRepository{db: db}
}

func (r *SQLiteShoppingListRepository) Create(ctx context.Context, list *domain.ShoppingList) error {
	list.CreatedAt = time.Now().UTC()
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO shopping_lists (user_id, meal_plan_id, created_at) VALUES (?, ?, ?)`,
		list.UserID, list.MealPlanID, list.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting shopping list: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("getting last insert id: %w", err)
	}
	list.ID = id
	return nil
}

func (r *SQLiteShoppingListRepository) GetByID(ctx context.Context, id int64) (*domain.ShoppingList, error) {
	var list domain.ShoppingList
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, meal_plan_id, created_at FROM shopping_lists WHERE id = ?`, id,
	).Scan(&list.ID, &list.UserID, &list.MealPlanID, &list.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shopping list not found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying shopping list: %w", err)
	}

	items, err := r.GetItems(ctx, list.ID)
	if err != nil {
		return nil, err
	}
	list.Items = items
	return &list, nil
}

func (r *SQLiteShoppingListRepository) GetCurrentByUser(ctx context.Context, userID int64) (*domain.ShoppingList, error) {
	var list domain.ShoppingList
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, meal_plan_id, created_at FROM shopping_lists
		 WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID,
	).Scan(&list.ID, &list.UserID, &list.MealPlanID, &list.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no shopping list found")
	}
	if err != nil {
		return nil, fmt.Errorf("querying current shopping list: %w", err)
	}

	items, err := r.GetItems(ctx, list.ID)
	if err != nil {
		return nil, err
	}
	list.Items = items
	return &list, nil
}

func (r *SQLiteShoppingListRepository) CreateItems(ctx context.Context, listID int64, items []domain.ShoppingListItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback()

	for i := range items {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO shopping_list_items (shopping_list_id, ingredient_name, quantity, unit, category, is_checked) VALUES (?, ?, ?, ?, ?, 0)`,
			listID, items[i].IngredientName, items[i].Quantity, items[i].Unit, items[i].Category,
		)
		if err != nil {
			return fmt.Errorf("inserting shopping item: %w", err)
		}
		id, _ := result.LastInsertId()
		items[i].ID = id
		items[i].ShoppingListID = listID
	}

	return tx.Commit()
}

func (r *SQLiteShoppingListRepository) GetItems(ctx context.Context, listID int64) ([]domain.ShoppingListItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, shopping_list_id, ingredient_name, quantity, unit, category, is_checked
		 FROM shopping_list_items WHERE shopping_list_id = ? ORDER BY category, ingredient_name`, listID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying shopping items: %w", err)
	}
	defer rows.Close()

	var items []domain.ShoppingListItem
	for rows.Next() {
		var item domain.ShoppingListItem
		var checked int
		if err := rows.Scan(&item.ID, &item.ShoppingListID, &item.IngredientName, &item.Quantity, &item.Unit, &item.Category, &checked); err != nil {
			return nil, fmt.Errorf("scanning shopping item: %w", err)
		}
		item.IsChecked = checked != 0
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *SQLiteShoppingListRepository) ToggleItem(ctx context.Context, itemID, userID int64) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE shopping_list_items SET is_checked = 1 - is_checked
		 WHERE id = ? AND shopping_list_id IN (SELECT id FROM shopping_lists WHERE user_id = ?)`,
		itemID, userID,
	)
	if err != nil {
		return fmt.Errorf("toggling item: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

func (r *SQLiteShoppingListRepository) DeleteByID(ctx context.Context, id, userID int64) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM shopping_lists WHERE id = ? AND user_id = ?`, id, userID,
	)
	if err != nil {
		return fmt.Errorf("deleting shopping list: %w", err)
	}
	return nil
}
