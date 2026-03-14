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
import { friends } from '@/services/api';
import type { FriendInfo, FriendInvite, LeaderboardEntry } from '@/types';

const PRIMARY = '#4CAF50';

export default function FriendsScreen() {
  const [friendList, setFriendList] = useState<FriendInfo[]>([]);
  const [invites, setInvites] = useState<FriendInvite[]>([]);
  const [leaderboard, setLeaderboard] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<'friends' | 'leaderboard'>('friends');
  const [inviteEmail, setInviteEmail] = useState('');
  const [showInvite, setShowInvite] = useState(false);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [f, i, l] = await Promise.all([
        friends.list(),
        friends.getPendingInvites(),
        friends.getLeaderboard(),
      ]);
      setFriendList(f || []);
      setInvites(i || []);
      setLeaderboard(l || []);
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

  const handleInvite = async () => {
    if (!inviteEmail.trim()) return;
    try {
      await friends.invite(inviteEmail.trim());
      Alert.alert('Einladung gesendet', `Einladung an ${inviteEmail} gesendet`);
      setInviteEmail('');
      setShowInvite(false);
    } catch (err: any) {
      Alert.alert('Fehler', err.message || 'Einladung fehlgeschlagen');
    }
  };

  const handleRespond = async (inviteId: number, accept: boolean) => {
    try {
      await friends.respondToInvite(inviteId, accept);
      fetchData();
    } catch {
      Alert.alert('Fehler', 'Aktion fehlgeschlagen');
    }
  };

  const handleRemove = async (friendId: number, name: string) => {
    Alert.alert('Freund entfernen', `${name} wirklich entfernen?`, [
      { text: 'Abbrechen', style: 'cancel' },
      {
        text: 'Entfernen',
        style: 'destructive',
        onPress: async () => {
          await friends.remove(friendId);
          fetchData();
        },
      },
    ]);
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
        {/* Tabs */}
        <View style={styles.tabRow}>
          <TouchableOpacity
            style={[styles.tab, tab === 'friends' && styles.tabActive]}
            onPress={() => setTab('friends')}
          >
            <Text style={[styles.tabText, tab === 'friends' && styles.tabTextActive]}>
              Freunde
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[styles.tab, tab === 'leaderboard' && styles.tabActive]}
            onPress={() => setTab('leaderboard')}
          >
            <Text style={[styles.tabText, tab === 'leaderboard' && styles.tabTextActive]}>
              Rangliste
            </Text>
          </TouchableOpacity>
        </View>

        {tab === 'friends' ? (
          <>
            {/* Pending Invites */}
            {invites.length > 0 && (
              <View style={styles.card}>
                <Text style={styles.cardTitle}>Einladungen</Text>
                {invites.map((inv) => (
                  <View key={inv.id} style={styles.inviteItem}>
                    <Text style={styles.inviteEmail}>Von: {inv.to_email}</Text>
                    <View style={styles.inviteActions}>
                      <TouchableOpacity
                        style={[styles.actionBtn, { backgroundColor: PRIMARY }]}
                        onPress={() => handleRespond(inv.id, true)}
                      >
                        <Ionicons name="checkmark" size={18} color="#fff" />
                      </TouchableOpacity>
                      <TouchableOpacity
                        style={[styles.actionBtn, { backgroundColor: '#F44336' }]}
                        onPress={() => handleRespond(inv.id, false)}
                      >
                        <Ionicons name="close" size={18} color="#fff" />
                      </TouchableOpacity>
                    </View>
                  </View>
                ))}
              </View>
            )}

            {/* Invite Button */}
            <TouchableOpacity
              style={styles.addButton}
              onPress={() => setShowInvite(!showInvite)}
            >
              <Ionicons name={showInvite ? 'close' : 'person-add'} size={18} color="#fff" />
              <Text style={styles.addButtonText}>
                {showInvite ? 'Abbrechen' : 'Freund einladen'}
              </Text>
            </TouchableOpacity>

            {showInvite && (
              <View style={styles.card}>
                <TextInput
                  style={styles.input}
                  placeholder="E-Mail-Adresse"
                  placeholderTextColor="#BDBDBD"
                  keyboardType="email-address"
                  autoCapitalize="none"
                  value={inviteEmail}
                  onChangeText={setInviteEmail}
                />
                <TouchableOpacity style={styles.saveButton} onPress={handleInvite}>
                  <Text style={styles.saveText}>Einladung senden</Text>
                </TouchableOpacity>
              </View>
            )}

            {/* Friends List */}
            <View style={styles.card}>
              <Text style={styles.cardTitle}>Deine Freunde ({friendList.length})</Text>
              {friendList.length === 0 ? (
                <Text style={styles.emptyText}>Noch keine Freunde hinzugefuegt</Text>
              ) : (
                friendList.map((f) => (
                  <View key={f.user_id} style={styles.friendItem}>
                    <View style={styles.friendAvatar}>
                      <Text style={styles.friendAvatarText}>
                        {f.name[0]?.toUpperCase() || '?'}
                      </Text>
                    </View>
                    <View style={{ flex: 1 }}>
                      <Text style={styles.friendName}>{f.name}</Text>
                      <Text style={styles.friendEmail}>{f.email}</Text>
                    </View>
                    <TouchableOpacity onPress={() => handleRemove(f.user_id, f.name)}>
                      <Ionicons name="trash-outline" size={20} color="#F44336" />
                    </TouchableOpacity>
                  </View>
                ))
              )}
            </View>
          </>
        ) : (
          /* Leaderboard */
          <View style={styles.card}>
            <Text style={styles.cardTitle}>Wochen-Rangliste</Text>
            {leaderboard.length === 0 ? (
              <Text style={styles.emptyText}>Keine Daten verfuegbar</Text>
            ) : (
              leaderboard.map((entry) => (
                <View key={entry.user_id} style={styles.leaderItem}>
                  <Text style={styles.leaderRank}>#{entry.rank}</Text>
                  <View style={{ flex: 1 }}>
                    <Text style={styles.leaderName}>{entry.name}</Text>
                    <Text style={styles.leaderDetail}>
                      {entry.week_calories} kcal diese Woche
                    </Text>
                  </View>
                  {entry.rank === 1 && (
                    <Ionicons name="trophy" size={24} color="#FFC107" />
                  )}
                </View>
              ))
            )}
          </View>
        )}
      </ScrollView>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, backgroundColor: '#F5F5F5' },
  scrollContent: { padding: 16, paddingBottom: 32 },
  tabRow: { flexDirection: 'row', marginBottom: 16, gap: 8 },
  tab: {
    flex: 1,
    paddingVertical: 10,
    borderRadius: 10,
    backgroundColor: '#E0E0E0',
    alignItems: 'center',
  },
  tabActive: { backgroundColor: PRIMARY },
  tabText: { fontSize: 14, fontWeight: '600', color: '#757575' },
  tabTextActive: { color: '#fff' },
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
  inviteItem: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
  },
  inviteEmail: { fontSize: 14, color: '#212121' },
  inviteActions: { flexDirection: 'row', gap: 8 },
  actionBtn: {
    width: 32,
    height: 32,
    borderRadius: 16,
    justifyContent: 'center',
    alignItems: 'center',
  },
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
  },
  saveText: { color: '#fff', fontSize: 15, fontWeight: '600' },
  emptyText: { color: '#BDBDBD', fontStyle: 'italic', textAlign: 'center', paddingVertical: 16 },
  friendItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 10,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
    gap: 12,
  },
  friendAvatar: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: PRIMARY,
    justifyContent: 'center',
    alignItems: 'center',
  },
  friendAvatarText: { fontSize: 18, fontWeight: 'bold', color: '#fff' },
  friendName: { fontSize: 15, fontWeight: '500', color: '#212121' },
  friendEmail: { fontSize: 12, color: '#757575' },
  leaderItem: {
    flexDirection: 'row',
    alignItems: 'center',
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: '#F5F5F5',
    gap: 12,
  },
  leaderRank: { fontSize: 18, fontWeight: 'bold', color: PRIMARY, width: 36 },
  leaderName: { fontSize: 15, fontWeight: '500', color: '#212121' },
  leaderDetail: { fontSize: 12, color: '#757575', marginTop: 2 },
});
