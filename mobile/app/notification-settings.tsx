import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  SafeAreaView,
  Switch,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { useFocusEffect } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { notifications as notifApi } from '@/services/api';
import type { UserNotification } from '@/types';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

const NOTIFICATION_TYPES = [
  { type: 'meal_reminder_breakfast', label: 'Fruehstueck-Erinnerung', icon: 'sunny-outline', defaultTime: '07:30' },
  { type: 'meal_reminder_lunch', label: 'Mittagessen-Erinnerung', icon: 'restaurant-outline', defaultTime: '12:00' },
  { type: 'meal_reminder_dinner', label: 'Abendessen-Erinnerung', icon: 'moon-outline', defaultTime: '18:30' },
  { type: 'water_reminder', label: 'Wasser-Erinnerung (stuendlich)', icon: 'water-outline', defaultTime: '09:00' },
  { type: 'weekly_plan_reminder', label: 'Wochenplan-Erinnerung', icon: 'calendar-outline', defaultTime: '19:00' },
];

export default function NotificationSettingsScreen() {
  const [settings, setSettings] = useState<UserNotification[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    try {
      const data = await notifApi.getSettings();
      setSettings(data || []);
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

  const getSettingForType = (type: string): UserNotification | undefined => {
    return settings.find((s) => s.type === type);
  };

  const toggleNotification = async (type: string, defaultTime: string) => {
    const existing = getSettingForType(type);
    const newEnabled = existing ? !existing.enabled : true;

    const updated: UserNotification[] = settings.map((s) =>
      s.type === type ? { ...s, enabled: newEnabled } : s,
    );

    if (!existing) {
      updated.push({
        id: 0,
        user_id: 0,
        type,
        time: defaultTime,
        enabled: true,
      });
    }

    setSettings(updated);

    try {
      const result = await notifApi.updateSettings(
        updated.map((s) => ({ ...s })),
      );
      if (result) setSettings(result);
    } catch (err: any) {
      Alert.alert('Fehler', err.message);
      fetchData();
    }
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
        <Text style={styles.title}>Benachrichtigungen</Text>
        <Text style={styles.subtitle}>Verwalte deine Erinnerungen</Text>

        <View style={styles.card}>
          {NOTIFICATION_TYPES.map((nt, index) => {
            const setting = getSettingForType(nt.type);
            const isEnabled = setting?.enabled ?? false;

            return (
              <View
                key={nt.type}
                style={[
                  styles.notifRow,
                  index < NOTIFICATION_TYPES.length - 1 && styles.notifRowBorder,
                ]}
              >
                <Ionicons name={nt.icon as any} size={24} color={isEnabled ? PRIMARY : '#BDBDBD'} />
                <View style={styles.notifInfo}>
                  <Text style={styles.notifLabel}>{nt.label}</Text>
                  {setting && (
                    <Text style={styles.notifTime}>{setting.time} Uhr</Text>
                  )}
                </View>
                <Switch
                  value={isEnabled}
                  onValueChange={() => toggleNotification(nt.type, nt.defaultTime)}
                  trackColor={{ false: '#E0E0E0', true: '#A5D6A7' }}
                  thumbColor={isEnabled ? PRIMARY : '#BDBDBD'}
                />
              </View>
            );
          })}
        </View>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  title: { fontSize: 28, fontWeight: 'bold', color: '#212121', marginBottom: 4 },
  subtitle: { fontSize: 16, color: '#757575', marginBottom: 20 },
  card: {
    backgroundColor: '#fff',
    borderRadius: 16,
    padding: 4,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.08,
    shadowRadius: 8,
    elevation: 2,
  },
  notifRow: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 16,
  },
  notifRowBorder: {
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
  },
  notifInfo: { flex: 1, marginLeft: 14 },
  notifLabel: { fontSize: 15, fontWeight: '500', color: '#212121' },
  notifTime: { fontSize: 12, color: '#757575', marginTop: 2 },
});
