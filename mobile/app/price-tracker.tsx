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
  Alert,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useFocusEffect } from 'expo-router';
import { prices } from '@/services/api';
import type { PriceLog, SpendingSummary } from '@/types';

const PRIMARY = '#4CAF50';

export default function PriceTrackerScreen() {
  const [logs, setLogs] = useState<PriceLog[]>([]);
  const [spending, setSpending] = useState<SpendingSummary | null>(null);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);

  // Form state
  const [itemName, setItemName] = useState('');
  const [priceEur, setPriceEur] = useState('');
  const [store, setStore] = useState('');
  const [saving, setSaving] = useState(false);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [logData, spendData] = await Promise.all([
        prices.getLogs(30),
        prices.getSpending(),
      ]);
      setLogs(logData || []);
      setSpending(spendData);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchData();
    }, [])
  );

  const handleAdd = async () => {
    if (!itemName.trim() || !priceEur.trim()) {
      Alert.alert('Fehler', 'Produkt und Preis sind erforderlich');
      return;
    }

    const cents = Math.round(parseFloat(priceEur.replace(',', '.')) * 100);
    if (isNaN(cents) || cents <= 0) {
      Alert.alert('Fehler', 'Bitte einen gueltigen Preis eingeben');
      return;
    }

    setSaving(true);
    try {
      await prices.log({
        item_name: itemName.trim(),
        price_cents: cents,
        store: store.trim(),
        date: new Date().toISOString().split('T')[0],
      });
      setItemName('');
      setPriceEur('');
      setStore('');
      setShowForm(false);
      fetchData();
    } catch {
      Alert.alert('Fehler', 'Preis konnte nicht gespeichert werden');
    } finally {
      setSaving(false);
    }
  };

  const formatPrice = (cents: number) => {
    return (cents / 100).toFixed(2).replace('.', ',') + ' EUR';
  };

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
        {/* Monthly Spending */}
        {spending && (
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Ausgaben diesen Monat</Text>
            <Text style={styles.spendingAmount}>{formatPrice(spending.total_cents)}</Text>
            <Text style={styles.spendingDetail}>{spending.item_count} Eintraege</Text>
          </View>
        )}

        {/* Add Button */}
        <TouchableOpacity
          style={styles.addButton}
          onPress={() => setShowForm(!showForm)}
        >
          <Ionicons name={showForm ? 'close' : 'add'} size={20} color="#fff" />
          <Text style={styles.addButtonText}>
            {showForm ? 'Abbrechen' : 'Preis eintragen'}
          </Text>
        </TouchableOpacity>

        {/* Add Form */}
        {showForm && (
          <View style={styles.card}>
            <TextInput
              style={styles.input}
              placeholder="Produkt (z.B. Milch 1L)"
              placeholderTextColor="#BDBDBD"
              value={itemName}
              onChangeText={setItemName}
            />
            <TextInput
              style={styles.input}
              placeholder="Preis in EUR (z.B. 1,49)"
              placeholderTextColor="#BDBDBD"
              keyboardType="decimal-pad"
              value={priceEur}
              onChangeText={setPriceEur}
            />
            <TextInput
              style={styles.input}
              placeholder="Laden (optional)"
              placeholderTextColor="#BDBDBD"
              value={store}
              onChangeText={setStore}
            />
            <TouchableOpacity style={styles.saveButton} onPress={handleAdd} disabled={saving}>
              {saving ? (
                <ActivityIndicator color="#fff" />
              ) : (
                <Text style={styles.saveText}>Speichern</Text>
              )}
            </TouchableOpacity>
          </View>
        )}

        {/* Recent Logs */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Letzte Eintraege</Text>
          {logs.length === 0 ? (
            <Text style={styles.emptyText}>Noch keine Preise eingetragen</Text>
          ) : (
            logs.map((log) => (
              <View key={log.id} style={styles.logItem}>
                <View style={{ flex: 1 }}>
                  <Text style={styles.logName}>{log.item_name}</Text>
                  <Text style={styles.logDetail}>
                    {log.store ? `${log.store} · ` : ''}{log.date}
                  </Text>
                </View>
                <Text style={styles.logPrice}>{formatPrice(log.price_cents)}</Text>
              </View>
            ))
          )}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#F5F5F5' },
  scrollContent: { padding: 16, paddingBottom: 32 },
  card: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 16,
    marginBottom: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
  },
  cardTitle: { fontSize: 16, fontWeight: '600', color: '#212121', marginBottom: 12 },
  spendingAmount: { fontSize: 28, fontWeight: 'bold', color: PRIMARY, textAlign: 'center' },
  spendingDetail: { fontSize: 13, color: '#757575', textAlign: 'center', marginTop: 4 },
  addButton: {
    backgroundColor: PRIMARY,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    borderRadius: 12,
    paddingVertical: 12,
    marginBottom: 16,
  },
  addButtonText: { color: '#fff', fontSize: 15, fontWeight: '600' },
  input: {
    borderWidth: 1,
    borderColor: '#E0E0E0',
    borderRadius: 10,
    padding: 12,
    fontSize: 15,
    color: '#212121',
    marginBottom: 10,
  },
  saveButton: {
    backgroundColor: PRIMARY,
    borderRadius: 10,
    paddingVertical: 12,
    alignItems: 'center',
    marginTop: 4,
  },
  saveText: { color: '#fff', fontSize: 15, fontWeight: '600' },
  emptyText: { color: '#BDBDBD', fontStyle: 'italic', textAlign: 'center', paddingVertical: 16 },
  logItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
  },
  logName: { fontSize: 15, color: '#212121', fontWeight: '500' },
  logDetail: { fontSize: 12, color: '#757575', marginTop: 2 },
  logPrice: { fontSize: 15, fontWeight: '600', color: PRIMARY },
});
