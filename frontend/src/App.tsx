import { Link } from '@tanstack/react-router';

export const Hello = () => {
  return (
    <div>
      <h1>Hello World</h1>
      <p>Click on the links above to see the code splitting in action.</p>
      <Link to="/">Go to Home</Link>
    </div>
  );
};
