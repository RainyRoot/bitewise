import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  SafeAreaView,
  ActivityIndicator,
} from 'react-native';
import { useFocusEffect } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { achievements as achievementsApi } from '@/services/api';
import type { Achievement, UserAchievement } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const ICON_MAP: Record<string, string> = {
  nutrition: 'nutrition-outline',
  flame: 'flame-outline',
  water: 'water-outline',
  restaurant: 'restaurant-outline',
  calendar: 'calendar-outline',
  barcode: 'barcode-outline',
  leaf: 'leaf-outline',
  sunny: 'sunny-outline',
  heart: 'heart-outline',
  trophy: 'trophy-outline',
};

export default function AchievementsScreen() {
  const [allAchievements, setAllAchievements] = useState<Achievement[]>([]);
  const [unlocked, setUnlocked] = useState<UserAchievement[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    try {
      const [all, mine] = await Promise.all([
        achievementsApi.getAll().catch(() => []),
        achievementsApi.getMine().catch(() => []),
      ]);
      setAllAchievements(all || []);
      setUnlocked(mine || []);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchData();
    }, [])
  );

  const unlockedIds = new Set((unlocked || []).map((u) => u.achievement_id));
  const unlockedMap = new Map((unlocked || []).map((u) => [u.achievement_id, u]));

  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <ActivityIndicator size="large" color={PRIMARY} style={{ marginTop: 48 }} />
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Text style={styles.title}>Achievements</Text>
        <Text style={styles.subtitle}>
          {unlockedIds.size} von {allAchievements.length} freigeschaltet
        </Text>

        <View style={styles.progressBar}>
          <View
            style={[
              styles.progressFill,
              {
                width: allAchievements.length > 0
                  ? `${(unlockedIds.size / allAchievements.length) * 100}%`
                  : '0%',
              },
            ]}
          />
        </View>

        {allAchievements.map((a) => {
          const isUnlocked = unlockedIds.has(a.id);
          const ua = unlockedMap.get(a.id);
          const iconName = ICON_MAP[a.icon] || 'star-outline';

          return (
            <View
              key={a.id}
              style={[styles.achievementCard, !isUnlocked && styles.achievementLocked]}
            >
              <View style={[styles.iconCircle, isUnlocked && styles.iconCircleUnlocked]}>
                <Ionicons
                  name={iconName as any}
                  size={28}
                  color={isUnlocked ? '#fff' : '#BDBDBD'}
                />
              </View>
              <View style={styles.achievementInfo}>
                <Text style={[styles.achievementName, !isUnlocked && styles.textLocked]}>
                  {a.description}
                </Text>
                <Text style={styles.achievementCategory}>{a.category}</Text>
                {isUnlocked && ua && (
                  <Text style={styles.unlockedDate}>
                    Freigeschaltet am {new Date(ua.unlocked_at).toLocaleDateString('de-DE')}
                  </Text>
                )}
              </View>
              {isUnlocked && (
                <Ionicons name="checkmark-circle" size={24} color={PRIMARY} />
              )}
            </View>
          );
        })}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  title: { fontSize: 28, fontWeight: 'bold', color: '#212121', marginBottom: 4 },
  subtitle: { fontSize: 16, color: '#757575', marginBottom: 16 },
  progressBar: {
    height: 8,
    backgroundColor: '#E0E0E0',
    borderRadius: 4,
    overflow: 'hidden',
    marginBottom: 24,
  },
  progressFill: { height: '100%', backgroundColor: PRIMARY, borderRadius: 4 },
  achievementCard: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 16,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
  },
  achievementLocked: { opacity: 0.55 },
  iconCircle: {
    width: 52,
    height: 52,
    borderRadius: 26,
    backgroundColor: '#F5F5F5',
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 14,
  },
  iconCircleUnlocked: { backgroundColor: PRIMARY },
  achievementInfo: { flex: 1 },
  achievementName: { fontSize: 15, fontWeight: '600', color: '#212121' },
  textLocked: { color: '#9E9E9E' },
  achievementCategory: { fontSize: 12, color: '#757575', marginTop: 2, textTransform: 'capitalize' },
  unlockedDate: { fontSize: 11, color: PRIMARY, marginTop: 4 },
});
