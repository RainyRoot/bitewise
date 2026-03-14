import AsyncStorage from '@react-native-async-storage/async-storage';
import type {
  Achievement,
  AuthResponse,
  CreateRecipeRequest,
  DiaryEntry,
  DiaryEntryRequest,
  FoodItem,
  FoodLog,
  FoodLogRequest,
  FriendInfo,
  FriendInvite,
  LeaderboardEntry,
  LoginRequest,
  MealPlan,
  MealPlanEntry,
  MonthlyMoodSummary,
  MonthlyStats,
  NutritionSummary,
  PantryItem,
  PantryMatch,
  PriceLog,
  PriceLogRequest,
  PriceTrend,
  Recipe,
  RecipeFilter,
  RegisterRequest,
  SeasonalResponse,
  SharedRecipe,
  ShoppingList,
  SpendingSummary,
  StoreComparison,
  StreakInfo,
  User,
  UserAchievement,
  UserNotification,
  WaterLog,
  WaterLogRequest,
  WeeklyStats,
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

  create: async (data: CreateRecipeRequest): Promise<Recipe> => {
    return request<Recipe>('/api/v1/recipes', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getMine: async (): Promise<Recipe[]> => {
    return request<Recipe[]>('/api/v1/recipes/mine');
  },

  delete: async (id: number): Promise<void> => {
    return request<void>(`/api/v1/recipes/${id}`, {
      method: 'DELETE',
    });
  },

  share: async (recipeId: number): Promise<SharedRecipe> => {
    return request<SharedRecipe>(`/api/v1/recipes/${recipeId}/share`, {
      method: 'POST',
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

// Achievements
export const achievements = {
  getAll: async (): Promise<Achievement[]> => {
    return request<Achievement[]>('/api/v1/achievements');
  },

  getMine: async (): Promise<UserAchievement[]> => {
    return request<UserAchievement[]>('/api/v1/achievements/mine');
  },
};

// Sharing
export const sharing = {
  getSharedRecipe: async (code: string): Promise<Recipe> => {
    return request<Recipe>(`/api/v1/shared/${code}`);
  },

  saveSharedRecipe: async (code: string): Promise<Recipe> => {
    return request<Recipe>(`/api/v1/shared/${code}/save`, {
      method: 'POST',
    });
  },
};

// Notifications
export const notifications = {
  getSettings: async (): Promise<UserNotification[]> => {
    return request<UserNotification[]>('/api/v1/notifications');
  },

  updateSettings: async (notifs: UserNotification[]): Promise<UserNotification[]> => {
    return request<UserNotification[]>('/api/v1/notifications', {
      method: 'PUT',
      body: JSON.stringify({ notifications: notifs }),
    });
  },
};

// Statistics
export const stats = {
  getWeekly: async (): Promise<WeeklyStats> => {
    return request<WeeklyStats>('/api/v1/stats/weekly');
  },

  getMonthly: async (): Promise<MonthlyStats> => {
    return request<MonthlyStats>('/api/v1/stats/monthly');
  },

  getStreaks: async (): Promise<StreakInfo> => {
    return request<StreakInfo>('/api/v1/stats/streaks');
  },
};

// Diary
export const diary = {
  createOrUpdate: async (data: DiaryEntryRequest): Promise<DiaryEntry> => {
    return request<DiaryEntry>('/api/v1/diary', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getByDate: async (date: string): Promise<DiaryEntry | null> => {
    return request<DiaryEntry | null>(`/api/v1/diary?date=${date}`);
  },

  getMonthly: async (month: string): Promise<MonthlyMoodSummary> => {
    return request<MonthlyMoodSummary>(`/api/v1/diary/monthly?month=${month}`);
  },

  delete: async (id: number): Promise<void> => {
    return request<void>(`/api/v1/diary/${id}`, {
      method: 'DELETE',
    });
  },
};

// Prices
export const prices = {
  log: async (data: PriceLogRequest): Promise<PriceLog> => {
    return request<PriceLog>('/api/v1/prices', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  getLogs: async (limit?: number): Promise<PriceLog[]> => {
    const qs = limit ? `?limit=${limit}` : '';
    return request<PriceLog[]>(`/api/v1/prices${qs}`);
  },

  getTrend: async (item: string): Promise<PriceTrend> => {
    return request<PriceTrend>(`/api/v1/prices/trends?item=${encodeURIComponent(item)}`);
  },

  compareStores: async (item: string): Promise<StoreComparison> => {
    return request<StoreComparison>(`/api/v1/prices/compare?item=${encodeURIComponent(item)}`);
  },

  getSpending: async (month?: string): Promise<SpendingSummary> => {
    const qs = month ? `?month=${month}` : '';
    return request<SpendingSummary>(`/api/v1/prices/spending${qs}`);
  },
};

// Friends
export const friends = {
  invite: async (email: string): Promise<FriendInvite> => {
    return request<FriendInvite>('/api/v1/friends/invite', {
      method: 'POST',
      body: JSON.stringify({ email }),
    });
  },

  getPendingInvites: async (): Promise<FriendInvite[]> => {
    return request<FriendInvite[]>('/api/v1/friends/invites');
  },

  respondToInvite: async (inviteId: number, accept: boolean): Promise<void> => {
    return request<void>(`/api/v1/friends/invites/${inviteId}`, {
      method: 'POST',
      body: JSON.stringify({ accept }),
    });
  },

  list: async (): Promise<FriendInfo[]> => {
    return request<FriendInfo[]>('/api/v1/friends');
  },

  remove: async (friendId: number): Promise<void> => {
    return request<void>(`/api/v1/friends/${friendId}`, {
      method: 'DELETE',
    });
  },

  getLeaderboard: async (): Promise<LeaderboardEntry[]> => {
    return request<LeaderboardEntry[]>('/api/v1/leaderboard');
  },
};

// Export
export const dataExport = {
  downloadCSV: async (): Promise<string> => {
    const token = await getToken();
    const response = await fetch(`${BASE_URL}/api/v1/export/csv`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    return response.text();
  },

  downloadJSON: async (): Promise<string> => {
    const token = await getToken();
    const response = await fetch(`${BASE_URL}/api/v1/export/json`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    });
    return response.text();
  },

  deleteAccount: async (): Promise<void> => {
    return request<void>('/api/v1/account', {
      method: 'DELETE',
    });
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
  achievements,
  sharing,
  notifications,
  stats,
  diary,
  prices,
  friends,
  dataExport,
};

export default api;
