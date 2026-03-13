import React, { useState, useEffect } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  SafeAreaView,
  ScrollView,
  Alert,
  ActivityIndicator,
} from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { BarCodeScanner } from 'expo-barcode-scanner';
import { nutritionLookup, tracking } from '@/services/api';
import type { FoodItem } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

export default function ScannerScreen() {
  const router = useRouter();
  const [barcode, setBarcode] = useState('');
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<FoodItem | null>(null);
  const [grams, setGrams] = useState('100');
  const [mealType, setMealType] = useState<'breakfast' | 'lunch' | 'dinner' | 'snack'>('lunch');
  const [hasPermission, setHasPermission] = useState<boolean | null>(null);
  const [showCamera, setShowCamera] = useState(false);
  const [scanned, setScanned] = useState(false);

  useEffect(() => {
    (async () => {
      const { status } = await BarCodeScanner.requestPermissionsAsync();
      setHasPermission(status === 'granted');
    })();
  }, []);

  const handleBarCodeScanned = ({ data }: { type: string; data: string }) => {
    if (scanned) return;
    setScanned(true);
    setBarcode(data);
    setShowCamera(false);
    lookupBarcode(data);
  };

  const lookupBarcode = async (code: string) => {
    if (!code.trim()) return;
    setLoading(true);
    setResult(null);
    try {
      const item = await nutritionLookup.lookupBarcode(code.trim());
      setResult(item);
    } catch {
      Alert.alert('Nicht gefunden', 'Produkt konnte nicht gefunden werden.');
    } finally {
      setLoading(false);
    }
  };

  const handleLookup = () => lookupBarcode(barcode);

  const handleLog = async () => {
    if (!result) return;
    const factor = (parseInt(grams, 10) || 100) / 100;
    try {
      await tracking.logFood({
        food_name: result.brand ? `${result.name} (${result.brand})` : result.name,
        barcode: result.barcode,
        meal_type: mealType,
        servings: 1,
        calories: Math.round(result.calories_per_100g * factor),
        protein_g: Math.round(result.protein_g_per_100g * factor * 10) / 10,
        carbs_g: Math.round(result.carbs_g_per_100g * factor * 10) / 10,
        fat_g: Math.round(result.fat_g_per_100g * factor * 10) / 10,
        fiber_g: Math.round(result.fiber_g_per_100g * factor * 10) / 10,
      });
      Alert.alert('Erfasst', `${result.name} wurde zu ${MEAL_LABELS[mealType]} hinzugefuegt.`);
      router.back();
    } catch (err: any) {
      Alert.alert('Fehler', err.message);
    }
  };

  const MEAL_LABELS: Record<string, string> = {
    breakfast: 'Fruehstueck',
    lunch: 'Mittagessen',
    dinner: 'Abendessen',
    snack: 'Snack',
  };

  if (showCamera) {
    return (
      <SafeAreaView style={styles.container}>
        <View style={styles.cameraContainer}>
          <BarCodeScanner
            onBarCodeScanned={scanned ? undefined : handleBarCodeScanned}
            style={StyleSheet.absoluteFillObject}
          />
          <View style={styles.overlay}>
            <View style={styles.scanWindow} />
          </View>
          <Text style={styles.cameraHint}>Barcode in das Feld halten</Text>
          <TouchableOpacity
            style={styles.cancelButton}
            onPress={() => {
              setShowCamera(false);
              setScanned(false);
            }}
          >
            <Text style={styles.cancelButtonText}>Abbrechen</Text>
          </TouchableOpacity>
        </View>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <View style={styles.header}>
          <Ionicons name="barcode-outline" size={48} color={PRIMARY} />
          <Text style={styles.title}>Barcode Scanner</Text>
          <Text style={styles.subtitle}>Barcode scannen oder eingeben</Text>
        </View>

        {/* Camera Button */}
        {hasPermission !== false && (
          <TouchableOpacity
            style={styles.cameraButton}
            onPress={() => {
              setScanned(false);
              setShowCamera(true);
            }}
          >
            <Ionicons name="camera-outline" size={24} color="#fff" />
            <Text style={styles.cameraButtonText}>Kamera scannen</Text>
          </TouchableOpacity>
        )}

        {hasPermission === false && (
          <Text style={styles.permissionText}>
            Kamera-Berechtigung nicht erteilt. Bitte in den Einstellungen aktivieren.
          </Text>
        )}

        {/* Manual Entry */}
        <Text style={styles.orText}>oder manuell eingeben</Text>
        <View style={styles.inputRow}>
          <TextInput
            style={styles.barcodeInput}
            placeholder="Barcode eingeben..."
            value={barcode}
            onChangeText={setBarcode}
            keyboardType="numeric"
            returnKeyType="search"
            onSubmitEditing={handleLookup}
          />
          <TouchableOpacity style={styles.searchButton} onPress={handleLookup}>
            <Ionicons name="search" size={22} color="#fff" />
          </TouchableOpacity>
        </View>

        {loading && (
          <ActivityIndicator size="large" color={PRIMARY} style={{ marginTop: 32 }} />
        )}

        {result && (
          <View style={styles.resultCard}>
            <Text style={styles.resultName}>{result.name}</Text>
            {result.brand ? <Text style={styles.resultBrand}>{result.brand}</Text> : null}

            <Text style={styles.sectionLabel}>Naehrwerte pro 100g</Text>
            <View style={styles.nutritionGrid}>
              <View style={styles.nutritionItem}>
                <Text style={styles.nutritionValue}>{Math.round(result.calories_per_100g)}</Text>
                <Text style={styles.nutritionLabel}>kcal</Text>
              </View>
              <View style={styles.nutritionItem}>
                <Text style={styles.nutritionValue}>{result.protein_g_per_100g.toFixed(1)}g</Text>
                <Text style={styles.nutritionLabel}>Protein</Text>
              </View>
              <View style={styles.nutritionItem}>
                <Text style={styles.nutritionValue}>{result.carbs_g_per_100g.toFixed(1)}g</Text>
                <Text style={styles.nutritionLabel}>Kohlenh.</Text>
              </View>
              <View style={styles.nutritionItem}>
                <Text style={styles.nutritionValue}>{result.fat_g_per_100g.toFixed(1)}g</Text>
                <Text style={styles.nutritionLabel}>Fett</Text>
              </View>
            </View>

            <Text style={styles.sectionLabel}>Menge (g)</Text>
            <TextInput
              style={styles.gramsInput}
              value={grams}
              onChangeText={setGrams}
              keyboardType="numeric"
            />

            <Text style={styles.sectionLabel}>Mahlzeit</Text>
            <View style={styles.mealSelector}>
              {(['breakfast', 'lunch', 'dinner', 'snack'] as const).map((type) => (
                <TouchableOpacity
                  key={type}
                  style={[styles.mealTab, mealType === type && styles.mealTabActive]}
                  onPress={() => setMealType(type)}
                >
                  <Text style={[styles.mealTabText, mealType === type && styles.mealTabTextActive]}>
                    {MEAL_LABELS[type]}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>

            <TouchableOpacity style={styles.logButton} onPress={handleLog}>
              <Ionicons name="add-circle-outline" size={20} color="#fff" />
              <Text style={styles.logButtonText}>Erfassen</Text>
            </TouchableOpacity>
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  header: { alignItems: 'center', marginBottom: 24, marginTop: 16 },
  title: { fontSize: 24, fontWeight: 'bold', color: '#212121', marginTop: 12 },
  subtitle: { fontSize: 14, color: '#757575', marginTop: 4 },
  cameraButton: {
    flexDirection: 'row',
    backgroundColor: PRIMARY,
    borderRadius: 12,
    paddingVertical: 16,
    justifyContent: 'center',
    alignItems: 'center',
    gap: 10,
    marginBottom: 16,
  },
  cameraButtonText: { color: '#fff', fontSize: 16, fontWeight: '600' },
  permissionText: { fontSize: 13, color: '#E53935', textAlign: 'center', marginBottom: 16 },
  orText: { fontSize: 13, color: '#BDBDBD', textAlign: 'center', marginBottom: 12 },
  inputRow: { flexDirection: 'row', gap: 8 },
  barcodeInput: {
    flex: 1,
    backgroundColor: '#fff',
    borderRadius: 12,
    paddingHorizontal: 16,
    paddingVertical: 14,
    fontSize: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.05,
    shadowRadius: 4,
    elevation: 1,
  },
  searchButton: { backgroundColor: PRIMARY, borderRadius: 12, width: 50, justifyContent: 'center', alignItems: 'center' },
  resultCard: { backgroundColor: '#fff', borderRadius: 16, padding: 20, marginTop: 20, shadowColor: '#000', shadowOffset: { width: 0, height: 2 }, shadowOpacity: 0.08, shadowRadius: 8, elevation: 2 },
  resultName: { fontSize: 20, fontWeight: 'bold', color: '#212121' },
  resultBrand: { fontSize: 14, color: '#757575', marginTop: 2 },
  sectionLabel: { fontSize: 14, fontWeight: '600', color: '#757575', marginTop: 16, marginBottom: 8 },
  nutritionGrid: { flexDirection: 'row', justifyContent: 'space-around' },
  nutritionItem: { alignItems: 'center' },
  nutritionValue: { fontSize: 18, fontWeight: 'bold', color: PRIMARY },
  nutritionLabel: { fontSize: 11, color: '#757575', marginTop: 2 },
  gramsInput: { backgroundColor: '#F5F5F5', borderRadius: 8, paddingHorizontal: 14, paddingVertical: 12, fontSize: 16 },
  mealSelector: { flexDirection: 'row', gap: 6 },
  mealTab: { flex: 1, paddingVertical: 8, borderRadius: 8, backgroundColor: '#F5F5F5', alignItems: 'center' },
  mealTabActive: { backgroundColor: PRIMARY },
  mealTabText: { fontSize: 11, fontWeight: '600', color: '#757575' },
  mealTabTextActive: { color: '#fff' },
  logButton: { flexDirection: 'row', backgroundColor: PRIMARY, borderRadius: 12, paddingVertical: 14, justifyContent: 'center', alignItems: 'center', gap: 8, marginTop: 20 },
  logButtonText: { color: '#fff', fontSize: 16, fontWeight: '600' },
  // Camera styles
  cameraContainer: { flex: 1 },
  overlay: {
    ...StyleSheet.absoluteFillObject,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: 'rgba(0,0,0,0.5)',
  },
  scanWindow: {
    width: 280,
    height: 160,
    borderWidth: 3,
    borderColor: PRIMARY,
    borderRadius: 16,
    backgroundColor: 'transparent',
  },
  cameraHint: {
    position: 'absolute',
    bottom: 120,
    alignSelf: 'center',
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  cancelButton: {
    position: 'absolute',
    bottom: 50,
    alignSelf: 'center',
    backgroundColor: 'rgba(0,0,0,0.6)',
    paddingHorizontal: 32,
    paddingVertical: 14,
    borderRadius: 12,
  },
  cancelButtonText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
