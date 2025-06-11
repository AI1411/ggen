import { createRootRoute, Outlet } from '@tanstack/react-router';

export const rootRoute = createRootRoute({
  component: () => (
    <div className="min-h-screen bg-background">
      <main className="container mx-auto py-8">
        <Outlet />
      </main>
    </div>
  ),
});