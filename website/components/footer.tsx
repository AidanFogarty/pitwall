import Link from "next/link";

export function Footer() {
  return (
    <footer className="mt-20">
      <div className="-mx-6 border-t border-border" />
      <div className="pt-12">
        <div className="text-sm text-muted-foreground flex items-center justify-between">
          <p>Pitwall</p>
          <div className="flex gap-6">
            <Link
              href="https://github.com/AidanFogarty/pitwall"
              className="hover:text-foreground transition"
            >
              GitHub
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
