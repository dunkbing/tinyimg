import { FreshContext } from "$fresh/server.ts";
import { kv, Stats, statsEntryKey } from "@/utils/kv.ts";

export const handler = {
  async GET(_req: Request, _ctx: FreshContext): Promise<Response> {
    const statsEntry = await kv.get(statsEntryKey);
    return new Response(JSON.stringify(statsEntry.value));
  },
  async POST(req: Request, _ctx: FreshContext): Promise<Response> {
    const { totalFiles, totalSize } = (await req.json()) as Stats;
    const statsEntry = await kv.get(statsEntryKey);
    const currentStats = statsEntry.value as Stats;
    await kv.set(statsEntryKey, {
      totalFiles: currentStats.totalFiles + totalFiles,
      totalSize: currentStats.totalSize + totalSize,
    });

    return new Response(JSON.stringify({ message: "ok" }));
  },
};
