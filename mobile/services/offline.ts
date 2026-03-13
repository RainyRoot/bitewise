import AsyncStorage from '@react-native-async-storage/async-storage';

const CACHE_PREFIX = 'bitewise_cache_';
const QUEUE_KEY = 'bitewise_offline_queue';

// Cache API responses for offline viewing
export const offlineCache = {
  set: async <T>(key: string, data: T, ttlMs: number = 30 * 60 * 1000): Promise<void> => {
    const entry = {
      data,
      expires: Date.now() + ttlMs,
    };
    await AsyncStorage.setItem(CACHE_PREFIX + key, JSON.stringify(entry));
  },

  get: async <T>(key: string): Promise<T | null> => {
    const raw = await AsyncStorage.getItem(CACHE_PREFIX + key);
    if (!raw) return null;
    try {
      const entry = JSON.parse(raw);
      if (entry.expires && entry.expires < Date.now()) {
        await AsyncStorage.removeItem(CACHE_PREFIX + key);
        return null;
      }
      return entry.data as T;
    } catch {
      return null;
    }
  },

  remove: async (key: string): Promise<void> => {
    await AsyncStorage.removeItem(CACHE_PREFIX + key);
  },
};

// Queue offline actions for later sync
interface QueuedAction {
  id: string;
  type: 'food_log' | 'water_log';
  data: any;
  createdAt: number;
}

export const offlineQueue = {
  add: async (type: QueuedAction['type'], data: any): Promise<void> => {
    const queue = await offlineQueue.getAll();
    queue.push({
      id: `${Date.now()}_${Math.random().toString(36).slice(2)}`,
      type,
      data,
      createdAt: Date.now(),
    });
    await AsyncStorage.setItem(QUEUE_KEY, JSON.stringify(queue));
  },

  getAll: async (): Promise<QueuedAction[]> => {
    const raw = await AsyncStorage.getItem(QUEUE_KEY);
    if (!raw) return [];
    try {
      return JSON.parse(raw);
    } catch {
      return [];
    }
  },

  remove: async (id: string): Promise<void> => {
    const queue = await offlineQueue.getAll();
    const filtered = queue.filter((a) => a.id !== id);
    await AsyncStorage.setItem(QUEUE_KEY, JSON.stringify(filtered));
  },

  clear: async (): Promise<void> => {
    await AsyncStorage.removeItem(QUEUE_KEY);
  },

  size: async (): Promise<number> => {
    const queue = await offlineQueue.getAll();
    return queue.length;
  },
};

// Sync queued actions when back online
export const syncOfflineQueue = async (
  handlers: {
    food_log: (data: any) => Promise<void>;
    water_log: (data: any) => Promise<void>;
  },
): Promise<number> => {
  const queue = await offlineQueue.getAll();
  let synced = 0;

  for (const action of queue) {
    try {
      const handler = handlers[action.type];
      if (handler) {
        await handler(action.data);
        await offlineQueue.remove(action.id);
        synced++;
      }
    } catch {
      // Keep in queue for next sync attempt
      break;
    }
  }

  return synced;
};

// Check if device is online
export const isOnline = async (): Promise<boolean> => {
  try {
    const controller = new AbortController();
    const timeout = setTimeout(() => controller.abort(), 3000);
    const response = await fetch(
      `${process.env.EXPO_PUBLIC_API_URL || 'http://localhost:8080'}/health`,
      { signal: controller.signal },
    );
    clearTimeout(timeout);
    return response.ok;
  } catch {
    return false;
  }
};
