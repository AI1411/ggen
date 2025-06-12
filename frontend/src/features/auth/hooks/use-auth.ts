import { useCallback } from 'react';
import { useAuthStore } from '@/store/auth-store';
import { User } from '@/lib/validations/example';

export function useAuth() {
  const { user, token, isAuthenticated, login, logout, updateUser } = useAuthStore();

  const handleLogin = useCallback(
    async (email: string, password: string): Promise<boolean> => {
      try {
        // In a real app, you would call an API here
        console.log('Logging in with', email, password);

        // Simulate API call
        await new Promise((resolve) => setTimeout(resolve, 1000));

        // Mock successful login
        const mockUser: User = {
          id: '123e4567-e89b-12d3-a456-426614174000',
          name: 'John Doe',
          email,
          role: 'user',
          createdAt: new Date().toISOString(),
        };

        login(mockUser, 'mock-jwt-token');
        return true;
      } catch (error) {
        console.error('Login failed', error);
        return false;
      }
    },
    [login],
  );

  const handleLogout = useCallback(() => {
    logout();
  }, [logout]);

  return {
    user,
    token,
    isAuthenticated,
    login: handleLogin,
    logout: handleLogout,
    updateUser,
  };
}
