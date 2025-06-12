import { Layout } from '@/components/layout/Header.tsx';
import { Outlet, createRootRoute, createRoute, createRouter } from '@tanstack/react-router';
import { Hello } from '../App.tsx';

const rootRoute = createRootRoute({});

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

const routeTree = rootRoute.addChildren([indexRoute.addChildren([helloRoute])]);

export const router = createRouter({ routeTree });

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
