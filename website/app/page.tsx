import { ComingSoon } from "@/components/coming-soon";
import { Features } from "@/components/features";
import { Footer } from "@/components/footer";
import { Header } from "@/components/header";
import { Hero } from "@/components/hero";
import { InstallTabs } from "@/components/install-tabs";

export default function Home() {
  return (
    <div className="min-h-screen bg-background text-foreground font-mono flex flex-col">
      <div className="max-w-5xl mx-auto border-l border-r border-border flex-1 flex flex-col">
        <Header />

        <main className="px-6 py-20">
          <Hero />
          <InstallTabs />

          <Features />
          <ComingSoon />

          <Footer />
        </main>
      </div>
    </div>
  );
}
