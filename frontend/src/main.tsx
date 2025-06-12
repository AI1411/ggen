import { queryClient } from '@/lib/tanstack-query';
import { QueryClientProvider } from '@tanstack/react-query';
import ReactDOM from 'react-dom/client';
import './index.css';
import { RouterProvider } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/router-devtools';
import { StrictMode } from 'react';
import AxiosProvider from './app/providers/AxiosProvider.tsx';
import { router } from './app/route.tsx';

// アプリケーションのエントリーポイント
ReactDOM.createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <AxiosProvider>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
        <TanStackRouterDevtools router={router} />
      </QueryClientProvider>
    </AxiosProvider>
  </StrictMode>,
);
