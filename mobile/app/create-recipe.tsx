import React, { useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  SafeAreaView,
  ScrollView,
  Alert,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { recipes } from '@/services/api';
import type { SimpleIngredient } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

export default function CreateRecipeScreen() {
  const router = useRouter();
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [prepTime, setPrepTime] = useState('');
  const [cookTime, setCookTime] = useState('');
  const [servings, setServings] = useState('4');
  const [calories, setCalories] = useState('');
  const [proteinG, setProteinG] = useState('');
  const [carbsG, setCarbsG] = useState('');
  const [fatG, setFatG] = useState('');
  const [instructions, setInstructions] = useState('');
  const [ingredients, setIngredients] = useState<SimpleIngredient[]>([
    { name: '', quantity: 0, unit: '' },
  ]);
  const [saving, setSaving] = useState(false);

  const addIngredient = () => {
    setIngredients([
      ...ingredients,
      { name: '', quantity: 0, unit: '' },
    ]);
  };

  const updateIngredient = (index: number, field: string, value: string) => {
    const updated = [...ingredients];
    if (field === 'name' || field === 'unit') {
      (updated[index] as any)[field] = value;
    } else {
      (updated[index] as any)[field] = parseFloat(value) || 0;
    }
    setIngredients(updated);
  };

  const removeIngredient = (index: number) => {
    if (ingredients.length <= 1) return;
    setIngredients(ingredients.filter((_, i) => i !== index));
  };

  const handleSave = async () => {
    if (!title.trim()) {
      Alert.alert('Fehler', 'Bitte gib einen Titel ein.');
      return;
    }

    setSaving(true);
    try {
      await recipes.create({
        title: title.trim(),
        description: description.trim(),
        prep_time_min: parseInt(prepTime, 10) || 0,
        cook_time_min: parseInt(cookTime, 10) || 0,
        servings: parseInt(servings, 10) || 4,
        calories_per_serving: parseInt(calories, 10) || 0,
        protein_g: parseFloat(proteinG) || 0,
        carbs_g: parseFloat(carbsG) || 0,
        fat_g: parseFloat(fatG) || 0,
        instructions: instructions.trim(),
        ingredients: ingredients
          .filter((i) => i.name.trim())
          .map((i) => ({ name: i.name, quantity: i.quantity || 0, unit: i.unit })),
      });
      Alert.alert('Gespeichert', 'Rezept wurde erstellt.');
      router.back();
    } catch (err: any) {
      Alert.alert('Fehler', err.message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Text style={styles.sectionTitle}>Rezept erstellen</Text>

        <Text style={styles.label}>Titel *</Text>
        <TextInput style={styles.input} value={title} onChangeText={setTitle} placeholder="Rezeptname" />

        <Text style={styles.label}>Beschreibung</Text>
        <TextInput
          style={[styles.input, { height: 80, textAlignVertical: 'top' }]}
          value={description}
          onChangeText={setDescription}
          placeholder="Kurze Beschreibung..."
          multiline
        />

        <View style={styles.row}>
          <View style={styles.halfField}>
            <Text style={styles.label}>Vorbereitung (min)</Text>
            <TextInput style={styles.input} value={prepTime} onChangeText={setPrepTime} keyboardType="numeric" />
          </View>
          <View style={styles.halfField}>
            <Text style={styles.label}>Kochen (min)</Text>
            <TextInput style={styles.input} value={cookTime} onChangeText={setCookTime} keyboardType="numeric" />
          </View>
        </View>

        <View style={styles.row}>
          <View style={styles.halfField}>
            <Text style={styles.label}>Portionen</Text>
            <TextInput style={styles.input} value={servings} onChangeText={setServings} keyboardType="numeric" />
          </View>
          <View style={styles.halfField}>
            <Text style={styles.label}>kcal / Portion</Text>
            <TextInput style={styles.input} value={calories} onChangeText={setCalories} keyboardType="numeric" />
          </View>
        </View>

        <Text style={styles.sectionTitle}>Naehrwerte pro Portion</Text>
        <View style={styles.row}>
          <View style={styles.thirdField}>
            <Text style={styles.label}>Protein (g)</Text>
            <TextInput style={styles.input} value={proteinG} onChangeText={setProteinG} keyboardType="numeric" />
          </View>
          <View style={styles.thirdField}>
            <Text style={styles.label}>Kohlenh. (g)</Text>
            <TextInput style={styles.input} value={carbsG} onChangeText={setCarbsG} keyboardType="numeric" />
          </View>
          <View style={styles.thirdField}>
            <Text style={styles.label}>Fett (g)</Text>
            <TextInput style={styles.input} value={fatG} onChangeText={setFatG} keyboardType="numeric" />
          </View>
        </View>

        <View style={styles.ingredientHeader}>
          <Text style={styles.sectionTitle}>Zutaten</Text>
          <TouchableOpacity onPress={addIngredient}>
            <Ionicons name="add-circle" size={28} color={PRIMARY} />
          </TouchableOpacity>
        </View>
        {ingredients.map((ing, i) => (
          <View key={i} style={styles.ingredientRow}>
            <TextInput
              style={[styles.input, { flex: 2 }]}
              value={ing.name}
              onChangeText={(v) => updateIngredient(i, 'name', v)}
              placeholder="Zutat"
            />
            <TextInput
              style={[styles.input, { flex: 1, marginHorizontal: 6 }]}
              value={ing.quantity ? String(ing.quantity) : ''}
              onChangeText={(v) => updateIngredient(i, 'quantity', v)}
              placeholder="Menge"
              keyboardType="numeric"
            />
            <TextInput
              style={[styles.input, { flex: 1 }]}
              value={ing.unit}
              onChangeText={(v) => updateIngredient(i, 'unit', v)}
              placeholder="Einheit"
            />
            <TouchableOpacity onPress={() => removeIngredient(i)} style={{ marginLeft: 6 }}>
              <Ionicons name="close-circle" size={22} color="#E53935" />
            </TouchableOpacity>
          </View>
        ))}

        <Text style={styles.sectionTitle}>Zubereitung</Text>
        <TextInput
          style={[styles.input, { height: 120, textAlignVertical: 'top' }]}
          value={instructions}
          onChangeText={setInstructions}
          placeholder="Schritt fuer Schritt Anleitung..."
          multiline
        />

        <TouchableOpacity style={styles.saveButton} onPress={handleSave} disabled={saving}>
          <Ionicons name="checkmark-circle-outline" size={20} color="#fff" />
          <Text style={styles.saveButtonText}>{saving ? 'Speichern...' : 'Rezept speichern'}</Text>
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 40 },
  sectionTitle: { fontSize: 18, fontWeight: '600', color: '#212121', marginTop: 16, marginBottom: 8 },
  label: { fontSize: 13, fontWeight: '500', color: '#757575', marginBottom: 4, marginTop: 8 },
  input: {
    backgroundColor: '#fff',
    borderRadius: 10,
    paddingHorizontal: 14,
    paddingVertical: 12,
    fontSize: 15,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.04,
    shadowRadius: 2,
    elevation: 1,
  },
  row: { flexDirection: 'row', gap: 10 },
  halfField: { flex: 1 },
  thirdField: { flex: 1 },
  ingredientHeader: { flexDirection: 'row', justifyContent: 'space-between', alignItems: 'center', marginTop: 16, marginBottom: 8 },
  ingredientRow: { flexDirection: 'row', alignItems: 'center', marginBottom: 8 },
  saveButton: {
    flexDirection: 'row',
    backgroundColor: PRIMARY,
    borderRadius: 12,
    paddingVertical: 16,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
    marginTop: 24,
  },
  saveButtonText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
