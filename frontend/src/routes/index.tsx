import { createFileRoute } from '@tanstack/react-router';

// 明示的に親ルートからのパスを指定
export const indexRoute = createFileRoute('/')({
  component: Index,
});

function Index() {
  return (
    <div className="flex flex-col items-center justify-center gap-6">
      <h1 className="text-4xl font-bold">Welcome to your new app</h1>
      <p className="text-xl text-muted-foreground">
        This is a starter template with React, TypeScript, Tailwind CSS, shadcn/ui, TanStack Query,
        and TanStack Router.
      </p>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-8">
        <FeatureCard
          title="React + TypeScript"
          description="Modern frontend development with type safety."
        />
        <FeatureCard
          title="Tailwind CSS + shadcn/ui"
          description="Beautiful, accessible components with utility-first CSS."
        />
        <FeatureCard
          title="TanStack Query"
          description="Powerful data fetching and caching for your API."
        />
        <FeatureCard
          title="TanStack Router"
          description="Type-safe routing for your React application."
        />
        <FeatureCard
          title="Orval"
          description="Generate API clients from OpenAPI specifications."
        />
        <FeatureCard title="Zod + Zustand" description="Schema validation and state management." />
      </div>
    </div>
  );
}

function FeatureCard({ title, description }: { title: string; description: string }) {
  return (
    <div className="bg-card rounded-lg shadow-sm p-6 border border-border">
      <h3 className="text-xl font-semibold mb-2">{title}</h3>
      <p className="text-muted-foreground">{description}</p>
    </div>
  );
}
