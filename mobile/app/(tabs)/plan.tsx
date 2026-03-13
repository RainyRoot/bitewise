import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const DAYS = ['Mo', 'Di', 'Mi', 'Do', 'Fr', 'Sa', 'So'];
const MEAL_TYPES = ['Frühstück', 'Mittagessen', 'Abendessen', 'Snack'] as const;

// TODO: Replace with actual meal plan data from API
interface PlaceholderMeal {
  name: string;
  calories: number;
  prepTime: number;
}

const PLACEHOLDER_PLAN: Record<string, Record<string, PlaceholderMeal | null>> = {
  Mo: {
    'Frühstück': { name: 'Overnight Oats mit Banane', calories: 380, prepTime: 5 },
    'Mittagessen': { name: 'Quinoa-Gemüse-Bowl', calories: 520, prepTime: 25 },
    'Abendessen': { name: 'Lachs mit Brokkoli', calories: 480, prepTime: 20 },
    'Snack': { name: 'Griechischer Joghurt', calories: 150, prepTime: 2 },
  },
  Di: {
    'Frühstück': { name: 'Vollkorn-Toast mit Avocado', calories: 340, prepTime: 10 },
    'Mittagessen': { name: 'Hähnchen-Wrap', calories: 490, prepTime: 15 },
    'Abendessen': { name: 'Pasta Primavera', calories: 550, prepTime: 30 },
    'Snack': { name: 'Nussmischung', calories: 200, prepTime: 1 },
  },
  Mi: {
    'Frühstück': null,
    'Mittagessen': null,
    'Abendessen': null,
    'Snack': null,
  },
  Do: {
    'Frühstück': null,
    'Mittagessen': null,
    'Abendessen': null,
    'Snack': null,
  },
  Fr: {
    'Frühstück': null,
    'Mittagessen': null,
    'Abendessen': null,
    'Snack': null,
  },
  Sa: {
    'Frühstück': null,
    'Mittagessen': null,
    'Abendessen': null,
    'Snack': null,
  },
  So: {
    'Frühstück': null,
    'Mittagessen': null,
    'Abendessen': null,
    'Snack': null,
  },
};

export default function PlanScreen() {
  const [selectedDay, setSelectedDay] = useState('Mo');

  const dayMeals = PLACEHOLDER_PLAN[selectedDay] || {};
  const totalCalories = MEAL_TYPES.reduce((sum, type) => {
    return sum + (dayMeals[type]?.calories || 0);
  }, 0);

  const handleGenerate = () => {
    // TODO: Call api.mealPlans.generate()
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Week Header */}
        <View style={styles.weekHeader}>
          <Text style={styles.weekTitle}>Diese Woche</Text>
          <Text style={styles.weekDates}>10. - 16. März 2026</Text>
        </View>

        {/* Day Selector */}
        <View style={styles.daySelector}>
          {DAYS.map((day) => (
            <TouchableOpacity
              key={day}
              style={[
                styles.dayButton,
                selectedDay === day && styles.dayButtonActive,
              ]}
              onPress={() => setSelectedDay(day)}
            >
              <Text
                style={[
                  styles.dayButtonText,
                  selectedDay === day && styles.dayButtonTextActive,
                ]}
              >
                {day}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {/* Day Calorie Summary */}
        {totalCalories > 0 && (
          <View style={styles.daySummary}>
            <Ionicons name="flame-outline" size={18} color={PRIMARY} />
            <Text style={styles.daySummaryText}>
              {totalCalories} kcal geplant
            </Text>
          </View>
        )}

        {/* Meals */}
        {MEAL_TYPES.map((mealType) => {
          const meal = dayMeals[mealType];
          return (
            <View key={mealType} style={styles.mealCard}>
              <View style={styles.mealHeader}>
                <Text style={styles.mealType}>{mealType}</Text>
                {meal && (
                  <TouchableOpacity>
                    <Ionicons name="swap-horizontal" size={20} color="#757575" />
                  </TouchableOpacity>
                )}
              </View>
              {meal ? (
                <View>
                  <Text style={styles.mealName}>{meal.name}</Text>
                  <View style={styles.mealMeta}>
                    <View style={styles.mealMetaItem}>
                      <Ionicons name="flame-outline" size={14} color="#757575" />
                      <Text style={styles.mealMetaText}>{meal.calories} kcal</Text>
                    </View>
                    <View style={styles.mealMetaItem}>
                      <Ionicons name="time-outline" size={14} color="#757575" />
                      <Text style={styles.mealMetaText}>{meal.prepTime} Min.</Text>
                    </View>
                  </View>
                </View>
              ) : (
                <Text style={styles.emptyMeal}>Kein Rezept geplant</Text>
              )}
            </View>
          );
        })}

        {/* Generate Button */}
        <TouchableOpacity style={styles.generateButton} onPress={handleGenerate}>
          <Ionicons name="sparkles" size={20} color="#fff" />
          <Text style={styles.generateButtonText}>Plan generieren</Text>
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: BACKGROUND,
  },
  scrollContent: {
    padding: 16,
    paddingBottom: 32,
  },
  weekHeader: {
    marginBottom: 16,
  },
  weekTitle: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#212121',
  },
  weekDates: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  daySelector: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 16,
  },
  dayButton: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: '#fff',
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 1,
  },
  dayButtonActive: {
    backgroundColor: PRIMARY,
  },
  dayButtonText: {
    fontSize: 14,
    fontWeight: '600',
    color: '#212121',
  },
  dayButtonTextActive: {
    color: '#fff',
  },
  daySummary: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 12,
    gap: 6,
  },
  daySummaryText: {
    fontSize: 14,
    color: PRIMARY,
    fontWeight: '500',
  },
  mealCard: {
    backgroundColor: '#fff',
    borderRadius: 12,
    padding: 16,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 6,
    elevation: 1,
  },
  mealHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 8,
  },
  mealType: {
    fontSize: 12,
    fontWeight: '600',
    color: '#757575',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  mealName: {
    fontSize: 16,
    fontWeight: '500',
    color: '#212121',
    marginBottom: 8,
  },
  mealMeta: {
    flexDirection: 'row',
    gap: 16,
  },
  mealMetaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  mealMetaText: {
    fontSize: 13,
    color: '#757575',
  },
  emptyMeal: {
    fontSize: 14,
    color: '#BDBDBD',
    fontStyle: 'italic',
  },
  generateButton: {
    flexDirection: 'row',
    backgroundColor: PRIMARY,
    borderRadius: 12,
    paddingVertical: 16,
    justifyContent: 'center',
    alignItems: 'center',
    marginTop: 8,
    gap: 8,
  },
  generateButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
