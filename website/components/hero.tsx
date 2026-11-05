import Link from "next/link";
import { Button } from "@/components/ui/button";
import { getLatestVersion } from "@/lib/release";

const REPO_URL = "https://github.com/AidanFogarty/pitwall";

export async function Hero() {
  const version = await getLatestVersion();

  return (
    <div className="max-w-2xl md:max-w-5xl">
      <p className="text-muted-foreground text-sm mb-8">
        <Link
          href={`${REPO_URL}/releases`}
          className="hover:text-foreground transition"
        >
          {version ? `${version} — see what's new` : "See what's new"}
        </Link>
      </p>

      <h1 className="text-5xl font-bold mb-6 leading-tight">
        A F1 live timing client built for the terminal
      </h1>

      <p className="text-muted-foreground text-lg mb-8 leading-relaxed">
        Live timing, sector times, tyre strategy and race control, right in your
        terminal. Open source and fast.
      </p>

      <Button asChild size="lg" className="mb-16">
        <Link href={`${REPO_URL}#readme`}>
          Read docs
          <span>→</span>
        </Link>
      </Button>
    </div>
  );
}
