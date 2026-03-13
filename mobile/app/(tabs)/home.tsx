import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
} from 'react-native';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

// TODO: Replace with actual user data from API
const PLACEHOLDER_USER = 'Max';
const CALORIE_TARGET = 2200;
const CALORIES_CONSUMED = 1450;

// TODO: Replace with actual water data from API
const WATER_TARGET = 2500;

interface MealSummary {
  type: string;
  name: string;
  calories: number;
}

// TODO: Replace with actual meal data from API
const PLACEHOLDER_MEALS: MealSummary[] = [
  { type: 'Frühstück', name: 'Haferflocken mit Beeren', calories: 350 },
  { type: 'Mittagessen', name: 'Hähnchensalat', calories: 520 },
  { type: 'Snack', name: 'Apfel & Mandeln', calories: 180 },
];

export default function HomeScreen() {
  // TODO: Fetch water logs from API and calculate total
  const [waterConsumed, setWaterConsumed] = useState(1200);

  const caloriePercent = Math.round((CALORIES_CONSUMED / CALORIE_TARGET) * 100);
  const waterPercent = Math.round((waterConsumed / WATER_TARGET) * 100);

  const addWater = (amount: number) => {
    // TODO: Call api.water.logWater({ amount_ml: amount })
    setWaterConsumed((prev) => Math.min(prev + amount, WATER_TARGET + 1000));
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Greeting */}
        <Text style={styles.greeting}>Hallo, {PLACEHOLDER_USER}!</Text>
        <Text style={styles.subtitle}>Dein heutiger Überblick</Text>

        {/* Calorie Progress */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Kalorien</Text>
          <View style={styles.calorieCircle}>
            <Text style={styles.calorieNumber}>{CALORIES_CONSUMED}</Text>
            <Text style={styles.calorieLabel}>von {CALORIE_TARGET} kcal</Text>
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
              <Text style={styles.macroValue}>85g</Text>
              <Text style={styles.macroLabel}>Protein</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValue}>165g</Text>
              <Text style={styles.macroLabel}>Kohlenhydrate</Text>
            </View>
            <View style={styles.macroItem}>
              <Text style={styles.macroValue}>48g</Text>
              <Text style={styles.macroLabel}>Fett</Text>
            </View>
          </View>
        </View>

        {/* Water Tracking */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Wasser</Text>
          <Text style={styles.waterText}>
            {waterConsumed} / {WATER_TARGET} ml ({waterPercent}%)
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
            <TouchableOpacity
              style={styles.waterButton}
              onPress={() => addWater(100)}
            >
              <Text style={styles.waterButtonText}>+100ml</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={styles.waterButton}
              onPress={() => addWater(250)}
            >
              <Text style={styles.waterButtonText}>+250ml</Text>
            </TouchableOpacity>
            <TouchableOpacity
              style={styles.waterButton}
              onPress={() => addWater(500)}
            >
              <Text style={styles.waterButtonText}>+500ml</Text>
            </TouchableOpacity>
          </View>
        </View>

        {/* Today's Meals */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Heutige Mahlzeiten</Text>
          {PLACEHOLDER_MEALS.map((meal, index) => (
            <View key={index} style={styles.mealRow}>
              <View>
                <Text style={styles.mealType}>{meal.type}</Text>
                <Text style={styles.mealName}>{meal.name}</Text>
              </View>
              <Text style={styles.mealCalories}>{meal.calories} kcal</Text>
            </View>
          ))}
          {/* TODO: Add "Abendessen" entry when logged */}
          <View style={styles.mealRow}>
            <View>
              <Text style={styles.mealType}>Abendessen</Text>
              <Text style={[styles.mealName, { color: '#999' }]}>
                Noch nicht erfasst
              </Text>
            </View>
            <Text style={[styles.mealCalories, { color: '#999' }]}>-- kcal</Text>
          </View>
        </View>
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
  greeting: {
    fontSize: 28,
    fontWeight: 'bold',
    color: '#212121',
    marginBottom: 4,
  },
  subtitle: {
    fontSize: 16,
    color: '#757575',
    marginBottom: 20,
  },
  card: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 20,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
  },
  cardTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#212121',
    marginBottom: 12,
  },
  calorieCircle: {
    alignItems: 'center',
    marginBottom: 16,
    paddingVertical: 20,
    borderWidth: 6,
    borderColor: PRIMARY,
    borderRadius: 100,
    width: 180,
    height: 180,
    justifyContent: 'center',
    alignSelf: 'center',
  },
  calorieNumber: {
    fontSize: 36,
    fontWeight: 'bold',
    color: PRIMARY,
  },
  calorieLabel: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  caloriePercent: {
    fontSize: 16,
    fontWeight: '600',
    color: PRIMARY,
    marginTop: 4,
  },
  progressBarBackground: {
    height: 8,
    backgroundColor: '#E0E0E0',
    borderRadius: 4,
    overflow: 'hidden',
    marginBottom: 12,
  },
  progressBarFill: {
    height: '100%',
    backgroundColor: PRIMARY,
    borderRadius: 4,
  },
  macroRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  macroItem: {
    alignItems: 'center',
  },
  macroValue: {
    fontSize: 18,
    fontWeight: '600',
    color: '#212121',
  },
  macroLabel: {
    fontSize: 12,
    color: '#757575',
    marginTop: 2,
  },
  waterText: {
    fontSize: 16,
    color: '#2196F3',
    marginBottom: 8,
    fontWeight: '500',
  },
  waterButtons: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    marginTop: 8,
  },
  waterButton: {
    backgroundColor: '#E3F2FD',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 20,
  },
  waterButtonText: {
    color: '#2196F3',
    fontWeight: '600',
    fontSize: 14,
  },
  mealRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#F0F0F0',
  },
  mealType: {
    fontSize: 12,
    color: '#757575',
    textTransform: 'uppercase',
    letterSpacing: 0.5,
  },
  mealName: {
    fontSize: 16,
    color: '#212121',
    marginTop: 2,
  },
  mealCalories: {
    fontSize: 16,
    fontWeight: '600',
    color: PRIMARY,
  },
});
