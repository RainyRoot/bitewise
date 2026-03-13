package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/rainyroot/bitewise/backend/internal/domain"
	"github.com/rainyroot/bitewise/backend/internal/repository"
)

type ShoppingService struct {
	shopping  repository.ShoppingListRepository
	mealPlans repository.MealPlanRepository
}

func NewShoppingService(shopping repository.ShoppingListRepository, mealPlans repository.MealPlanRepository) *ShoppingService {
	return &ShoppingService{shopping: shopping, mealPlans: mealPlans}
}

func (s *ShoppingService) GenerateFromMealPlan(ctx context.Context, userID, mealPlanID int64) (*domain.ShoppingList, error) {
	plan, err := s.mealPlans.GetByID(ctx, mealPlanID)
	if err != nil {
		return nil, fmt.Errorf("getting meal plan: %w", err)
	}
	if plan.UserID != userID {
		return nil, fmt.Errorf("meal plan not found")
	}

	// Collect and merge ingredients from all entries
	type ingredientKey struct {
		name string
		unit string
	}
	merged := map[ingredientKey]*domain.ShoppingListItem{}

	for _, entry := range plan.Entries {
		if entry.Recipe == nil {
			continue
		}
		for _, ing := range entry.Recipe.Ingredients {
			key := ingredientKey{
				name: strings.ToLower(ing.Name),
				unit: strings.ToLower(ing.Unit),
			}
			if existing, ok := merged[key]; ok {
				existing.Quantity += ing.Quantity * float64(entry.Servings)
			} else {
				merged[key] = &domain.ShoppingListItem{
					IngredientName: ing.Name,
					Quantity:       ing.Quantity * float64(entry.Servings),
					Unit:           ing.Unit,
					Category:       categorizeIngredient(ing.Name),
				}
			}
		}
	}

	list := &domain.ShoppingList{
		UserID:     userID,
		MealPlanID: mealPlanID,
	}
	if err := s.shopping.Create(ctx, list); err != nil {
		return nil, fmt.Errorf("creating shopping list: %w", err)
	}

	items := make([]domain.ShoppingListItem, 0, len(merged))
	for _, item := range merged {
		items = append(items, *item)
	}

	if len(items) > 0 {
		if err := s.shopping.CreateItems(ctx, list.ID, items); err != nil {
			return nil, fmt.Errorf("creating items: %w", err)
		}
	}

	return s.shopping.GetByID(ctx, list.ID)
}

func (s *ShoppingService) GetCurrent(ctx context.Context, userID int64) (*domain.ShoppingList, error) {
	return s.shopping.GetCurrentByUser(ctx, userID)
}

func (s *ShoppingService) ToggleItem(ctx context.Context, itemID, userID int64) error {
	return s.shopping.ToggleItem(ctx, itemID, userID)
}

// categorizeIngredient assigns a German category based on ingredient name.
func categorizeIngredient(name string) string {
	lower := strings.ToLower(name)

	fruitsVeggies := []string{"apfel", "banane", "tomate", "gurke", "salat", "spinat", "brokkoli",
		"paprika", "zwiebel", "kartoffel", "karotte", "möhre", "zucchini", "aubergine",
		"beeren", "erdbeere", "himbeere", "birne", "orange", "zitrone", "knoblauch",
		"pilz", "champignon", "kürbis", "avocado", "mais", "erbsen", "bohnen", "lauch"}
	for _, f := range fruitsVeggies {
		if strings.Contains(lower, f) {
			return "Obst & Gemüse"
		}
	}

	dairy := []string{"milch", "joghurt", "käse", "quark", "sahne", "butter", "schmand", "crème"}
	for _, d := range dairy {
		if strings.Contains(lower, d) {
			return "Milchprodukte"
		}
	}

	meat := []string{"hähnchen", "huhn", "rind", "schwein", "hack", "schinken", "wurst",
		"lachs", "fisch", "thunfisch", "garnele", "fleisch", "pute"}
	for _, m := range meat {
		if strings.Contains(lower, m) {
			return "Fleisch & Fisch"
		}
	}

	grains := []string{"mehl", "brot", "nudel", "pasta", "reis", "hafer", "müsli",
		"toast", "semmel", "vollkorn", "quinoa", "couscous"}
	for _, g := range grains {
		if strings.Contains(lower, g) {
			return "Getreide & Backwaren"
		}
	}

	spices := []string{"salz", "pfeffer", "zimt", "oregano", "basilikum", "thymian",
		"paprikapulver", "curry", "muskat", "koriander", "gewürz"}
	for _, sp := range spices {
		if strings.Contains(lower, sp) {
			return "Gewürze"
		}
	}

	oils := []string{"öl", "olivenöl", "essig", "sojasoße", "senf", "ketchup", "honig"}
	for _, o := range oils {
		if strings.Contains(lower, o) {
			return "Öle & Soßen"
		}
	}

	return "Sonstiges"
}
