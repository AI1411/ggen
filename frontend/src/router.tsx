import { createRouter, RouterProvider } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/router-devtools';
import { StrictMode } from 'react';

// Import the generated route tree
import { routeTree } from './routeTree.gen';

const router = createRouter({ routeTree });

// Register the router for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

// Router component
export function Router() {
  return (
    <StrictMode>
      <RouterProvider router={router} />
      {import.meta.env.DEV && <TanStackRouterDevtools router={router} />}
    </StrictMode>
  );
}