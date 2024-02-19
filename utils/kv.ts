export const kv = await Deno.openKv();

export const statsEntryKey = ["stats"];

export type Stats = {
  totalFiles: number;
  totalSize: number;
};

const statsEntry = await kv.get(statsEntryKey);

if (!statsEntry.value) {
  await kv.set(statsEntryKey, {
    totalFiles: 0,
    totalSize: 0,
  });
}
