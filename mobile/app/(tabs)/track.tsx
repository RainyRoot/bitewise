import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  TextInput,
  SafeAreaView,
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

type MealType = 'breakfast' | 'lunch' | 'dinner' | 'snack';

const MEAL_LABELS: Record<MealType, string> = {
  breakfast: 'Frühstück',
  lunch: 'Mittagessen',
  dinner: 'Abendessen',
  snack: 'Snack',
};

interface QuickFood {
  name: string;
  calories: number;
  protein_g: number;
  carbs_g: number;
  fat_g: number;
}

const QUICK_FOODS: QuickFood[] = [
  { name: 'Apfel', calories: 52, protein_g: 0.3, carbs_g: 14, fat_g: 0.2 },
  { name: 'Banane', calories: 89, protein_g: 1.1, carbs_g: 23, fat_g: 0.3 },
  { name: 'Vollkornbrot (1 Scheibe)', calories: 120, protein_g: 4, carbs_g: 22, fat_g: 1.5 },
  { name: 'Ei (gekocht)', calories: 78, protein_g: 6, carbs_g: 0.6, fat_g: 5 },
  { name: 'Joghurt (150g)', calories: 90, protein_g: 5, carbs_g: 12, fat_g: 2 },
  { name: 'Hähnchenbrust (100g)', calories: 165, protein_g: 31, carbs_g: 0, fat_g: 3.6 },
];

// TODO: Replace with actual data from API
const PLACEHOLDER_LOGS = {
  breakfast: [
    { id: 1, name: 'Haferflocken mit Beeren', calories: 350, protein_g: 12, carbs_g: 55, fat_g: 8 },
  ],
  lunch: [
    { id: 2, name: 'Hähnchensalat', calories: 520, protein_g: 35, carbs_g: 20, fat_g: 28 },
  ],
  dinner: [],
  snack: [
    { id: 3, name: 'Apfel & Mandeln', calories: 180, protein_g: 4, carbs_g: 22, fat_g: 9 },
  ],
};

export default function TrackScreen() {
  const [selectedMeal, setSelectedMeal] = useState<MealType>('lunch');
  const [showManualEntry, setShowManualEntry] = useState(false);
  const [foodName, setFoodName] = useState('');
  const [calories, setCalories] = useState('');

  const totalCalories = Object.values(PLACEHOLDER_LOGS).reduce(
    (sum, meals) => sum + meals.reduce((s, m) => s + m.calories, 0),
    0,
  );

  const handleQuickAdd = (food: QuickFood) => {
    // TODO: Call api.tracking.logFood({ name: food.name, meal_type: selectedMeal, ...food, amount: 1, unit: 'portion' })
    Alert.alert('Hinzugefügt', `${food.name} zu ${MEAL_LABELS[selectedMeal]} hinzugefügt`);
  };

  const handleManualAdd = () => {
    if (!foodName || !calories) return;
    // TODO: Call api.tracking.logFood(...)
    Alert.alert('Hinzugefügt', `${foodName} (${calories} kcal) erfasst`);
    setFoodName('');
    setCalories('');
    setShowManualEntry(false);
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Daily Summary */}
        <View style={styles.summaryCard}>
          <Text style={styles.summaryTitle}>Heute</Text>
          <Text style={styles.summaryCalories}>{totalCalories} kcal</Text>
          <Text style={styles.summarySubtext}>von 2.200 kcal Ziel</Text>
        </View>

        {/* Meal Type Selector */}
        <View style={styles.mealSelector}>
          {(Object.keys(MEAL_LABELS) as MealType[]).map((type) => (
            <TouchableOpacity
              key={type}
              style={[
                styles.mealTab,
                selectedMeal === type && styles.mealTabActive,
              ]}
              onPress={() => setSelectedMeal(type)}
            >
              <Text
                style={[
                  styles.mealTabText,
                  selectedMeal === type && styles.mealTabTextActive,
                ]}
              >
                {MEAL_LABELS[type]}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {/* Current Logs for Selected Meal */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>{MEAL_LABELS[selectedMeal]}</Text>
          {PLACEHOLDER_LOGS[selectedMeal].length > 0 ? (
            PLACEHOLDER_LOGS[selectedMeal].map((log) => (
              <View key={log.id} style={styles.logRow}>
                <View style={styles.logInfo}>
                  <Text style={styles.logName}>{log.name}</Text>
                  <Text style={styles.logMacros}>
                    P: {log.protein_g}g | K: {log.carbs_g}g | F: {log.fat_g}g
                  </Text>
                </View>
                <Text style={styles.logCalories}>{log.calories} kcal</Text>
              </View>
            ))
          ) : (
            <Text style={styles.emptyText}>Noch nichts erfasst</Text>
          )}
        </View>

        {/* Quick Add */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Schnell hinzufügen</Text>
          <View style={styles.quickGrid}>
            {QUICK_FOODS.map((food, index) => (
              <TouchableOpacity
                key={index}
                style={styles.quickItem}
                onPress={() => handleQuickAdd(food)}
              >
                <Text style={styles.quickName}>{food.name}</Text>
                <Text style={styles.quickCalories}>{food.calories} kcal</Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        {/* Manual Entry */}
        {showManualEntry ? (
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Manuell erfassen</Text>
            <TextInput
              style={styles.input}
              placeholder="Lebensmittel"
              value={foodName}
              onChangeText={setFoodName}
            />
            <TextInput
              style={styles.input}
              placeholder="Kalorien (kcal)"
              value={calories}
              onChangeText={setCalories}
              keyboardType="numeric"
            />
            <View style={styles.manualButtons}>
              <TouchableOpacity
                style={styles.cancelButton}
                onPress={() => setShowManualEntry(false)}
              >
                <Text style={styles.cancelButtonText}>Abbrechen</Text>
              </TouchableOpacity>
              <TouchableOpacity style={styles.addButton} onPress={handleManualAdd}>
                <Text style={styles.addButtonText}>Hinzufügen</Text>
              </TouchableOpacity>
            </View>
          </View>
        ) : (
          <View style={styles.actionButtons}>
            <TouchableOpacity
              style={styles.actionButton}
              onPress={() => setShowManualEntry(true)}
            >
              <Ionicons name="create-outline" size={20} color={PRIMARY} />
              <Text style={styles.actionButtonText}>Manuell</Text>
            </TouchableOpacity>
            <TouchableOpacity style={styles.actionButton}>
              <Ionicons name="barcode-outline" size={20} color={PRIMARY} />
              <Text style={styles.actionButtonText}>Scannen</Text>
            </TouchableOpacity>
          </View>
        )}
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
  summaryCard: {
    backgroundColor: PRIMARY,
    borderRadius: 16,
    padding: 20,
    alignItems: 'center',
    marginBottom: 16,
  },
  summaryTitle: {
    fontSize: 14,
    color: 'rgba(255,255,255,0.8)',
    fontWeight: '500',
  },
  summaryCalories: {
    fontSize: 36,
    fontWeight: 'bold',
    color: '#fff',
    marginTop: 4,
  },
  summarySubtext: {
    fontSize: 14,
    color: 'rgba(255,255,255,0.7)',
    marginTop: 2,
  },
  mealSelector: {
    flexDirection: 'row',
    marginBottom: 16,
    gap: 8,
  },
  mealTab: {
    flex: 1,
    paddingVertical: 10,
    borderRadius: 8,
    backgroundColor: '#fff',
    alignItems: 'center',
  },
  mealTabActive: {
    backgroundColor: PRIMARY,
  },
  mealTabText: {
    fontSize: 12,
    fontWeight: '600',
    color: '#757575',
  },
  mealTabTextActive: {
    color: '#fff',
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
  logRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#F0F0F0',
  },
  logInfo: {
    flex: 1,
  },
  logName: {
    fontSize: 15,
    color: '#212121',
    fontWeight: '500',
  },
  logMacros: {
    fontSize: 12,
    color: '#757575',
    marginTop: 2,
  },
  logCalories: {
    fontSize: 15,
    fontWeight: '600',
    color: PRIMARY,
  },
  emptyText: {
    fontSize: 14,
    color: '#BDBDBD',
    fontStyle: 'italic',
    textAlign: 'center',
    paddingVertical: 16,
  },
  quickGrid: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  quickItem: {
    backgroundColor: '#F5F5F5',
    borderRadius: 8,
    paddingHorizontal: 12,
    paddingVertical: 10,
    minWidth: '30%',
    flexGrow: 1,
  },
  quickName: {
    fontSize: 13,
    fontWeight: '500',
    color: '#212121',
  },
  quickCalories: {
    fontSize: 11,
    color: '#757575',
    marginTop: 2,
  },
  input: {
    backgroundColor: '#F5F5F5',
    borderRadius: 8,
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 15,
    marginBottom: 10,
  },
  manualButtons: {
    flexDirection: 'row',
    gap: 12,
    marginTop: 4,
  },
  cancelButton: {
    flex: 1,
    paddingVertical: 12,
    borderRadius: 8,
    backgroundColor: '#F5F5F5',
    alignItems: 'center',
  },
  cancelButtonText: {
    color: '#757575',
    fontWeight: '600',
  },
  addButton: {
    flex: 1,
    paddingVertical: 12,
    borderRadius: 8,
    backgroundColor: PRIMARY,
    alignItems: 'center',
  },
  addButtonText: {
    color: '#fff',
    fontWeight: '600',
  },
  actionButtons: {
    flexDirection: 'row',
    gap: 12,
  },
  actionButton: {
    flex: 1,
    flexDirection: 'row',
    backgroundColor: '#fff',
    borderRadius: 12,
    paddingVertical: 16,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 6,
    elevation: 1,
  },
  actionButtonText: {
    fontSize: 15,
    fontWeight: '600',
    color: PRIMARY,
  },
});
