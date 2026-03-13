import AsyncStorage from '@react-native-async-storage/async-storage';
import type {
  AuthResponse,
  FoodItem,
  FoodLog,
  FoodLogRequest,
  LoginRequest,
  MealPlan,
  MealPlanEntry,
  NutritionSummary,
  PantryItem,
  PantryMatch,
  Recipe,
  RecipeFilter,
  RegisterRequest,
  SeasonalResponse,
  ShoppingList,
  User,
  WaterLog,
  WaterLogRequest,
} from '@/types';

const TOKEN_KEY = 'bitewise_token';
const BASE_URL = process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080';

let cachedToken: string | null = null;

const getToken = async (): Promise<string | null> => {
  if (cachedToken) return cachedToken;
  cachedToken = await AsyncStorage.getItem(TOKEN_KEY);
  return cachedToken;
};

const setToken = async (token: string): Promise<void> => {
  cachedToken = token;
  await AsyncStorage.setItem(TOKEN_KEY, token);
};

const clearToken = async (): Promise<void> => {
  cachedToken = null;
  await AsyncStorage.removeItem(TOKEN_KEY);
};

const request = async <T>(
  path: string,
  options: RequestInit = {},
): Promise<T> => {
  const token = await getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ message: 'Request failed' }));
    throw new Error(error.message || error.error || `HTTP ${response.status}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
};

// Auth
export const auth = {
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const res = await request<AuthResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    await setToken(res.token);
    return res;
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const res = await request<AuthResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    await setToken(res.token);
    return res;
  },

  getMe: async (): Promise<User> => {
    return request<User>('/api/v1/profile');
  },

  logout: async (): Promise<void> => {
    await clearToken();
  },
};

// Profile
export const profile = {
  getProfile: async (): Promise<User> => {
    return request<User>('/api/v1/profile');
  },

  updateProfile: async (data: Partial<User>): Promise<User> => {
    return request<User>('/api/v1/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  getAllergies: async (): Promise<string[]> => {
    const res = await request<{ allergen: string }[]>('/api/v1/profile/allergies');
    return (res || []).map((a) => a.allergen);
  },

  updateAllergies: async (allergies: string[]): Promise<void> => {
    return request<void>('/api/v1/profile/allergies', {
      method: 'PUT',
      body: JSON.stringify({ allergies }),
    });
  },

  updatePreferences: async (preferences: Record<string, string>): Promise<void> => {
    return request<void>('/api/v1/profile/preferences', {
      method: 'PUT',
      body: JSON.stringify(preferences),
    });
  },
};

// Recipes
export const recipes = {
  list: async (filter?: RecipeFilter): Promise<{ recipes: Recipe[]; total: number }> => {
    const params = new URLSearchParams();
    if (filter?.query) params.set('q', filter.query);
    if (filter?.category) params.set('category', filter.category);
    if (filter?.max_prep_time) params.set('max_prep_time', String(filter.max_prep_time));
    if (filter?.limit) params.set('limit', String(filter.limit));
    if (filter?.page) params.set('offset', String(((filter.page || 1) - 1) * (filter.limit || 20)));
    const qs = params.toString();
    return request<{ recipes: Recipe[]; total: number }>(`/api/v1/recipes${qs ? `?${qs}` : ''}`);
  },

  getById: async (id: number): Promise<Recipe> => {
    return request<Recipe>(`/api/v1/recipes/${id}`);
  },

  getFavorites: async (): Promise<Recipe[]> => {
    return request<Recipe[]>('/api/v1/recipes/favorites');
  },

  addFavorite: async (recipeId: number): Promise<void> => {
    return request<void>(`/api/v1/recipes/${recipeId}/favorite`, {
      method: 'POST',
    });
  },

  removeFavorite: async (recipeId: number): Promise<void> => {
    return request<void>(`/api/v1/recipes/${recipeId}/favorite`, {
      method: 'DELETE',
    });
  },
};

// Meal Plans
export const mealPlans = {
  generate: async (weekStartDate?: string): Promise<MealPlan> => {
    return request<MealPlan>('/api/v1/meal-plans/generate', {
      method: 'POST',
      body: JSON.stringify({ week_start_date: weekStartDate || '' }),
    });
  },

  getCurrent: async (): Promise<MealPlan> => {
    return request<MealPlan>('/api/v1/meal-plans/current');
  },

  getById: async (id: number): Promise<MealPlan> => {
    return request<MealPlan>(`/api/v1/meal-plans/${id}`);
  },

  updateEntry: async (planId: number, entryId: number, data: Partial<MealPlanEntry>): Promise<MealPlanEntry> => {
    return request<MealPlanEntry>(`/api/v1/meal-plans/${planId}/entries/${entryId}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    });
  },

  activate: async (planId: number): Promise<void> => {
    return request<void>(`/api/v1/meal-plans/${planId}/activate`, {
      method: 'POST',
    });
  },
};

// Food Tracking
export const tracking = {
  logFood: async (data: FoodLogRequest): Promise<FoodLog> => {
    return request<FoodLog>('/api/v1/tracking/food', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getFoodLogs: async (date: string): Promise<FoodLog[]> => {
    return request<FoodLog[]>(`/api/v1/tracking/food?date=${date}`);
  },

  deleteFood: async (id: number): Promise<void> => {
    return request<void>(`/api/v1/tracking/food/${id}`, {
      method: 'DELETE',
    });
  },

  getSummary: async (date: string): Promise<NutritionSummary> => {
    return request<NutritionSummary>(`/api/v1/tracking/summary?date=${date}`);
  },
};

// Water Tracking
export const water = {
  logWater: async (data: WaterLogRequest): Promise<WaterLog> => {
    return request<WaterLog>('/api/v1/tracking/water', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getWaterLogs: async (date: string): Promise<WaterLog[]> => {
    return request<WaterLog[]>(`/api/v1/tracking/water?date=${date}`);
  },
};

// Nutrition (Barcode / Food Search)
export const nutritionLookup = {
  lookupBarcode: async (code: string): Promise<FoodItem> => {
    return request<FoodItem>(`/api/v1/nutrition/barcode/${code}`);
  },

  searchFood: async (query: string): Promise<FoodItem[]> => {
    return request<FoodItem[]>(`/api/v1/nutrition/search?q=${encodeURIComponent(query)}`);
  },
};

// Shopping Lists
export const shoppingLists = {
  generate: async (mealPlanId: number): Promise<ShoppingList> => {
    return request<ShoppingList>('/api/v1/shopping-lists', {
      method: 'POST',
      body: JSON.stringify({ meal_plan_id: mealPlanId }),
    });
  },

  getCurrent: async (): Promise<ShoppingList> => {
    return request<ShoppingList>('/api/v1/shopping-lists/current');
  },

  toggleItem: async (itemId: number): Promise<void> => {
    return request<void>(`/api/v1/shopping-lists/items/${itemId}`, {
      method: 'PATCH',
    });
  },
};

// Pantry
export const pantry = {
  setItems: async (items: string[]): Promise<PantryItem[]> => {
    return request<PantryItem[]>('/api/v1/pantry', {
      method: 'POST',
      body: JSON.stringify({ items }),
    });
  },

  getItems: async (): Promise<PantryItem[]> => {
    return request<PantryItem[]>('/api/v1/pantry');
  },

  findRecipes: async (): Promise<PantryMatch[]> => {
    return request<PantryMatch[]>('/api/v1/pantry/recipes');
  },
};

// Seasonal
export const seasonal = {
  getCurrent: async (): Promise<SeasonalResponse> => {
    return request<SeasonalResponse>('/api/v1/seasonal');
  },
};

export const api = {
  auth,
  profile,
  recipes,
  mealPlans,
  tracking,
  water,
  nutritionLookup,
  shoppingLists,
  pantry,
  seasonal,
};

export default api;
