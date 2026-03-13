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
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  category: string;
  cuisine: string;
  tags: string[];
  ingredients: Ingredient[];
  instructions: string[];
  is_favorite?: boolean;
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
  day: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  locked: boolean;
}

export interface MealPlan {
  id: number;
  user_id: number;
  week_start: string;
  week_end: string;
  entries: MealPlanEntry[];
  created_at: string;
}

// Tracking types
export interface FoodLog {
  id: number;
  user_id: number;
  recipe_id?: number;
  recipe?: Recipe;
  name: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  amount: number;
  unit: string;
  logged_at: string;
}

export interface FoodLogRequest {
  recipe_id?: number;
  name: string;
  meal_type: 'breakfast' | 'lunch' | 'dinner' | 'snack';
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
  amount: number;
  unit: string;
  logged_at?: string;
}

export interface WaterLog {
  id: number;
  user_id: number;
  amount_ml: number;
  logged_at: string;
}

export interface WaterLogRequest {
  amount_ml: number;
  logged_at?: string;
}

export interface NutritionSummary {
  date: string;
  total_calories: number;
  total_protein_g: number;
  total_carbs_g: number;
  total_fat_g: number;
  calorie_target: number;
  protein_target_g: number;
  carbs_target_g: number;
  fat_target_g: number;
  water_ml: number;
  water_target_ml: number;
  meals: {
    breakfast: FoodLog[];
    lunch: FoodLog[];
    dinner: FoodLog[];
    snack: FoodLog[];
  };
}
