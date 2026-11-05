import path from "node:path";
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  turbopack: {
    // Pin the root to this folder. Without it Next picks up a stray
    // package-lock.json further up the tree and warns about it.
    root: path.join(__dirname),
  },
};

export default nextConfig;
