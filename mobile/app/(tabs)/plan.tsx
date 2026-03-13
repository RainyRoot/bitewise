import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useFocusEffect } from 'expo-router';
import { mealPlans } from '@/services/api';
import type { MealPlan, MealPlanEntry } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const DAYS = ['Mo', 'Di', 'Mi', 'Do', 'Fr', 'Sa', 'So'];
const MEAL_TYPE_LABELS: Record<string, string> = {
  breakfast: 'Frühstück',
  lunch: 'Mittagessen',
  dinner: 'Abendessen',
  snack: 'Snack',
};
const MEAL_TYPES = ['breakfast', 'lunch', 'dinner', 'snack'] as const;

export default function PlanScreen() {
  const [selectedDay, setSelectedDay] = useState(0); // 0=Mo, 6=So
  const [plan, setPlan] = useState<MealPlan | null>(null);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);

  const fetchPlan = async () => {
    try {
      const current = await mealPlans.getCurrent();
      setPlan(current);
    } catch {
      setPlan(null);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchPlan();
    }, [])
  );

  const handleGenerate = async () => {
    setGenerating(true);
    try {
      const newPlan = await mealPlans.generate();
      setPlan(newPlan);
    } catch (err: any) {
      Alert.alert('Fehler', err.message || 'Plan konnte nicht generiert werden.');
    } finally {
      setGenerating(false);
    }
  };

  // Get entries for selected day
  const dayEntries = (plan?.entries || []).filter((e) => e.day_of_week === selectedDay);
  const totalCalories = dayEntries.reduce(
    (sum, e) => sum + (e.recipe?.calories_per_serving || e.recipe?.calories || 0),
    0
  );

  // Calculate week dates from plan
  const weekDates = plan?.week_start_date
    ? (() => {
        const start = new Date(plan.week_start_date);
        const end = new Date(start);
        end.setDate(end.getDate() + 6);
        const fmt = (d: Date) => `${d.getDate()}. ${['Jan', 'Feb', 'Mär', 'Apr', 'Mai', 'Jun', 'Jul', 'Aug', 'Sep', 'Okt', 'Nov', 'Dez'][d.getMonth()]}`;
        return `${fmt(start)} - ${fmt(end)} ${end.getFullYear()}`;
      })()
    : '';

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
        <View style={styles.weekHeader}>
          <Text style={styles.weekTitle}>
            {plan ? 'Wochenplan' : 'Kein Plan vorhanden'}
          </Text>
          {weekDates ? <Text style={styles.weekDates}>{weekDates}</Text> : null}
        </View>

        {/* Day Selector */}
        <View style={styles.daySelector}>
          {DAYS.map((day, idx) => (
            <TouchableOpacity
              key={day}
              style={[styles.dayButton, selectedDay === idx && styles.dayButtonActive]}
              onPress={() => setSelectedDay(idx)}
            >
              <Text style={[styles.dayButtonText, selectedDay === idx && styles.dayButtonTextActive]}>
                {day}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {totalCalories > 0 && (
          <View style={styles.daySummary}>
            <Ionicons name="flame-outline" size={18} color={PRIMARY} />
            <Text style={styles.daySummaryText}>{totalCalories} kcal geplant</Text>
          </View>
        )}

        {/* Meals */}
        {MEAL_TYPES.map((mealType) => {
          const entry = dayEntries.find((e) => e.meal_type === mealType);
          return (
            <View key={mealType} style={styles.mealCard}>
              <View style={styles.mealHeader}>
                <Text style={styles.mealType}>{MEAL_TYPE_LABELS[mealType]}</Text>
                {entry?.is_locked && (
                  <Ionicons name="lock-closed" size={16} color="#757575" />
                )}
              </View>
              {entry?.recipe ? (
                <View>
                  <Text style={styles.mealName}>{entry.recipe.title}</Text>
                  <View style={styles.mealMeta}>
                    <View style={styles.mealMetaItem}>
                      <Ionicons name="flame-outline" size={14} color="#757575" />
                      <Text style={styles.mealMetaText}>
                        {entry.recipe.calories_per_serving || entry.recipe.calories || 0} kcal
                      </Text>
                    </View>
                    <View style={styles.mealMetaItem}>
                      <Ionicons name="time-outline" size={14} color="#757575" />
                      <Text style={styles.mealMetaText}>{entry.recipe.prep_time_min} Min.</Text>
                    </View>
                  </View>
                </View>
              ) : (
                <Text style={styles.emptyMeal}>Kein Rezept geplant</Text>
              )}
            </View>
          );
        })}

        <TouchableOpacity
          style={styles.generateButton}
          onPress={handleGenerate}
          disabled={generating}
        >
          {generating ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <>
              <Ionicons name="sparkles" size={20} color="#fff" />
              <Text style={styles.generateButtonText}>
                {plan ? 'Neuen Plan generieren' : 'Plan generieren'}
              </Text>
            </>
          )}
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  weekHeader: { marginBottom: 16 },
  weekTitle: { fontSize: 24, fontWeight: 'bold', color: '#212121' },
  weekDates: { fontSize: 14, color: '#757575', marginTop: 2 },
  daySelector: { flexDirection: 'row', justifyContent: 'space-between', marginBottom: 16 },
  dayButton: { width: 44, height: 44, borderRadius: 22, backgroundColor: '#fff', justifyContent: 'center', alignItems: 'center', shadowColor: '#000', shadowOffset: { width: 0, height: 1 }, shadowOpacity: 0.05, shadowRadius: 4, elevation: 1 },
  dayButtonActive: { backgroundColor: PRIMARY },
  dayButtonText: { fontSize: 14, fontWeight: '600', color: '#212121' },
  dayButtonTextActive: { color: '#fff' },
  daySummary: { flexDirection: 'row', alignItems: 'center', marginBottom: 12, gap: 6 },
  daySummaryText: { fontSize: 14, color: PRIMARY, fontWeight: '500' },
  mealCard: { backgroundColor: '#fff', borderRadius: 12, padding: 16, marginBottom: 12, shadowColor: '#000', shadowOffset: { width: 0, height: 2 }, shadowOpacity: 0.06, shadowRadius: 6, elevation: 1 },
  mealHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 },
  mealType: { fontSize: 12, fontWeight: '600', color: '#757575', textTransform: 'uppercase', letterSpacing: 0.5 },
  mealName: { fontSize: 16, fontWeight: '500', color: '#212121', marginBottom: 8 },
  mealMeta: { flexDirection: 'row', gap: 16 },
  mealMetaItem: { flexDirection: 'row', alignItems: 'center', gap: 4 },
  mealMetaText: { fontSize: 13, color: '#757575' },
  emptyMeal: { fontSize: 14, color: '#BDBDBD', fontStyle: 'italic' },
  generateButton: { flexDirection: 'row', backgroundColor: PRIMARY, borderRadius: 12, paddingVertical: 16, justifyContent: 'center', alignItems: 'center', marginTop: 8, gap: 8 },
  generateButtonText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
