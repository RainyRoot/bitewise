import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { AuthProvider } from '@/hooks/useAuth';

export default function RootLayout() {
  return (
    <AuthProvider>
      <StatusBar style="auto" />
      <Stack
        screenOptions={{
          headerStyle: { backgroundColor: '#4CAF50' },
          headerTintColor: '#fff',
          headerTitleStyle: { fontWeight: 'bold' },
        }}
      >
        <Stack.Screen
          name="(tabs)"
          options={{ headerShown: false }}
        />
        <Stack.Screen
          name="login"
          options={{
            title: 'Anmelden',
            headerShown: false,
          }}
        />
        <Stack.Screen
          name="register"
          options={{
            title: 'Registrieren',
            headerShown: false,
          }}
        />
        <Stack.Screen
          name="recipe/[id]"
          options={{
            title: 'Rezept',
          }}
        />
        <Stack.Screen
          name="scanner"
          options={{
            title: 'Barcode Scanner',
          }}
        />
        <Stack.Screen
          name="achievements"
          options={{
            title: 'Achievements',
          }}
        />
        <Stack.Screen
          name="create-recipe"
          options={{
            title: 'Rezept erstellen',
          }}
        />
        <Stack.Screen
          name="notification-settings"
          options={{
            title: 'Benachrichtigungen',
          }}
        />
        <Stack.Screen
          name="stats"
          options={{
            title: 'Statistiken',
          }}
        />
      </Stack>
    </AuthProvider>
  );
}
