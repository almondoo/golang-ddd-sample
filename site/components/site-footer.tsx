export function SiteFooter() {
  return (
    <footer className="bg-foreground px-6 py-14 text-center text-background">
      <p>
        <a
          href="https://github.com/almondoo/golang-ddd-sample"
          target="_blank"
          rel="noopener noreferrer"
          className="font-bold text-secondary underline decoration-secondary/50 underline-offset-4 hover:decoration-secondary"
        >
          github.com/almondoo/golang-ddd-sample
        </a>
      </p>
      <p className="mt-4 text-sm text-background/70">
        このサイトは絵本です。本当の教科書は、コードそのもの。
      </p>
    </footer>
  );
}
