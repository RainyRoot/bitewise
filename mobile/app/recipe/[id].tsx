import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
} from 'react-native';
import { useLocalSearchParams, useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

// TODO: Replace with actual API call using recipe ID
const PLACEHOLDER_RECIPE = {
  id: 1,
  title: 'Overnight Oats mit Beeren',
  description: 'Gesunde und einfache Overnight Oats mit frischen Beeren, Chiasamen und einem Hauch Honig. Perfekt als schnelles Frühstück.',
  calories: 380,
  protein_g: 12,
  carbs_g: 55,
  fat_g: 8,
  fiber_g: 7,
  prepTime: 5,
  cookTime: 0,
  servings: 1,
  difficulty: 'Einfach',
  source: 'chefkoch.de',
  isFavorite: true,
  ingredients: [
    { name: 'Haferflocken', quantity: '50g' },
    { name: 'Milch / Pflanzenmilch', quantity: '150ml' },
    { name: 'Griechischer Joghurt', quantity: '80g' },
    { name: 'Chiasamen', quantity: '1 EL' },
    { name: 'Honig', quantity: '1 TL' },
    { name: 'Frische Beeren', quantity: '80g' },
    { name: 'Nüsse (optional)', quantity: '15g' },
  ],
  instructions: [
    'Haferflocken, Milch, Joghurt und Chiasamen in einem Glas verrühren.',
    'Mit Honig süßen und gut durchmischen.',
    'Über Nacht (mind. 4 Stunden) im Kühlschrank ziehen lassen.',
    'Am nächsten Morgen mit frischen Beeren toppen.',
    'Optional mit gehackten Nüssen garnieren.',
  ],
};

export default function RecipeDetailScreen() {
  const { id } = useLocalSearchParams();
  const router = useRouter();
  const [isFavorite, setIsFavorite] = useState(PLACEHOLDER_RECIPE.isFavorite);

  const recipe = PLACEHOLDER_RECIPE; // TODO: Fetch from API using id

  const toggleFavorite = () => {
    // TODO: Call api.recipes.addFavorite(id) or api.recipes.removeFavorite(id)
    setIsFavorite(!isFavorite);
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Hero / Image Placeholder */}
        <View style={styles.heroImage}>
          <Ionicons name="restaurant-outline" size={64} color="#BDBDBD" />
        </View>

        {/* Title & Actions */}
        <View style={styles.titleRow}>
          <Text style={styles.title}>{recipe.title}</Text>
          <TouchableOpacity onPress={toggleFavorite}>
            <Ionicons
              name={isFavorite ? 'heart' : 'heart-outline'}
              size={28}
              color={isFavorite ? '#E91E63' : '#BDBDBD'}
            />
          </TouchableOpacity>
        </View>

        <Text style={styles.description}>{recipe.description}</Text>

        {/* Quick Info */}
        <View style={styles.infoRow}>
          <View style={styles.infoItem}>
            <Ionicons name="time-outline" size={20} color={PRIMARY} />
            <Text style={styles.infoValue}>{recipe.prepTime} Min.</Text>
            <Text style={styles.infoLabel}>Zubereitung</Text>
          </View>
          <View style={styles.infoItem}>
            <Ionicons name="people-outline" size={20} color={PRIMARY} />
            <Text style={styles.infoValue}>{recipe.servings}</Text>
            <Text style={styles.infoLabel}>Portionen</Text>
          </View>
          <View style={styles.infoItem}>
            <Ionicons name="speedometer-outline" size={20} color={PRIMARY} />
            <Text style={styles.infoValue}>{recipe.difficulty}</Text>
            <Text style={styles.infoLabel}>Schwierigkeit</Text>
          </View>
        </View>

        {/* Nutrition */}
        <View style={styles.card}>
          <Text style={styles.sectionTitle}>Nährwerte pro Portion</Text>
          <View style={styles.nutritionGrid}>
            <View style={styles.nutritionItem}>
              <Text style={styles.nutritionValue}>{recipe.calories}</Text>
              <Text style={styles.nutritionLabel}>kcal</Text>
            </View>
            <View style={styles.nutritionItem}>
              <Text style={styles.nutritionValue}>{recipe.protein_g}g</Text>
              <Text style={styles.nutritionLabel}>Protein</Text>
            </View>
            <View style={styles.nutritionItem}>
              <Text style={styles.nutritionValue}>{recipe.carbs_g}g</Text>
              <Text style={styles.nutritionLabel}>Kohlenhydrate</Text>
            </View>
            <View style={styles.nutritionItem}>
              <Text style={styles.nutritionValue}>{recipe.fat_g}g</Text>
              <Text style={styles.nutritionLabel}>Fett</Text>
            </View>
          </View>
        </View>

        {/* Ingredients */}
        <View style={styles.card}>
          <Text style={styles.sectionTitle}>Zutaten</Text>
          {recipe.ingredients.map((ing, index) => (
            <View key={index} style={styles.ingredientRow}>
              <Text style={styles.ingredientQuantity}>{ing.quantity}</Text>
              <Text style={styles.ingredientName}>{ing.name}</Text>
            </View>
          ))}
        </View>

        {/* Instructions */}
        <View style={styles.card}>
          <Text style={styles.sectionTitle}>Zubereitung</Text>
          {recipe.instructions.map((step, index) => (
            <View key={index} style={styles.stepRow}>
              <View style={styles.stepNumber}>
                <Text style={styles.stepNumberText}>{index + 1}</Text>
              </View>
              <Text style={styles.stepText}>{step}</Text>
            </View>
          ))}
        </View>

        {/* Source */}
        <Text style={styles.source}>Quelle: {recipe.source}</Text>

        {/* Add to Meal Plan Button */}
        <TouchableOpacity style={styles.addButton}>
          <Ionicons name="calendar-outline" size={20} color="#fff" />
          <Text style={styles.addButtonText}>Zum Wochenplan hinzufügen</Text>
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
    paddingBottom: 32,
  },
  heroImage: {
    height: 200,
    backgroundColor: '#E0E0E0',
    justifyContent: 'center',
    alignItems: 'center',
  },
  titleRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    padding: 16,
    paddingBottom: 0,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#212121',
    flex: 1,
    marginRight: 12,
  },
  description: {
    fontSize: 14,
    color: '#757575',
    lineHeight: 20,
    paddingHorizontal: 16,
    marginTop: 8,
    marginBottom: 16,
  },
  infoRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
    backgroundColor: '#fff',
    marginHorizontal: 16,
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 6,
    elevation: 1,
  },
  infoItem: {
    alignItems: 'center',
  },
  infoValue: {
    fontSize: 16,
    fontWeight: '600',
    color: '#212121',
    marginTop: 4,
  },
  infoLabel: {
    fontSize: 11,
    color: '#757575',
    marginTop: 2,
  },
  card: {
    backgroundColor: '#fff',
    marginHorizontal: 16,
    borderRadius: 12,
    padding: 16,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 6,
    elevation: 1,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    color: '#212121',
    marginBottom: 12,
  },
  nutritionGrid: {
    flexDirection: 'row',
    justifyContent: 'space-around',
  },
  nutritionItem: {
    alignItems: 'center',
  },
  nutritionValue: {
    fontSize: 20,
    fontWeight: 'bold',
    color: PRIMARY,
  },
  nutritionLabel: {
    fontSize: 12,
    color: '#757575',
    marginTop: 2,
  },
  ingredientRow: {
    flexDirection: 'row',
    paddingVertical: 8,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
  },
  ingredientQuantity: {
    fontSize: 14,
    fontWeight: '600',
    color: PRIMARY,
    width: 80,
  },
  ingredientName: {
    fontSize: 14,
    color: '#212121',
    flex: 1,
  },
  stepRow: {
    flexDirection: 'row',
    marginBottom: 12,
  },
  stepNumber: {
    width: 28,
    height: 28,
    borderRadius: 14,
    backgroundColor: PRIMARY,
    justifyContent: 'center',
    alignItems: 'center',
    marginRight: 12,
    flexShrink: 0,
  },
  stepNumberText: {
    color: '#fff',
    fontWeight: 'bold',
    fontSize: 14,
  },
  stepText: {
    fontSize: 14,
    color: '#424242',
    lineHeight: 20,
    flex: 1,
    paddingTop: 4,
  },
  source: {
    fontSize: 12,
    color: '#BDBDBD',
    textAlign: 'center',
    marginBottom: 16,
  },
  addButton: {
    flexDirection: 'row',
    backgroundColor: PRIMARY,
    marginHorizontal: 16,
    borderRadius: 12,
    paddingVertical: 16,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
  },
  addButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
