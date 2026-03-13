import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  TextInput,
  SafeAreaView,
  Image,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter } from 'expo-router';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const CATEGORIES = ['Alle', 'Frühstück', 'Mittag', 'Abend', 'Snack', 'Vegan', 'Low-Carb'];

// TODO: Replace with actual data from API
interface PlaceholderRecipe {
  id: number;
  title: string;
  calories: number;
  prepTime: number;
  category: string;
  isFavorite: boolean;
  imageUrl?: string;
}

const PLACEHOLDER_RECIPES: PlaceholderRecipe[] = [
  { id: 1, title: 'Overnight Oats mit Beeren', calories: 380, prepTime: 5, category: 'Frühstück', isFavorite: true },
  { id: 2, title: 'Quinoa-Gemüse-Bowl', calories: 520, prepTime: 25, category: 'Mittag', isFavorite: false },
  { id: 3, title: 'Lachs mit Brokkoli & Reis', calories: 480, prepTime: 20, category: 'Abend', isFavorite: true },
  { id: 4, title: 'Griechischer Joghurt mit Honig', calories: 150, prepTime: 2, category: 'Snack', isFavorite: false },
  { id: 5, title: 'Pasta Primavera', calories: 550, prepTime: 30, category: 'Abend', isFavorite: false },
  { id: 6, title: 'Avocado-Toast mit Ei', calories: 340, prepTime: 10, category: 'Frühstück', isFavorite: true },
  { id: 7, title: 'Buddha Bowl', calories: 450, prepTime: 20, category: 'Vegan', isFavorite: false },
  { id: 8, title: 'Hähnchen-Wrap', calories: 490, prepTime: 15, category: 'Mittag', isFavorite: false },
];

export default function RecipesScreen() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('Alle');

  const filteredRecipes = PLACEHOLDER_RECIPES.filter((recipe) => {
    const matchesSearch = recipe.title.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesCategory = selectedCategory === 'Alle' || recipe.category === selectedCategory;
    return matchesSearch && matchesCategory;
  });

  const toggleFavorite = (id: number) => {
    // TODO: Call api.recipes.addFavorite(id) or api.recipes.removeFavorite(id)
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Search */}
        <View style={styles.searchContainer}>
          <Ionicons name="search-outline" size={20} color="#757575" style={styles.searchIcon} />
          <TextInput
            style={styles.searchInput}
            placeholder="Rezepte suchen..."
            value={searchQuery}
            onChangeText={setSearchQuery}
          />
        </View>

        {/* Categories */}
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          style={styles.categoryScroll}
          contentContainerStyle={styles.categoryContent}
        >
          {CATEGORIES.map((cat) => (
            <TouchableOpacity
              key={cat}
              style={[
                styles.categoryChip,
                selectedCategory === cat && styles.categoryChipActive,
              ]}
              onPress={() => setSelectedCategory(cat)}
            >
              <Text
                style={[
                  styles.categoryText,
                  selectedCategory === cat && styles.categoryTextActive,
                ]}
              >
                {cat}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>

        {/* Results Count */}
        <Text style={styles.resultCount}>{filteredRecipes.length} Rezepte</Text>

        {/* Recipe Cards */}
        {filteredRecipes.map((recipe) => (
          <TouchableOpacity
            key={recipe.id}
            style={styles.recipeCard}
            onPress={() => router.push(`/recipe/${recipe.id}`)}
          >
            <View style={styles.recipePlaceholderImage}>
              <Ionicons name="restaurant-outline" size={32} color="#BDBDBD" />
            </View>
            <View style={styles.recipeInfo}>
              <Text style={styles.recipeTitle}>{recipe.title}</Text>
              <View style={styles.recipeMeta}>
                <View style={styles.recipeMetaItem}>
                  <Ionicons name="flame-outline" size={14} color="#757575" />
                  <Text style={styles.recipeMetaText}>{recipe.calories} kcal</Text>
                </View>
                <View style={styles.recipeMetaItem}>
                  <Ionicons name="time-outline" size={14} color="#757575" />
                  <Text style={styles.recipeMetaText}>{recipe.prepTime} Min.</Text>
                </View>
              </View>
              <View style={styles.recipeTags}>
                <View style={styles.tag}>
                  <Text style={styles.tagText}>{recipe.category}</Text>
                </View>
              </View>
            </View>
            <TouchableOpacity
              style={styles.favoriteButton}
              onPress={() => toggleFavorite(recipe.id)}
            >
              <Ionicons
                name={recipe.isFavorite ? 'heart' : 'heart-outline'}
                size={22}
                color={recipe.isFavorite ? '#E91E63' : '#BDBDBD'}
              />
            </TouchableOpacity>
          </TouchableOpacity>
        ))}
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
  searchContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 12,
    paddingHorizontal: 12,
    marginBottom: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 1,
  },
  searchIcon: {
    marginRight: 8,
  },
  searchInput: {
    flex: 1,
    paddingVertical: 14,
    fontSize: 15,
  },
  categoryScroll: {
    marginBottom: 12,
  },
  categoryContent: {
    gap: 8,
  },
  categoryChip: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    borderRadius: 20,
    backgroundColor: '#fff',
  },
  categoryChipActive: {
    backgroundColor: PRIMARY,
  },
  categoryText: {
    fontSize: 13,
    fontWeight: '600',
    color: '#757575',
  },
  categoryTextActive: {
    color: '#fff',
  },
  resultCount: {
    fontSize: 13,
    color: '#757575',
    marginBottom: 12,
  },
  recipeCard: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    borderRadius: 12,
    marginBottom: 12,
    overflow: 'hidden',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 6,
    elevation: 1,
  },
  recipePlaceholderImage: {
    width: 90,
    backgroundColor: '#F5F5F5',
    justifyContent: 'center',
    alignItems: 'center',
  },
  recipeInfo: {
    flex: 1,
    padding: 12,
  },
  recipeTitle: {
    fontSize: 15,
    fontWeight: '600',
    color: '#212121',
    marginBottom: 6,
  },
  recipeMeta: {
    flexDirection: 'row',
    gap: 12,
    marginBottom: 6,
  },
  recipeMetaItem: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  recipeMetaText: {
    fontSize: 12,
    color: '#757575',
  },
  recipeTags: {
    flexDirection: 'row',
    gap: 6,
  },
  tag: {
    backgroundColor: '#E8F5E9',
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 4,
  },
  tagText: {
    fontSize: 11,
    color: PRIMARY,
    fontWeight: '500',
  },
  favoriteButton: {
    padding: 12,
    justifyContent: 'center',
  },
});
