import React, { useState, useEffect, useCallback } from 'react';
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
import { shoppingLists, mealPlans } from '@/services/api';
import type { ShoppingList, ShoppingListItem } from '@/types';
import { useFocusEffect } from 'expo-router';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

export default function ShoppingScreen() {
  const [list, setList] = useState<ShoppingList | null>(null);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);

  const fetchList = async () => {
    try {
      const current = await shoppingLists.getCurrent();
      setList(current);
    } catch {
      setList(null);
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchList();
    }, [])
  );

  const handleGenerate = async () => {
    setGenerating(true);
    try {
      const plan = await mealPlans.getCurrent();
      const newList = await shoppingLists.generate(plan.id);
      setList(newList);
    } catch (err: any) {
      Alert.alert('Fehler', err.message || 'Kein aktiver Essensplan gefunden.');
    } finally {
      setGenerating(false);
    }
  };

  const handleToggle = async (itemId: number) => {
    try {
      await shoppingLists.toggleItem(itemId);
      setList((prev) => {
        if (!prev) return prev;
        return {
          ...prev,
          items: prev.items.map((item) =>
            item.id === itemId ? { ...item, is_checked: !item.is_checked } : item
          ),
        };
      });
    } catch (err: any) {
      Alert.alert('Fehler', err.message);
    }
  };

  // Group items by category
  const groupedItems = (list?.items || []).reduce<Record<string, ShoppingListItem[]>>(
    (groups, item) => {
      const cat = item.category || 'Sonstiges';
      if (!groups[cat]) groups[cat] = [];
      groups[cat].push(item);
      return groups;
    },
    {}
  );

  const checkedCount = (list?.items || []).filter((i) => i.is_checked).length;
  const totalCount = (list?.items || []).length;

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
        {list && totalCount > 0 ? (
          <>
            <View style={styles.header}>
              <Text style={styles.title}>Einkaufsliste</Text>
              <Text style={styles.subtitle}>
                {checkedCount} von {totalCount} erledigt
              </Text>
            </View>

            <View style={styles.progressBarBackground}>
              <View
                style={[
                  styles.progressBarFill,
                  { width: `${totalCount > 0 ? (checkedCount / totalCount) * 100 : 0}%` },
                ]}
              />
            </View>

            {Object.entries(groupedItems).map(([category, items]) => (
              <View key={category} style={styles.categorySection}>
                <Text style={styles.categoryTitle}>{category}</Text>
                {items.map((item) => (
                  <TouchableOpacity
                    key={item.id}
                    style={styles.itemRow}
                    onPress={() => handleToggle(item.id)}
                  >
                    <Ionicons
                      name={item.is_checked ? 'checkbox' : 'square-outline'}
                      size={22}
                      color={item.is_checked ? PRIMARY : '#BDBDBD'}
                    />
                    <View style={styles.itemInfo}>
                      <Text
                        style={[
                          styles.itemName,
                          item.is_checked && styles.itemChecked,
                        ]}
                      >
                        {item.ingredient_name}
                      </Text>
                      {item.quantity > 0 && (
                        <Text style={styles.itemQuantity}>
                          {item.quantity} {item.unit}
                        </Text>
                      )}
                    </View>
                  </TouchableOpacity>
                ))}
              </View>
            ))}
          </>
        ) : (
          <View style={styles.emptyState}>
            <Ionicons name="cart-outline" size={64} color="#BDBDBD" />
            <Text style={styles.emptyTitle}>Keine Einkaufsliste</Text>
            <Text style={styles.emptySubtitle}>
              Erstelle eine Einkaufsliste aus deinem aktuellen Essensplan.
            </Text>
          </View>
        )}

        <TouchableOpacity
          style={styles.generateButton}
          onPress={handleGenerate}
          disabled={generating}
        >
          {generating ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <>
              <Ionicons name="refresh-outline" size={20} color="#fff" />
              <Text style={styles.generateButtonText}>
                {list ? 'Neue Liste generieren' : 'Aus Essensplan generieren'}
              </Text>
            </>
          )}
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
  header: {
    marginBottom: 12,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#212121',
  },
  subtitle: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  progressBarBackground: {
    height: 6,
    backgroundColor: '#E0E0E0',
    borderRadius: 3,
    overflow: 'hidden',
    marginBottom: 20,
  },
  progressBarFill: {
    height: '100%',
    backgroundColor: PRIMARY,
    borderRadius: 3,
  },
  categorySection: {
    marginBottom: 16,
  },
  categoryTitle: {
    fontSize: 14,
    fontWeight: '600',
    color: PRIMARY,
    textTransform: 'uppercase',
    letterSpacing: 0.5,
    marginBottom: 8,
  },
  itemRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 10,
    padding: 14,
    marginBottom: 6,
    gap: 12,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.04,
    shadowRadius: 3,
    elevation: 1,
  },
  itemInfo: {
    flex: 1,
  },
  itemName: {
    fontSize: 15,
    color: '#212121',
    fontWeight: '500',
  },
  itemChecked: {
    textDecorationLine: 'line-through',
    color: '#BDBDBD',
  },
  itemQuantity: {
    fontSize: 12,
    color: '#757575',
    marginTop: 2,
  },
  emptyState: {
    alignItems: 'center',
    paddingVertical: 48,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: '600',
    color: '#757575',
    marginTop: 16,
  },
  emptySubtitle: {
    fontSize: 14,
    color: '#BDBDBD',
    textAlign: 'center',
    marginTop: 8,
    paddingHorizontal: 32,
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
