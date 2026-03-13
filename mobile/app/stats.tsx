import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  SafeAreaView,
  ActivityIndicator,
  Dimensions,
} from 'react-native';
import { useFocusEffect } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { stats as statsApi } from '@/services/api';
import type { WeeklyStats, StreakInfo } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';
const SCREEN_WIDTH = Dimensions.get('window').width;
const BAR_MAX_HEIGHT = 120;

const DAY_LABELS = ['Mo', 'Di', 'Mi', 'Do', 'Fr', 'Sa', 'So'];

export default function StatsScreen() {
  const [weekly, setWeekly] = useState<WeeklyStats | null>(null);
  const [streaks, setStreaks] = useState<StreakInfo | null>(null);
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    try {
      const [w, s] = await Promise.all([
        statsApi.getWeekly().catch(() => null),
        statsApi.getStreaks().catch(() => null),
      ]);
      setWeekly(w);
      setStreaks(s);
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

  if (loading) {
    return (
      <SafeAreaView style={styles.container}>
        <ActivityIndicator size="large" color={PRIMARY} style={{ marginTop: 48 }} />
      </SafeAreaView>
    );
  }

  const days = weekly?.days || [];
  const maxCal = Math.max(...days.map((d) => d.calories), 1);
  const barWidth = (SCREEN_WIDTH - 80) / 7;

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        <Text style={styles.title}>Statistiken</Text>

        {/* Streaks */}
        <View style={styles.streakRow}>
          <View style={styles.streakCard}>
            <Ionicons name="flame" size={32} color="#FF9800" />
            <Text style={styles.streakValue}>{streaks?.current_streak || 0}</Text>
            <Text style={styles.streakLabel}>Aktuelle Serie</Text>
          </View>
          <View style={styles.streakCard}>
            <Ionicons name="trophy" size={32} color="#FFC107" />
            <Text style={styles.streakValue}>{streaks?.longest_streak || 0}</Text>
            <Text style={styles.streakLabel}>Laengste Serie</Text>
          </View>
        </View>

        {/* Weekly Calories Bar Chart */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Kalorien diese Woche</Text>
          {days.length === 0 ? (
            <Text style={styles.emptyText}>Noch keine Daten fuer diese Woche</Text>
          ) : (
            <>
              <View style={styles.chartContainer}>
                {DAY_LABELS.map((label, i) => {
                  const day = days.find((d) => {
                    const date = new Date(d.date);
                    const dayOfWeek = date.getDay();
                    const adjustedDay = dayOfWeek === 0 ? 6 : dayOfWeek - 1;
                    return adjustedDay === i;
                  });
                  const cal = day?.calories || 0;
                  const height = maxCal > 0 ? (cal / maxCal) * BAR_MAX_HEIGHT : 0;

                  return (
                    <View key={label} style={[styles.barColumn, { width: barWidth }]}>
                      <Text style={styles.barValue}>{cal > 0 ? cal : ''}</Text>
                      <View style={[styles.bar, { height: Math.max(height, 2) }]} />
                      <Text style={styles.barLabel}>{label}</Text>
                    </View>
                  );
                })}
              </View>
              <View style={styles.avgRow}>
                <Text style={styles.avgLabel}>Durchschnitt</Text>
                <Text style={styles.avgValue}>{weekly?.avg_calories || 0} kcal</Text>
              </View>
            </>
          )}
        </View>

        {/* Weekly Macros */}
        {days.length > 0 && (
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Makros (Wochendurchschnitt)</Text>
            <View style={styles.macroRow}>
              <View style={styles.macroItem}>
                <Text style={[styles.macroValue, { color: '#4CAF50' }]}>
                  {Math.round(days.reduce((s, d) => s + d.protein_g, 0) / days.length)}g
                </Text>
                <Text style={styles.macroLabel}>Protein</Text>
              </View>
              <View style={styles.macroItem}>
                <Text style={[styles.macroValue, { color: '#FF9800' }]}>
                  {Math.round(days.reduce((s, d) => s + d.carbs_g, 0) / days.length)}g
                </Text>
                <Text style={styles.macroLabel}>Kohlenhydrate</Text>
              </View>
              <View style={styles.macroItem}>
                <Text style={[styles.macroValue, { color: '#F44336' }]}>
                  {Math.round(days.reduce((s, d) => s + d.fat_g, 0) / days.length)}g
                </Text>
                <Text style={styles.macroLabel}>Fett</Text>
              </View>
            </View>
          </View>
        )}

        {/* Water Average */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Wasser diese Woche</Text>
          <View style={styles.avgRow}>
            <Ionicons name="water" size={24} color="#2196F3" />
            <Text style={[styles.avgValue, { color: '#2196F3', marginLeft: 8 }]}>
              {weekly?.avg_water_ml || 0} ml / Tag
            </Text>
          </View>
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  title: { fontSize: 28, fontWeight: 'bold', color: '#212121', marginBottom: 16 },
  streakRow: { flexDirection: 'row', gap: 12, marginBottom: 16 },
  streakCard: {
    flex: 1,
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 20,
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
  },
  streakValue: { fontSize: 36, fontWeight: 'bold', color: '#212121', marginTop: 8 },
  streakLabel: { fontSize: 13, color: '#757575', marginTop: 4 },
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
  cardTitle: { fontSize: 18, fontWeight: '600', color: '#212121', marginBottom: 16 },
  emptyText: { fontSize: 14, color: '#BDBDBD', fontStyle: 'italic', textAlign: 'center', paddingVertical: 20 },
  chartContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-end',
    height: BAR_MAX_HEIGHT + 40,
    paddingHorizontal: 4,
  },
  barColumn: { alignItems: 'center', justifyContent: 'flex-end' },
  bar: { backgroundColor: PRIMARY, borderRadius: 4, width: 24 },
  barValue: { fontSize: 10, color: '#757575', marginBottom: 4 },
  barLabel: { fontSize: 12, color: '#757575', marginTop: 6 },
  avgRow: { flexDirection: 'row', alignItems: 'center', justifyContent: 'space-between', marginTop: 12, paddingTop: 12, borderTopWidth: 1, borderTopColor: '#F0F0F0' },
  avgLabel: { fontSize: 14, color: '#757575' },
  avgValue: { fontSize: 16, fontWeight: '600', color: PRIMARY },
  macroRow: { flexDirection: 'row', justifyContent: 'space-around' },
  macroItem: { alignItems: 'center' },
  macroValue: { fontSize: 22, fontWeight: 'bold' },
  macroLabel: { fontSize: 12, color: '#757575', marginTop: 4 },
});
