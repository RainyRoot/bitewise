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
import { useFocusEffect } from 'expo-router';
import { diary } from '@/services/api';
import type { DiaryEntry } from '@/types';

const PRIMARY = '#4CAF50';
const MOODS = [
  { key: 'great', label: 'Super', icon: 'happy', color: '#4CAF50' },
  { key: 'good', label: 'Gut', icon: 'happy-outline', color: '#8BC34A' },
  { key: 'neutral', label: 'Okay', icon: 'remove-circle-outline', color: '#FFC107' },
  { key: 'bad', label: 'Schlecht', icon: 'sad-outline', color: '#FF9800' },
  { key: 'terrible', label: 'Mies', icon: 'sad', color: '#F44336' },
] as const;

export default function DiaryScreen() {
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [entry, setEntry] = useState<DiaryEntry | null>(null);
  const [mood, setMood] = useState('neutral');
  const [energy, setEnergy] = useState(5);
  const [notes, setNotes] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const fetchEntry = async () => {
    try {
      setLoading(true);
      const e = await diary.getByDate(date);
      setEntry(e);
      if (e) {
        setMood(e.mood);
        setEnergy(e.energy_level);
        setNotes(e.notes);
      } else {
        setMood('neutral');
        setEnergy(5);
        setNotes('');
      }
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  };

  useFocusEffect(
    useCallback(() => {
      fetchEntry();
    }, [date])
  );

  const handleSave = async () => {
    setSaving(true);
    try {
      const saved = await diary.createOrUpdate({
        date,
        mood,
        energy_level: energy,
        notes,
      });
      setEntry(saved);
    } catch {
      // ignore
    } finally {
      setSaving(false);
    }
  };

  const changeDate = (days: number) => {
    const d = new Date(date);
    d.setDate(d.getDate() + days);
    setDate(d.toISOString().split('T')[0]);
  };

  const formatDate = (d: string) => {
    const date = new Date(d);
    return date.toLocaleDateString('de-DE', { weekday: 'long', day: 'numeric', month: 'long' });
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
        {/* Date Navigation */}
        <View style={styles.dateNav}>
          <TouchableOpacity onPress={() => changeDate(-1)}>
            <Ionicons name="chevron-back" size={28} color={PRIMARY} />
          </TouchableOpacity>
          <Text style={styles.dateText}>{formatDate(date)}</Text>
          <TouchableOpacity onPress={() => changeDate(1)}>
            <Ionicons name="chevron-forward" size={28} color={PRIMARY} />
          </TouchableOpacity>
        </View>

        {/* Mood Selection */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Stimmung</Text>
          <View style={styles.moodRow}>
            {MOODS.map((m) => (
              <TouchableOpacity
                key={m.key}
                style={[styles.moodItem, mood === m.key && { backgroundColor: m.color + '20' }]}
                onPress={() => setMood(m.key)}
              >
                <Ionicons
                  name={m.icon as any}
                  size={32}
                  color={mood === m.key ? m.color : '#BDBDBD'}
                />
                <Text style={[styles.moodLabel, mood === m.key && { color: m.color }]}>
                  {m.label}
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        {/* Energy Level */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Energielevel: {energy}/10</Text>
          <View style={styles.energyRow}>
            {[1, 2, 3, 4, 5, 6, 7, 8, 9, 10].map((level) => (
              <TouchableOpacity
                key={level}
                style={[
                  styles.energyDot,
                  { backgroundColor: level <= energy ? PRIMARY : '#E0E0E0' },
                ]}
                onPress={() => setEnergy(level)}
              >
                <Text style={[styles.energyText, { color: level <= energy ? '#fff' : '#757575' }]}>
                  {level}
                </Text>
              </TouchableOpacity>
            ))}
          </View>
        </View>

        {/* Notes */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Notizen</Text>
          <TextInput
            style={styles.notesInput}
            multiline
            numberOfLines={4}
            placeholder="Wie war dein Tag? Was hast du gegessen?"
            placeholderTextColor="#BDBDBD"
            value={notes}
            onChangeText={setNotes}
          />
        </View>

        {/* Save Button */}
        <TouchableOpacity style={styles.saveButton} onPress={handleSave} disabled={saving}>
          {saving ? (
            <ActivityIndicator color="#fff" />
          ) : (
            <Text style={styles.saveText}>
              {entry ? 'Aktualisieren' : 'Speichern'}
            </Text>
          )}
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#F5F5F5' },
  scrollContent: { padding: 16, paddingBottom: 32 },
  dateNav: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 16,
    paddingHorizontal: 8,
  },
  dateText: { fontSize: 18, fontWeight: '600', color: '#212121' },
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
  moodRow: { flexDirection: 'row', justifyContent: 'space-around' },
  moodItem: {
    alignItems: 'center',
    padding: 8,
    borderRadius: 12,
    minWidth: 56,
  },
  moodLabel: { fontSize: 11, color: '#757575', marginTop: 4 },
  energyRow: { flexDirection: 'row', justifyContent: 'space-between' },
  energyDot: {
    width: 28,
    height: 28,
    borderRadius: 14,
    justifyContent: 'center',
    alignItems: 'center',
  },
  energyText: { fontSize: 11, fontWeight: '600' },
  notesInput: {
    borderWidth: 1,
    borderColor: '#E0E0E0',
    borderRadius: 12,
    padding: 12,
    fontSize: 15,
    color: '#212121',
    textAlignVertical: 'top',
    minHeight: 100,
  },
  saveButton: {
    backgroundColor: PRIMARY,
    borderRadius: 12,
    paddingVertical: 14,
    alignItems: 'center',
    marginTop: 4,
  },
  saveText: { color: '#fff', fontSize: 16, fontWeight: '600' },
});
