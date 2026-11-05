import Image from "next/image";
import type { ReactNode } from "react";
import { ThemedImage } from "@/components/themed-image";
import qualifyingModeDark from "@/public/images/qualifying-mode-dark-mode.png";
import qualifyingModeLight from "@/public/images/qualifying-mode-light-mode.png";
import raceControlDark from "@/public/images/race-control-dark-mode.png";
import timingTowerDark from "@/public/images/timing-tower-dark-mode.png";
import timingTowerLight from "@/public/images/timing-tower-light-mode.png";

const imageClassName = "rounded-lg border border-border";

interface Feature {
  title: string;
  description: string;
  image: ReactNode;
}

const features: Feature[] = [
  {
    title: "Timing Tower",
    description:
      "Live positions, gaps and lap times, with colour-coded sectors and tyre compounds.",
    image: (
      <ThemedImage
        lightSrc={timingTowerLight}
        darkSrc={timingTowerDark}
        alt="Pitwall Timing Tower"
        className={imageClassName}
      />
    ),
  },
  {
    title: "Qualifying & Practice",
    description: "Follow every lap with live sector and mini sector progress.",
    image: (
      <ThemedImage
        lightSrc={qualifyingModeLight}
        darkSrc={qualifyingModeDark}
        alt="Pitwall Qualifying Mode"
        className={imageClassName}
      />
    ),
  },
  {
    title: "Race Control",
    description:
      "Every official message in a full-screen view, with flags and penalties highlighted.",
    image: (
      <Image
        src={raceControlDark}
        alt="Pitwall Race Control screen"
        className={imageClassName}
        style={{ width: "100%", height: "auto" }}
      />
    ),
  },
];

export function Features() {
  return (
    <div className="mt-10">
      <div className="-mx-6 border-t border-border" />
      <div className="pt-12">
        <h2 className="text-3xl font-bold mb-12">Features</h2>

        <div className="space-y-16 md:space-y-24">
          {features.map((feature, index) => (
            <div
              key={feature.title}
              className="grid items-center gap-8 md:grid-cols-5 md:gap-12"
            >
              <div
                className={`md:col-span-2 ${
                  index % 2 === 1 ? "md:order-2" : ""
                }`}
              >
                <h3 className="text-xl font-semibold mb-3">{feature.title}</h3>
                <p className="text-muted-foreground leading-relaxed">
                  {feature.description}
                </p>
              </div>

              <div
                className={`md:col-span-3 ${
                  index % 2 === 1 ? "md:order-1" : ""
                }`}
              >
                {feature.image}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
