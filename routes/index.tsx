import { Handlers, PageProps } from "$fresh/server.ts";

import Form from "@/islands/Form.tsx";
import Head from "@/components/Head.tsx";
import config from "@/utils/config.ts";
import { Promote } from "@/components/Promote.tsx";
import { StatsView } from "@/components/Stats.tsx";
import { kv, statsEntryKey } from "@/utils/kv.ts";

const uploadUrl = `${config.apiUrl}/upload`;
const downloadUrl = `${config.apiUrl}/download-all`;

export const handler: Handlers<unknown> = {
  async GET(_req, ctx) {
    try {
      const statsEntry = await kv.get(statsEntryKey);

      return ctx.render(statsEntry.value);
    } catch (error) {
      return ctx.render({});
    }
  },
};

export default function Home(ctx: PageProps<unknown>) {
  const data = ctx.data;

  return (
    <div class="flex flex-col justify-center items-center">
      <Head href={ctx.url.href}>
        <link
          as="fetch"
          crossOrigin="anonymous"
          href={ctx.url.href}
          rel="preload"
        />
      </Head>
      <Form uploadUrl={uploadUrl} downloadUrl={downloadUrl} />
      <StatsView {...data as any} />
      <div class="mt-4" />
      <Promote />
    </div>
  );
}
