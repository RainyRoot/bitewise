import React from 'react';
import {
  View,
  Text,
  StyleSheet,
  ScrollView,
  TouchableOpacity,
  SafeAreaView,
} from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { useAuth } from '@/hooks/useAuth';
import { useRouter } from 'expo-router';

const PRIMARY = '#4CAF50';
const BACKGROUND = '#F5F5F5';

// TODO: Replace with actual user data
const PLACEHOLDER_STATS = {
  streak: 5,
  recipesCooked: 23,
  waterStreak: 3,
};

const ALLERGENS = ['Laktose', 'Nüsse'];

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

  const handleLogout = async () => {
    await logout();
    router.replace('/login');
  };

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

        {/* Quick Stats */}
        <View style={styles.statsRow}>
          <View style={styles.statItem}>
            <Ionicons name="flame" size={24} color="#FF9800" />
            <Text style={styles.statValue}>{PLACEHOLDER_STATS.streak}</Text>
            <Text style={styles.statLabel}>Tage Streak</Text>
          </View>
          <View style={styles.statItem}>
            <Ionicons name="restaurant" size={24} color={PRIMARY} />
            <Text style={styles.statValue}>{PLACEHOLDER_STATS.recipesCooked}</Text>
            <Text style={styles.statLabel}>Rezepte</Text>
          </View>
          <View style={styles.statItem}>
            <Ionicons name="water" size={24} color="#2196F3" />
            <Text style={styles.statValue}>{PLACEHOLDER_STATS.waterStreak}</Text>
            <Text style={styles.statLabel}>Wasser-Streak</Text>
          </View>
        </View>

        {/* Allergies */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Allergien & Unverträglichkeiten</Text>
          <View style={styles.allergenRow}>
            {ALLERGENS.map((a) => (
              <View key={a} style={styles.allergenChip}>
                <Text style={styles.allergenText}>{a}</Text>
              </View>
            ))}
            <TouchableOpacity style={styles.allergenAddChip}>
              <Ionicons name="add" size={16} color={PRIMARY} />
            </TouchableOpacity>
          </View>
        </View>

        {/* Menu */}
        <View style={styles.card}>
          <MenuItem icon="body-outline" label="Körperdaten" value="75 kg, 178 cm" />
          <MenuItem icon="fitness-outline" label="Aktivitätslevel" value="Moderat" />
          <MenuItem icon="nutrition-outline" label="Kalorienziel" value="2.200 kcal" />
          <MenuItem icon="water-outline" label="Wasserziel" value="2.500 ml" />
          <MenuItem icon="heart-outline" label="Favoriten" value="6 Rezepte" />
          <MenuItem icon="trophy-outline" label="Achievements" value="3 / 12" />
        </View>

        {/* Settings */}
        <View style={styles.card}>
          <MenuItem icon="notifications-outline" label="Benachrichtigungen" />
          <MenuItem icon="moon-outline" label="Erscheinungsbild" value="System" />
          <MenuItem icon="information-circle-outline" label="Über BiteWise" />
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
  container: {
    flex: 1,
    backgroundColor: BACKGROUND,
  },
  scrollContent: {
    padding: 16,
    paddingBottom: 32,
  },
  profileHeader: {
    alignItems: 'center',
    marginBottom: 20,
  },
  avatar: {
    width: 80,
    height: 80,
    borderRadius: 40,
    backgroundColor: PRIMARY,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 12,
  },
  avatarText: {
    fontSize: 32,
    fontWeight: 'bold',
    color: '#fff',
  },
  userName: {
    fontSize: 22,
    fontWeight: 'bold',
    color: '#212121',
  },
  userEmail: {
    fontSize: 14,
    color: '#757575',
    marginTop: 2,
  },
  statsRow: {
    flexDirection: 'row',
    justifyContent: 'space-around',
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
  statItem: {
    alignItems: 'center',
  },
  statValue: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#212121',
    marginTop: 4,
  },
  statLabel: {
    fontSize: 11,
    color: '#757575',
    marginTop: 2,
  },
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
  cardTitle: {
    fontSize: 16,
    fontWeight: '600',
    color: '#212121',
    marginBottom: 12,
  },
  allergenRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 8,
  },
  allergenChip: {
    backgroundColor: '#FFEBEE',
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 16,
  },
  allergenText: {
    fontSize: 13,
    color: '#E53935',
    fontWeight: '500',
  },
  allergenAddChip: {
    backgroundColor: '#E8F5E9',
    width: 32,
    height: 32,
    borderRadius: 16,
    justifyContent: 'center',
    alignItems: 'center',
  },
  menuItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 14,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
  },
  menuLabel: {
    fontSize: 15,
    color: '#212121',
    marginLeft: 12,
    flex: 1,
  },
  menuRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 4,
  },
  menuValue: {
    fontSize: 13,
    color: '#757575',
  },
  logoutButton: {
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    gap: 8,
    paddingVertical: 16,
    marginTop: 8,
  },
  logoutText: {
    fontSize: 16,
    fontWeight: '600',
    color: '#E53935',
  },
});
