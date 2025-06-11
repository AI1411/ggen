# Frontend Project

A modern React application built with TypeScript, Tailwind CSS, and other cutting-edge technologies.

## Tech Stack

- **React**: A JavaScript library for building user interfaces
- **TypeScript**: Typed JavaScript at scale
- **Vite**: Next generation frontend tooling
- **Tailwind CSS**: A utility-first CSS framework
- **shadcn/ui**: Beautifully designed components built with Radix UI and Tailwind CSS
- **TanStack Query**: Powerful asynchronous state management for fetching, caching, and updating data
- **TanStack Router**: Type-safe routing for React applications
- **Orval**: Generate API clients from OpenAPI specifications
- **Biome**: Fast and reliable linter and formatter
- **Zod**: TypeScript-first schema validation with static type inference
- **Zustand**: A small, fast and scalable state-management solution

## Project Structure

The project follows a feature-based architecture:

```
src/
├── api/                  # API related code
│   └── mutator/          # Custom API client configuration
├── components/           # Shared components
│   └── ui/               # UI components (shadcn/ui)
├── features/             # Feature modules
│   └── auth/             # Authentication feature
│       ├── components/   # Feature-specific components
│       ├── hooks/        # Feature-specific hooks
│       └── types/        # Feature-specific types
├── lib/                  # Utility functions and shared code
│   ├── validations/      # Zod schemas
│   └── utils.ts          # Utility functions
├── routes/               # Application routes
├── store/                # Global state management
├── App.tsx               # Main application component
└── main.tsx              # Application entry point
```

## Getting Started

### Prerequisites

- Node.js (v16 or later)
- npm or yarn

### Installation

1. Clone the repository
2. Install dependencies:

```bash
cd frontend
npm install
# or
yarn
```

### Development

Start the development server:

```bash
npm run dev
# or
yarn dev
```

### Building for Production

Build the application for production:

```bash
npm run build
# or
yarn build
```

### Linting and Formatting

Lint the code:

```bash
npm run lint
# or
yarn lint
```

Format the code:

```bash
npm run format
# or
yarn format
```

## Adding New Features

To add a new feature:

1. Create a new directory in `src/features/`
2. Add components, hooks, and types specific to the feature
3. Create routes in `src/routes/` if needed
4. Import and use the feature in your application

## API Integration

1. Place your OpenAPI specification in `openapi.yaml` at the root
2. Run Orval to generate API clients:

```bash
npm run generate-api
# or
yarn generate-api
```

## Styling Components

The project uses Tailwind CSS for styling. For complex components, consider using shadcn/ui or creating custom components in the `components/ui` directory.

## State Management

- Use Zustand for global state management
- Use TanStack Query for server state
- Use React's built-in state management (useState, useReducer) for component-local state