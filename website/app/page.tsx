import { Header } from "@/components/header";
import { Hero } from "@/components/hero";
import { InstallTabs } from "@/components/install-tabs";
import { ThemedImage } from "@/components/themed-image";
import { Features } from "@/components/features";
import { ComingSoon } from "@/components/coming-soon";
import { Footer } from "@/components/footer";
import timingTowerLight from "@/public/images/timing-tower-light-mode.png";
import timingTowerDark from "@/public/images/timing-tower-dark-mode.png";
import qualifyingModeLight from "@/public/images/qualifying-mode-light-mode.png";
import qualifyingModeDark from "@/public/images/qualifying-mode-dark-mode.png";

export default function Home() {
  return (
    <div className="min-h-screen bg-background text-foreground font-mono flex flex-col">
      <div className="max-w-5xl mx-auto border-l border-r border-border flex-1 flex flex-col">
        <Header />

        <main className="px-6 py-20">
          <Hero />
          <InstallTabs />

          <div className="mt-12">
            <ThemedImage
              lightSrc={timingTowerLight}
              darkSrc={timingTowerDark}
              alt="Pitwall Timing Tower"
              className="rounded-lg border border-border"
            />
          </div>

          <Features />

          <div className="mt-12">
            <ThemedImage
              lightSrc={qualifyingModeLight}
              darkSrc={qualifyingModeDark}
              alt="Pitwall Qualifying Mode"
              className="rounded-lg border border-border"
            />
          </div>

          <ComingSoon />

          <Footer />
        </main>
      </div>
    </div>
  );
}
