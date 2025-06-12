import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useAuthStore } from '@/store/auth-store';
import { loginFormSchema, type LoginFormValues } from '@/lib/validations/example';

export function Login() {
  const [isLoading, setIsLoading] = useState(false);
  const login = useAuthStore((state) => state.login);

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginFormSchema),
    defaultValues: {
      email: '',
      password: '',
      rememberMe: false,
    },
  });

  async function onSubmit(data: LoginFormValues) {
    setIsLoading(true);
    try {
      // In a real app, you would call an API here
      console.log('Login form submitted', data);

      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1000));

      // Mock successful login
      login(
        {
          id: '123e4567-e89b-12d3-a456-426614174000',
          name: 'John Doe',
          email: data.email,
          role: 'user',
          createdAt: new Date().toISOString(),
        },
        'mock-jwt-token',
      );
    } catch (error) {
      console.error('Login failed', error);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="space-y-2 text-center">
        <h1 className="text-3xl font-bold">Login</h1>
        <p className="text-gray-500 dark:text-gray-400">
          Enter your credentials to access your account
        </p>
      </div>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="space-y-2">
          <label htmlFor="email" className="text-sm font-medium">
            Email
          </label>
          <input
            id="email"
            type="email"
            placeholder="name@example.com"
            className="w-full px-3 py-2 border rounded-md"
            {...form.register('email')}
          />
          {form.formState.errors.email && (
            <p className="text-sm text-red-500">{form.formState.errors.email.message}</p>
          )}
        </div>
        <div className="space-y-2">
          <label htmlFor="password" className="text-sm font-medium">
            Password
          </label>
          <input
            id="password"
            type="password"
            placeholder="••••••••"
            className="w-full px-3 py-2 border rounded-md"
            {...form.register('password')}
          />
          {form.formState.errors.password && (
            <p className="text-sm text-red-500">{form.formState.errors.password.message}</p>
          )}
        </div>
        <div className="flex items-center space-x-2">
          <input
            id="rememberMe"
            type="checkbox"
            className="h-4 w-4"
            {...form.register('rememberMe')}
          />
          <label htmlFor="rememberMe" className="text-sm font-medium">
            Remember me
          </label>
        </div>
        <button
          type="submit"
          className="w-full py-2 px-4 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
          disabled={isLoading}
        >
          {isLoading ? 'Loading...' : 'Login'}
        </button>
      </form>
    </div>
  );
}
