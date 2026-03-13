// User & Auth types
export interface User {
  id: number;
  email: string;
  name: string;
  height_cm?: number;
  weight_kg?: number;
  birth_date?: string;
  gender?: 'male' | 'female' | 'other';
  activity_level?: 'sedentary' | 'light' | 'moderate' | 'active' | 'very_active';
  goal?: 'lose' | 'maintain' | 'gain';
  calorie_target?: number;
  daily_water_ml_goal?: number;
  created_at: string;
  updated_at: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

// Recipe types
export interface Ingredient {
  id: number;
  name: string;
  amount: number;
  quantity?: number;
  unit: string;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
}

export interface Recipe {
  id: number;
  title: string;
  description: string;
  image_url?: string;
  prep_time_min: number;
  cook_time_min: number;
  servings: number;
  calories: number;
  calories_per_serving?: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  fiber_g?: number;
  difficulty?: string;
  category: string;
  cuisine: string;
  tags: string[];
  allergens?: string[];
  categories?: string[];
  ingredients: Ingredient[];
  instructions: string[];
  is_favorite?: boolean;
  source_url?: string;
  source_site?: string;
  created_at: string;
}

export interface RecipeFilter {
  query?: string;
  category?: string;
  cuisine?: string;
  max_calories?: number;
  max_prep_time?: number;
  tags?: string[];
  page?: number;
  limit?: number;
}

// Meal Plan types
export interface MealPlanEntry {
  id: number;
  meal_plan_id: number;
  recipe_id: number;
  recipe?: Recipe;
  day_of_week: number;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  servings: number;
  is_locked: boolean;
}

export interface MealPlan {
  id: number;
  user_id: number;
  week_start_date: string;
  status: string;
  entries: MealPlanEntry[];
  created_at: string;
}

// Tracking types
export interface FoodLog {
  id: number;
  user_id: number;
  date: string;
  food_name: string;
  barcode?: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  servings: number;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  fiber_g?: number;
  created_at?: string;
}

export interface FoodLogRequest {
  date?: string;
  food_name: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  barcode?: string;
  servings?: number;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  fiber_g?: number;
}

export interface WaterLog {
  id: number;
  user_id: number;
  date?: string;
  amount_ml: number;
  logged_at: string;
}

export interface WaterLogRequest {
  date?: string;
  amount_ml: number;
}

export interface NutritionSummary {
  date: string;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  fiber_g?: number;
  water_ml: number;
  meal_count: number;
}

// Nutrition / Barcode types
export interface FoodItem {
  barcode: string;
  name: string;
  brand?: string;
  image_url?: string;
  calories_per_100g: number;
  protein_g_per_100g: number;
  carbs_g_per_100g: number;
  fat_g_per_100g: number;
  fiber_g_per_100g: number;
  serving_size?: string;
}

// Shopping List types
export interface ShoppingListItem {
  id: number;
  shopping_list_id: number;
  ingredient_name: string;
  quantity: number;
  unit: string;
  category: string;
  is_checked: boolean;
}

export interface ShoppingList {
  id: number;
  user_id: number;
  meal_plan_id: number;
  items: ShoppingListItem[];
  created_at: string;
}

// Pantry types
export interface PantryItem {
  id: number;
  user_id: number;
  name: string;
  created_at: string;
}

export interface PantryMatch {
  recipe: Recipe;
  matched_count: number;
  total_count: number;
  match_percent: number;
}

// Seasonal types
export interface SeasonalItem {
  name: string;
  category: string;
}

export interface SeasonalResponse {
  month: string;
  items: SeasonalItem[];
}
