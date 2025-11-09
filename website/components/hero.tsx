import Link from "next/link";
import { Button } from "@/components/ui/button";

export function Hero() {
  return (
    <div className="max-w-2xl md:max-w-5xl">
      <p className="text-muted-foreground text-sm mb-8">
        <Link href={"https://github.com/AidanFogarty/pitwall"}>
          See what&apos;s new!
        </Link>
      </p>

      <h1 className="text-5xl font-bold mb-6 leading-tight">
        A F1 live timing client built for the terminal
      </h1>

      <p className="text-muted-foreground text-lg mb-8 leading-relaxed">
        Pitwall brings real-time Formula 1 data directly to your terminal. Track
        live timing, driver standings, and race updates with a lightweight,
        open-source client designed for speed.
      </p>

      <Button asChild size="lg" className="mb-16">
        <Link href="https://github.com/AidanFogarty/pitwall">
          Read docs
          <span>→</span>
        </Link>
      </Button>
    </div>
  );
}
