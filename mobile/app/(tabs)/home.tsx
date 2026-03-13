import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
  ActivityIndicator,
} from 'react-native';
import { useFocusEffect } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { useAuth } from '@/hooks/useAuth';
import { tracking, water } from '@/services/api';
import { offlineCache, isOnline, syncOfflineQueue } from '@/services/offline';
import type { NutritionSummary, FoodLog } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const MEAL_LABELS: Record<string, string> = {
  breakfast: 'Frühstück',
  lunch: 'Mittagessen',
  dinner: 'Abendessen',
  snack: 'Snack',
};

export default function HomeScreen() {
  const { user } = useAuth();
  const [summary, setSummary] = useState<NutritionSummary | null>(null);
  const [foodLogs, setFoodLogs] = useState<FoodLog[]>([]);
  const [waterTotal, setWaterTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [offline, setOffline] = useState(false);

  const calorieTarget = user?.calorie_target || 2200;
  const waterTarget = user?.daily_water_ml_goal || 2500;
  const today = new Date().toISOString().split('T')[0];

  const fetchData = async () => {
    const online = await isOnline();
    setOffline(!online);

    if (online) {
      // Sync any queued offline actions
      await syncOfflineQueue({
        food_log: async (data) => { await tracking.logFood(data); },
        water_log: async (data) => { await water.logWater(data); },
      });
    }

    try {
      const [summaryData, logs, waterLogs] = await Promise.all([
        tracking.getSummary(today).catch(() => null),
        tracking.getFoodLogs(today).catch(() => []),
        water.getWaterLogs(today).catch(() => []),
      ]);
      setSummary(summaryData);
      setFoodLogs(logs || []);
      const totalWater = (waterLogs || []).reduce((sum, w) => sum + w.amount_ml, 0);
      setWaterTotal(totalWater);

      // Cache for offline
      if (online) {
        await offlineCache.set('home_summary', summaryData);
        await offlineCache.set('home_food_logs', logs);
        await offlineCache.set('home_water_total', totalWater);
      }
    } catch {
      // Try loading from cache
      const cachedSummary = await offlineCache.get<NutritionSummary>('home_summary');
      const cachedLogs = await offlineCache.get<FoodLog[]>('home_food_logs');
      const cachedWater = await offlineCache.get<number>('home_water_total');
      if (cachedSummary) setSummary(cachedSummary);
      if (cachedLogs) setFoodLogs(cachedLogs);
      if (cachedWater) setWaterTotal(cachedWater);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchData();
    }, [])
  );

  const addWater = async (amount: number) => {
    try {
      await water.logWater({ date: today, amount_ml: amount });
      setWaterTotal((prev) => prev + amount);
    } catch {
      // ignore
    }
  };

  const caloriesConsumed = summary?.calories || 0;
  const caloriePercent = calorieTarget > 0 ? Math.round((caloriesConsumed / calorieTarget) * 100) : 0;
  const waterPercent = waterTarget > 0 ? Math.round((waterTotal / waterTarget) * 100) : 0;
  const proteinG = summary?.protein_g || 0;
  const carbsG = summary?.carbs_g || 0;
  const fatG = summary?.fat_g || 0;

  // Group food logs by meal type
  const mealGroups = foodLogs.reduce<Record<string, FoodLog[]>>((groups, log) => {
    if (!groups[log.meal_type]) groups[log.meal_type] = [];
    groups[log.meal_type].push(log);
    return groups;
  }, {});

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
        {offline && (
          <View style={styles.offlineBanner}>
            <Ionicons name="cloud-offline-outline" size={16} color="#fff" />
            <Text style={styles.offlineText}>Offline-Modus - Daten aus Cache</Text>
          </View>
        )}
        <Text style={styles.greeting}>Hallo, {user?.name || 'User'}!</Text>
        <Text style={styles.subtitle}>Dein heutiger Überblick</Text>

        {/* Calorie Progress */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Kalorien</Text>
          <View style={styles.calorieCircle}>
            <Text style={styles.calorieNumber}>{caloriesConsumed}</Text>
            <Text style={styles.calorieLabel}>von {calorieTarget} kcal</Text>
            <Text style={styles.caloriePercent}>{caloriePercent}%</Text>
          </View>
          <View style={styles.progressBarBackground}>
            <View
              style={[
                styles.progressBarFill,
                { width: `${Math.min(caloriePercent, 100)}%` },
              ]}
            />
          </View>
          <View style={styles.macroRow}>
            <View style={styles.macroItem}>
              <Text style={styles.macroValue}>{Math.round(proteinG)}g</Text>
              <Text style={styles.macroLabel}>Protein</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValue}>{Math.round(carbsG)}g</Text>
              <Text style={styles.macroLabel}>Kohlenhydrate</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValue}>{Math.round(fatG)}g</Text>
              <Text style={styles.macroLabel}>Fett</Text>
            </View>
          </View>
        </View>

        {/* Water Tracking */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Wasser</Text>
          <Text style={styles.waterText}>
            {waterTotal} / {waterTarget} ml ({waterPercent}%)
          </Text>
          <View style={styles.progressBarBackground}>
            <View
              style={[
                styles.progressBarFill,
                {
                  width: `${Math.min(waterPercent, 100)}%`,
                  backgroundColor: '#2196F3',
                },
              ]}
            />
          </View>
          <View style={styles.waterButtons}>
            <TouchableOpacity style={styles.waterButton} onPress={() => addWater(100)}>
              <Text style={styles.waterButtonText}>+100ml</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.waterButton} onPress={() => addWater(250)}>
              <Text style={styles.waterButtonText}>+250ml</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.waterButton} onPress={() => addWater(500)}>
              <Text style={styles.waterButtonText}>+500ml</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Today's Meals */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Heutige Mahlzeiten</Text>
          {['breakfast', 'lunch', 'dinner', 'snack'].map((mealType) => {
            const logs = mealGroups[mealType] || [];
            const totalCals = logs.reduce((s, l) => s + l.calories, 0);
            return (
              <View key={mealType} style={styles.mealRow}>
                <View style={{ flex: 1 }}>
                  <Text style={styles.mealType}>{MEAL_LABELS[mealType]}</Text>
                  <Text style={[styles.mealName, logs.length === 0 && { color: '#999' }]}>
                    {logs.length > 0
                      ? logs.map((l) => l.food_name).join(', ')
                      : 'Noch nicht erfasst'}
                  </Text>
                </View>
                <Text style={[styles.mealCalories, logs.length === 0 && { color: '#999' }]}>
                  {logs.length > 0 ? `${totalCals} kcal` : '-- kcal'}
                </Text>
              </View>
            );
          })}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  greeting: { fontSize: 28, fontWeight: 'bold', color: '#212121', marginBottom: 4 },
  subtitle: { fontSize: 16, color: '#757575', marginBottom: 20 },
  card: { backgroundColor: '#fff', borderRadius: 16, padding: 20, marginBottom: 16, shadowColor: '#000', shadowOffset: { width: 0, height: 2 }, shadowOpacity: 0.08, shadowRadius: 8, elevation: 2 },
  cardTitle: { fontSize: 18, fontWeight: '600', color: '#212121', marginBottom: 12 },
  calorieCircle: { alignItems: 'center', marginBottom: 16, paddingVertical: 20, borderWidth: 6, borderColor: PRIMARY, borderRadius: 100, width: 180, height: 180, justifyContent: 'center', alignSelf: 'center' },
  calorieNumber: { fontSize: 36, fontWeight: 'bold', color: PRIMARY },
  calorieLabel: { fontSize: 14, color: '#757575', marginTop: 2 },
  caloriePercent: { fontSize: 16, fontWeight: '600', color: PRIMARY, marginTop: 4 },
  progressBarBackground: { height: 8, backgroundColor: '#E0E0E0', borderRadius: 4, overflow: 'hidden', marginBottom: 12 },
  progressBarFill: { height: '100%', backgroundColor: PRIMARY, borderRadius: 4 },
  macroRow: { flexDirection: 'row', justifyContent: 'space-around' },
  macroItem: { alignItems: 'center' },
  macroValue: { fontSize: 18, fontWeight: '600', color: '#212121' },
  macroLabel: { fontSize: 12, color: '#757575', marginTop: 2 },
  waterText: { fontSize: 16, color: '#2196F3', marginBottom: 8, fontWeight: '500' },
  waterButtons: { flexDirection: 'row', justifyContent: 'space-around', marginTop: 8 },
  waterButton: { backgroundColor: '#E3F2FD', paddingHorizontal: 20, paddingVertical: 10, borderRadius: 20 },
  waterButtonText: { color: '#2196F3', fontWeight: '600', fontSize: 14 },
  mealRow: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', paddingVertical: 12, borderBottomWidth: 1, borderBottomColor: '#F0F0F0' },
  mealType: { fontSize: 12, color: '#757575', textTransform: 'uppercase', letterSpacing: 0.5 },
  mealName: { fontSize: 16, color: '#212121', marginTop: 2 },
  mealCalories: { fontSize: 16, fontWeight: '600', color: PRIMARY },
  offlineBanner: { flexDirection: 'row', alignItems: 'center', justifyContent: 'center', gap: 8, backgroundColor: '#FF9800', paddingVertical: 8, borderRadius: 8, marginBottom: 12 },
  offlineText: { color: '#fff', fontSize: 13, fontWeight: '600' },
});
