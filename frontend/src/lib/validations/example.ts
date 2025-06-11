import { z } from 'zod';

// User schema
export const userSchema = z.object({
  id: z.string().uuid(),
  name: z.string().min(2, { message: 'Name must be at least 2 characters' }),
  email: z.string().email({ message: 'Invalid email address' }),
  role: z.enum(['admin', 'user', 'guest']),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime().optional(),
});

// Infer the TypeScript type from the schema
export type User = z.infer<typeof userSchema>;

// Login form schema
export const loginFormSchema = z.object({
  email: z.string().email({ message: 'Invalid email address' }),
  password: z.string().min(8, { message: 'Password must be at least 8 characters' }),
  rememberMe: z.boolean().default(false),
});

export type LoginFormValues = z.infer<typeof loginFormSchema>;

// Registration form schema
export const registrationFormSchema = z.object({
  name: z.string().min(2, { message: 'Name must be at least 2 characters' }),
  email: z.string().email({ message: 'Invalid email address' }),
  password: z.string().min(8, { message: 'Password must be at least 8 characters' }),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
});

export type RegistrationFormValues = z.infer<typeof registrationFormSchema>;