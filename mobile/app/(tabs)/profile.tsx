import React, { useState, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
  ActivityIndicator,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useAuth } from '@/hooks/useAuth';
import { useRouter, useFocusEffect } from 'expo-router';
import { profile } from '@/services/api';
import { useI18n } from '@/i18n';
import type { Locale } from '@/i18n';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

interface MenuItemProps {
  icon: string;
  label: string;
  value?: string;
  onPress?: () => void;
}

function MenuItem({ icon, label, value, onPress }: MenuItemProps) {
  return (
    <TouchableOpacity style={styles.menuItem} onPress={onPress}>
      <Ionicons name={icon as any} size={22} color="#757575" />
      <Text style={styles.menuLabel}>{label}</Text>
      <View style={styles.menuRight}>
        {value && <Text style={styles.menuValue}>{value}</Text>}
        <Ionicons name="chevron-forward" size={18} color="#BDBDBD" />
      </View>
    </TouchableOpacity>
  );
}

export default function ProfileScreen() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const { locale, setLocale, t } = useI18n();
  const [allergies, setAllergies] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    try {
      const allergyList = await profile.getAllergies();
      setAllergies(allergyList || []);
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

  const handleLogout = async () => {
    await logout();
    router.replace('/login');
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
        {/* Profile Header */}
        <View style={styles.profileHeader}>
          <View style={styles.avatar}>
            <Text style={styles.avatarText}>
              {user?.name?.[0]?.toUpperCase() || 'B'}
            </Text>
          </View>
          <Text style={styles.userName}>{user?.name || 'BiteWise User'}</Text>
          <Text style={styles.userEmail}>{user?.email || 'user@example.com'}</Text>
        </View>

        {/* Allergies */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>{t.allergies}</Text>
          <View style={styles.allergenRow}>
            {allergies.length > 0 ? (
              allergies.map((a) => (
                <View key={a} style={styles.allergenChip}>
                  <Text style={styles.allergenText}>{a}</Text>
                </View>
              ))
            ) : (
              <Text style={{ color: '#BDBDBD', fontStyle: 'italic' }}>{t.no_allergies}</Text>
            )}
          </View>
        </View>

        {/* Menu */}
        <View style={styles.card}>
          <MenuItem
            icon="body-outline"
            label="Körperdaten"
            value={user?.weight_kg ? `${user.weight_kg} kg, ${user.height_cm} cm` : ''}
          />
          <MenuItem
            icon="fitness-outline"
            label="Aktivitätslevel"
            value={user?.activity_level || ''}
          />
          <MenuItem
            icon="nutrition-outline"
            label="Kalorienziel"
            value={user?.calorie_target ? `${user.calorie_target} kcal` : ''}
          />
          <MenuItem
            icon="water-outline"
            label="Wasserziel"
            value={user?.daily_water_ml_goal ? `${user.daily_water_ml_goal} ml` : ''}
          />
          <MenuItem icon="heart-outline" label={t.favorites} />
          <MenuItem icon="trophy-outline" label={t.achievements} onPress={() => router.push('/achievements')} />
          <MenuItem icon="stats-chart-outline" label={t.statistics} onPress={() => router.push('/stats')} />
          <MenuItem icon="book-outline" label={t.my_recipes} onPress={() => router.push('/create-recipe')} />
          <MenuItem icon="journal-outline" label={t.diary_title} onPress={() => router.push('/diary')} />
          <MenuItem icon="pricetag-outline" label={t.prices_title} onPress={() => router.push('/price-tracker')} />
          <MenuItem icon="people-outline" label={t.friends_title} onPress={() => router.push('/friends')} />
        </View>

        {/* Settings */}
        <View style={styles.card}>
          <MenuItem icon="notifications-outline" label={t.notifications} onPress={() => router.push('/notification-settings')} />
          <MenuItem
            icon="language-outline"
            label="Sprache / Language"
            value={locale === 'de' ? 'Deutsch' : 'English'}
            onPress={() => setLocale(locale === 'de' ? 'en' as Locale : 'de' as Locale)}
          />
          <MenuItem icon="cloud-download-outline" label={t.export_title} onPress={() => router.push('/data-export')} />
          <MenuItem icon="moon-outline" label={t.appearance} value="System" />
          <MenuItem icon="information-circle-outline" label={t.about} />
        </View>

        {/* Logout */}
        <TouchableOpacity style={styles.logoutButton} onPress={handleLogout}>
          <Ionicons name="log-out-outline" size={20} color="#E53935" />
          <Text style={styles.logoutText}>Abmelden</Text>
        </TouchableOpacity>
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: BACKGROUND },
  scrollContent: { padding: 16, paddingBottom: 32 },
  profileHeader: { alignItems: 'center', marginBottom: 20 },
  avatar: { width: 80, height: 80, borderRadius: 40, backgroundColor: PRIMARY, justifyContent: 'center', alignItems: 'center', marginBottom: 12 },
  avatarText: { fontSize: 32, fontWeight: 'bold', color: '#fff' },
  userName: { fontSize: 22, fontWeight: 'bold', color: '#212121' },
  userEmail: { fontSize: 14, color: '#757575', marginTop: 2 },
  card: { backgroundColor: '#fff', borderRadius: 16, padding: 16, marginBottom: 16, shadowColor: '#000', shadowOffset: { width: 0, height: 2 }, shadowOpacity: 0.08, shadowRadius: 8, elevation: 2 },
  cardTitle: { fontSize: 16, fontWeight: '600', color: '#212121', marginBottom: 12 },
  allergenRow: { flexDirection: 'row', flexWrap: 'wrap', gap: 8 },
  allergenChip: { backgroundColor: '#FFEBEE', paddingHorizontal: 12, paddingVertical: 6, borderRadius: 16 },
  allergenText: { fontSize: 13, color: '#E53935', fontWeight: '500' },
  menuItem: { flexDirection: 'row', alignItems: 'center', paddingVertical: 14, borderBottomWidth: 1, borderBottomColor: '#F5F5F5' },
  menuLabel: { fontSize: 15, color: '#212121', marginLeft: 12, flex: 1 },
  menuRight: { flexDirection: 'row', alignItems: 'center', gap: 4 },
  menuValue: { fontSize: 13, color: '#757575' },
  logoutButton: { flexDirection: 'row', justifyContent: 'center', alignItems: 'center', gap: 8, paddingVertical: 16, marginTop: 8 },
  logoutText: { fontSize: 16, fontWeight: '600', color: '#E53935' },
});
