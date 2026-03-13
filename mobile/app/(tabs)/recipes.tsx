import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  TextInput,
  SafeAreaView,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useRouter, useFocusEffect } from 'expo-router';
import { recipes } from '@/services/api';
import type { Recipe } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const CATEGORIES = ['Alle', 'Frühstück', 'Mittag', 'Abend', 'Snack', 'Vegan', 'Low-Carb'];

export default function RecipesScreen() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('Alle');
  const [recipeList, setRecipeList] = useState<Recipe[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);

  const fetchRecipes = async () => {
    try {
      const filter: any = { limit: 50 };
      if (searchQuery.trim()) filter.query = searchQuery.trim();
      if (selectedCategory !== 'Alle') filter.category = selectedCategory;
      const result = await recipes.list(filter);
      setRecipeList(result.recipes || []);
      setTotal(result.total || 0);
    } catch {
      setRecipeList([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchRecipes();
    }, [searchQuery, selectedCategory])
  );

  const toggleFavorite = async (id: number, isFav: boolean) => {
    try {
      if (isFav) {
        await recipes.removeFavorite(id);
      } else {
        await recipes.addFavorite(id);
      }
      setRecipeList((prev) =>
        prev.map((r) => (r.id === id ? { ...r, is_favorite: !isFav } : r))
      );
    } catch {
      // ignore
    }
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
            returnKeyType="search"
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
              style={[styles.categoryChip, selectedCategory === cat && styles.categoryChipActive]}
              onPress={() => setSelectedCategory(cat)}
            >
              <Text style={[styles.categoryText, selectedCategory === cat && styles.categoryTextActive]}>
                {cat}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>

        {/* Create Recipe Button */}
        <TouchableOpacity
          style={styles.createButton}
          onPress={() => router.push('/create-recipe')}
        >
          <Ionicons name="add-circle-outline" size={20} color="#fff" />
          <Text style={styles.createButtonText}>Rezept erstellen</Text>
        </TouchableOpacity>

        <Text style={styles.resultCount}>{total} Rezepte</Text>

        {loading ? (
          <ActivityIndicator size="large" color={PRIMARY} style={{ marginTop: 32 }} />
        ) : (
          recipeList.map((recipe) => (
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
                    <Text style={styles.recipeMetaText}>
                      {recipe.calories_per_serving || recipe.calories || 0} kcal
                    </Text>
                  </View>
                  <View style={styles.recipeMetaItem}>
                    <Ionicons name="time-outline" size={14} color="#757575" />
                    <Text style={styles.recipeMetaText}>{recipe.prep_time_min} Min.</Text>
                  </View>
                </View>
                {(recipe.categories || []).length > 0 && (
                  <View style={styles.recipeTags}>
                    {(recipe.categories || []).slice(0, 2).map((cat) => (
                      <View key={cat} style={styles.tag}>
                        <Text style={styles.tagText}>{cat}</Text>
                      </View>
                    ))}
                  </View>
                )}
              </View>
              <TouchableOpacity
                style={styles.favoriteButton}
                onPress={() => toggleFavorite(recipe.id, !!recipe.is_favorite)}
              >
                <Ionicons
                  name={recipe.is_favorite ? 'heart' : 'heart-outline'}
                  size={22}
                  color={recipe.is_favorite ? '#E91E63' : '#BDBDBD'}
                />
              </TouchableOpacity>
            </TouchableOpacity>
          ))
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  searchContainer: { flexDirection: 'row', alignItems: 'center', backgroundColor: '#fff', borderRadius: 12, paddingHorizontal: 12, marginBottom: 12, shadowColor: '#000', shadowOffset: { width: 0, height: 1 }, shadowOpacity: 0.05, shadowRadius: 4, elevation: 1 },
  searchIcon: { marginRight: 8 },
  searchInput: { flex: 1, paddingVertical: 14, fontSize: 15 },
  categoryScroll: { marginBottom: 12 },
  categoryContent: { gap: 8 },
  categoryChip: { paddingHorizontal: 16, paddingVertical: 8, borderRadius: 20, backgroundColor: '#fff' },
  categoryChipActive: { backgroundColor: PRIMARY },
  categoryText: { fontSize: 13, fontWeight: '600', color: '#757575' },
  categoryTextActive: { color: '#fff' },
  resultCount: { fontSize: 13, color: '#757575', marginBottom: 12 },
  recipeCard: { flexDirection: 'row', backgroundColor: '#fff', borderRadius: 12, marginBottom: 12, overflow: 'hidden', shadowColor: '#000', shadowOffset: { width: 0, height: 2 }, shadowOpacity: 0.06, shadowRadius: 6, elevation: 1 },
  recipePlaceholderImage: { width: 90, backgroundColor: '#F5F5F5', justifyContent: 'center', alignItems: 'center' },
  recipeInfo: { flex: 1, padding: 12 },
  recipeTitle: { fontSize: 15, fontWeight: '600', color: '#212121', marginBottom: 6 },
  recipeMeta: { flexDirection: 'row', gap: 12, marginBottom: 6 },
  recipeMetaItem: { flexDirection: 'row', alignItems: 'center', gap: 4 },
  recipeMetaText: { fontSize: 12, color: '#757575' },
  recipeTags: { flexDirection: 'row', gap: 6 },
  tag: { backgroundColor: '#E8F5E9', paddingHorizontal: 8, paddingVertical: 3, borderRadius: 4 },
  tagText: { fontSize: 11, color: PRIMARY, fontWeight: '500' },
  favoriteButton: { padding: 12, justifyContent: 'center' },
  createButton: { flexDirection: 'row', backgroundColor: PRIMARY, borderRadius: 12, paddingVertical: 12, justifyContent: 'center', alignItems: 'center', gap: 8, marginBottom: 12 },
  createButtonText: { color: '#fff', fontSize: 15, fontWeight: '600' },
});
