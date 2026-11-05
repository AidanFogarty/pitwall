const RELEASES_API =
  "https://api.github.com/repos/AidanFogarty/pitwall/releases/latest";

/**
 * Latest published Pitwall release tag (e.g. "v0.0.7"), or null if GitHub is
 * unreachable. Revalidated hourly so the site never needs a manual version bump.
 */
export async function getLatestVersion(): Promise<string | null> {
  try {
    const res = await fetch(RELEASES_API, {
      headers: { Accept: "application/vnd.github+json" },
      next: { revalidate: 3600 },
    });
    if (!res.ok) return null;
    const data: { tag_name?: string } = await res.json();
    return data.tag_name ?? null;
  } catch {
    return null;
  }
}
