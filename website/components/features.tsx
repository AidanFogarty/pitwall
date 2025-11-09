const features = [
  {
    title: "Timing Tower",
    description:
      "Real-time race positions, sector times, gap analysis, and tyre strategy tracking",
  },
  {
    title: "Race Control",
    description:
      "Live feed of official race control messages and track status updates",
  },
  {
    title: "Session Intelligence",
    description:
      "Comprehensive session data including weather conditions and track information",
  },
];

export function Features() {
  return (
    <div className="mt-10">
      <div className="-mx-6 border-t border-border" />
      <div className="pt-12">
        <h2 className="text-3xl font-bold mb-4">Features</h2>
        <p className="text-muted-foreground text-lg mb-12 leading-relaxed max-w-2xl">
          Live race data and timing information delivered to your terminal.
        </p>

        <div className="space-y-6">
          {features.map((feature, index) => (
            <div key={index} className="flex gap-4">
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
