import React, { useState } from 'react';
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
import { useRouter } from 'expo-router';
import { useAuth } from '@/hooks/useAuth';
import { dataExport } from '@/services/api';

const PRIMARY = '#4CAF50';

export default function DataExportScreen() {
  const router = useRouter();
  const { logout } = useAuth();
  const [exporting, setExporting] = useState<string | null>(null);

  const handleExportCSV = async () => {
    setExporting('csv');
    try {
      const csv = await dataExport.downloadCSV();
      Alert.alert('Export erfolgreich', `${csv.split('\n').length - 1} Zeilen exportiert (CSV)`);
    } catch {
      Alert.alert('Fehler', 'CSV-Export fehlgeschlagen');
    } finally {
      setExporting(null);
    }
  };

  const handleExportJSON = async () => {
    setExporting('json');
    try {
      const json = await dataExport.downloadJSON();
      const data = JSON.parse(json);
      const logCount = data.food_logs?.length || 0;
      Alert.alert('Export erfolgreich', `Profil + ${logCount} Eintraege exportiert (JSON)`);
    } catch {
      Alert.alert('Fehler', 'JSON-Export fehlgeschlagen');
    } finally {
      setExporting(null);
    }
  };

  const handleDeleteAccount = () => {
    Alert.alert(
      'Konto loeschen',
      'Bist du sicher? Alle Daten werden unwiderruflich geloescht.',
      [
        { text: 'Abbrechen', style: 'cancel' },
        {
          text: 'Konto loeschen',
          style: 'destructive',
          onPress: async () => {
            try {
              await dataExport.deleteAccount();
              await logout();
              router.replace('/login');
            } catch {
              Alert.alert('Fehler', 'Konto konnte nicht geloescht werden');
            }
          },
        },
      ],
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      <ScrollView contentContainerStyle={styles.scrollContent}>
        {/* Info */}
        <View style={styles.card}>
          <Ionicons name="cloud-download-outline" size={48} color={PRIMARY} style={{ alignSelf: 'center', marginBottom: 12 }} />
          <Text style={styles.cardTitle}>Deine Daten</Text>
          <Text style={styles.cardDesc}>
            Du kannst jederzeit alle deine Daten herunterladen oder dein Konto vollstaendig loeschen.
          </Text>
        </View>

        {/* Export Options */}
        <View style={styles.card}>
          <Text style={styles.cardTitle}>Daten exportieren</Text>

          <TouchableOpacity style={styles.exportBtn} onPress={handleExportCSV} disabled={exporting !== null}>
            <Ionicons name="document-text-outline" size={24} color={PRIMARY} />
            <View style={{ flex: 1 }}>
              <Text style={styles.exportTitle}>CSV-Export</Text>
              <Text style={styles.exportDesc}>Ernaehrungstagebuch als Tabelle</Text>
            </View>
            {exporting === 'csv' ? (
              <ActivityIndicator color={PRIMARY} />
            ) : (
              <Ionicons name="download-outline" size={22} color={PRIMARY} />
            )}
          </TouchableOpacity>

          <TouchableOpacity style={styles.exportBtn} onPress={handleExportJSON} disabled={exporting !== null}>
            <Ionicons name="code-outline" size={24} color={PRIMARY} />
            <View style={{ flex: 1 }}>
              <Text style={styles.exportTitle}>JSON-Export</Text>
              <Text style={styles.exportDesc}>Alle Daten inkl. Profil</Text>
            </View>
            {exporting === 'json' ? (
              <ActivityIndicator color={PRIMARY} />
            ) : (
              <Ionicons name="download-outline" size={22} color={PRIMARY} />
            )}
          </TouchableOpacity>
        </View>

        {/* Danger Zone */}
        <View style={[styles.card, { borderColor: '#FFCDD2', borderWidth: 1 }]}>
          <Text style={[styles.cardTitle, { color: '#F44336' }]}>Gefahrenzone</Text>
          <Text style={styles.cardDesc}>
            Das Loeschen deines Kontos kann nicht rueckgaengig gemacht werden. Exportiere deine Daten vorher.
          </Text>
          <TouchableOpacity style={styles.deleteBtn} onPress={handleDeleteAccount}>
            <Ionicons name="trash-outline" size={18} color="#fff" />
            <Text style={styles.deleteText}>Konto endgueltig loeschen</Text>
          </TouchableOpacity>
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
  cardTitle: { fontSize: 16, fontWeight: '600', color: '#212121', marginBottom: 8 },
  cardDesc: { fontSize: 14, color: '#757575', lineHeight: 20, marginBottom: 4 },
  exportBtn: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    paddingVertical: 14,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
  },
  exportTitle: { fontSize: 15, fontWeight: '500', color: '#212121' },
  exportDesc: { fontSize: 12, color: '#757575', marginTop: 2 },
  deleteBtn: {
    backgroundColor: '#F44336',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    borderRadius: 12,
    paddingVertical: 14,
    marginTop: 12,
  },
  deleteText: { color: '#fff', fontSize: 15, fontWeight: '600' },
});
