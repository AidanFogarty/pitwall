const comingSoonFeatures = [
  {
    title: "Driver Tracker & Telemetry",
    description: "Live circuit map with driver positions and car telemetry",
  },
  {
    title: "Live Championship Standings",
    description: "Live driver and team standings, updated as points change",
  },
];

export function ComingSoon() {
  return (
    <div className="mt-10">
      <div className="-mx-6 border-t border-border" />
      <div className="pt-12">
        <h2 className="text-3xl font-bold mb-4">Coming Soon</h2>
        <p className="text-muted-foreground text-lg mb-12 leading-relaxed max-w-2xl">
          In development.
        </p>

        <div className="space-y-6">
          {comingSoonFeatures.map((feature) => (
            <div key={feature.title} className="flex gap-4">
              <div className="shrink-0 text-foreground font-bold">&gt;</div>
              <div>
                <h3 className="text-foreground font-semibold mb-1">
                  {feature.title}
                </h3>
                <p className="text-muted-foreground text-sm">
                  {feature.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
