import AsyncStorage from '@react-native-async-storage/async-storage';
import type {
  AuthResponse,
  FoodLog,
  FoodLogRequest,
  LoginRequest,
  MealPlan,
  MealPlanEntry,
  NutritionSummary,
  Recipe,
  RecipeFilter,
  RegisterRequest,
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
    throw new Error(error.message || `HTTP ${response.status}`);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json();
};

// Auth
export const auth = {
  register: async (data: RegisterRequest): Promise<AuthResponse> => {
    const res = await request<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    await setToken(res.token);
    return res;
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    const res = await request<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    await setToken(res.token);
    return res;
  },

  getMe: async (): Promise<User> => {
    return request<User>('/api/auth/me');
  },

  logout: async (): Promise<void> => {
    await clearToken();
  },
};

// Profile
export const profile = {
  getProfile: async (): Promise<User> => {
    return request<User>('/api/profile');
  },

  updateProfile: async (data: Partial<User>): Promise<User> => {
    return request<User>('/api/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  updateAllergies: async (allergies: string[]): Promise<void> => {
    return request<void>('/api/profile/allergies', {
      method: 'PUT',
      body: JSON.stringify({ allergies }),
    });
  },

  updatePreferences: async (preferences: Record<string, string>): Promise<void> => {
    return request<void>('/api/profile/preferences', {
      method: 'PUT',
      body: JSON.stringify(preferences),
    });
  },
};

// Recipes
export const recipes = {
  list: async (filter?: RecipeFilter): Promise<Recipe[]> => {
    const params = new URLSearchParams();
    if (filter?.query) params.set('query', filter.query);
    if (filter?.category) params.set('category', filter.category);
    if (filter?.cuisine) params.set('cuisine', filter.cuisine);
    if (filter?.max_calories) params.set('max_calories', String(filter.max_calories));
    if (filter?.max_prep_time) params.set('max_prep_time', String(filter.max_prep_time));
    if (filter?.tags) params.set('tags', filter.tags.join(','));
    if (filter?.page) params.set('page', String(filter.page));
    if (filter?.limit) params.set('limit', String(filter.limit));
    const qs = params.toString();
    return request<Recipe[]>(`/api/recipes${qs ? `?${qs}` : ''}`);
  },

  getById: async (id: number): Promise<Recipe> => {
    return request<Recipe>(`/api/recipes/${id}`);
  },

  addFavorite: async (recipeId: number): Promise<void> => {
    return request<void>(`/api/recipes/${recipeId}/favorite`, {
      method: 'POST',
    });
  },

  removeFavorite: async (recipeId: number): Promise<void> => {
    return request<void>(`/api/recipes/${recipeId}/favorite`, {
      method: 'DELETE',
    });
  },
};

// Meal Plans
export const mealPlans = {
  generate: async (): Promise<MealPlan> => {
    return request<MealPlan>('/api/meal-plans/generate', {
      method: 'POST',
    });
  },

  getCurrent: async (): Promise<MealPlan> => {
    return request<MealPlan>('/api/meal-plans/current');
  },

  getById: async (id: number): Promise<MealPlan> => {
    return request<MealPlan>(`/api/meal-plans/${id}`);
  },

  updateEntry: async (planId: number, entryId: number, data: Partial<MealPlanEntry>): Promise<MealPlanEntry> => {
    return request<MealPlanEntry>(`/api/meal-plans/${planId}/entries/${entryId}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  },

  lockEntry: async (planId: number, entryId: number, locked: boolean): Promise<MealPlanEntry> => {
    return request<MealPlanEntry>(`/api/meal-plans/${planId}/entries/${entryId}/lock`, {
      method: 'PUT',
      body: JSON.stringify({ locked }),
    });
  },

  regenerate: async (planId: number): Promise<MealPlan> => {
    return request<MealPlan>(`/api/meal-plans/${planId}/regenerate`, {
      method: 'POST',
    });
  },
};

// Food Tracking
export const tracking = {
  logFood: async (data: FoodLogRequest): Promise<FoodLog> => {
    return request<FoodLog>('/api/tracking/food', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getFoodLogs: async (date: string): Promise<FoodLog[]> => {
    return request<FoodLog[]>(`/api/tracking/food?date=${date}`);
  },

  deleteFood: async (id: number): Promise<void> => {
    return request<void>(`/api/tracking/food/${id}`, {
      method: 'DELETE',
    });
  },

  getSummary: async (date: string): Promise<NutritionSummary> => {
    return request<NutritionSummary>(`/api/tracking/summary?date=${date}`);
  },
};

// Water Tracking
export const water = {
  logWater: async (data: WaterLogRequest): Promise<WaterLog> => {
    return request<WaterLog>('/api/tracking/water', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getWaterLogs: async (date: string): Promise<WaterLog[]> => {
    return request<WaterLog[]>(`/api/tracking/water?date=${date}`);
  },

  getWaterSummary: async (date: string): Promise<{ total_ml: number; target_ml: number }> => {
    return request<{ total_ml: number; target_ml: number }>(`/api/tracking/water/summary?date=${date}`);
  },
};

export const api = {
  auth,
  profile,
  recipes,
  mealPlans,
  tracking,
  water,
};

export default api;
