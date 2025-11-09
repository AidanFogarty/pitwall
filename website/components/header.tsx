import Link from "next/link";

export function Header() {
  return (
    <header className="border-b border-border">
      <nav className="px-6 py-6 flex items-center justify-between">
        <Link href="/" className="text-2xl font-semibold tracking-wider">
          pitwall
        </Link>
        <div className="flex items-center gap-4 text-sm">
          <Link
            href="https://github.com/AidanFogarty/pitwall"
            className="text-muted-foreground hover:text-foreground transition"
          >
            GitHub
          </Link>
        </div>
      </nav>
    </header>
  );
}
