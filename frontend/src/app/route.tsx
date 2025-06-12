import { Layout } from '@/components/layout/Header.tsx';
import { Outlet, createRootRoute, createRoute, createRouter } from '@tanstack/react-router';
import { Hello } from '../App.tsx';
import { TopPage } from '../pages/TopPage.tsx';

const rootRoute = createRootRoute({});

// TOP page route
const topRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: () => (
    <Layout>
      <TopPage />
    </Layout>
  ),
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/hello',
  component: () => (
    <Layout>
      <Outlet />
    </Layout>
  ),
});

const helloRoute = createRoute({
  getParentRoute: () => indexRoute,
  path: '/foo',
  component: () => <Hello />,
});

const routeTree = rootRoute.addChildren([topRoute, indexRoute.addChildren([helloRoute])]);

export const router = createRouter({ routeTree });

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
