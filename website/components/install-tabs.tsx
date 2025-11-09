"use client";

import { useState } from "react";
import { Copy, Check } from "lucide-react";
import { Button } from "@/components/ui/button";

const installMethods: Record<string, string> = {
  brew: "brew install AidanFogarty/tap/pitwall",
  go: "go install github.com/AidanFogarty/pitwall@latest",
};

export function InstallTabs() {
  const [activeTab, setActiveTab] = useState("brew");
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(installMethods[activeTab]);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="max-w-2xl md:max-w-5xl border border-border rounded-lg overflow-hidden">
      <div className="flex border-b border-border bg-secondary/20">
        {Object.keys(installMethods).map((method) => (
          <Button
            key={method}
            onClick={() => setActiveTab(method)}
            variant="ghost"
            size="sm"
            className={`rounded-none border-0 relative ${
              activeTab === method ? "text-foreground" : "text-muted-foreground"
            }`}
          >
            {method}
            {activeTab === method && (
              <div className="absolute bottom-0 left-0 right-0 h-0.5 bg-foreground"></div>
            )}
          </Button>
        ))}
      </div>

      <div className="p-6 bg-background flex items-center justify-between">
        <code className="text-foreground text-sm">
          {installMethods[activeTab]}
        </code>
        <div className="flex items-center gap-3">
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={handleCopy}
            title={copied ? "Copied!" : "Copy to clipboard"}
            className={`transition-all duration-200 ${
              copied ? "text-primary" : ""
            }`}
          >
            {copied ? (
              <Check
                size={18}
                className="animate-in fade-in zoom-in duration-200"
              />
            ) : (
              <Copy size={18} />
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}
