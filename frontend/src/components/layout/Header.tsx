export const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div>
      <header>
        <h1>I'm a header</h1>
      </header>
      {children}
    </div>
  );
};
