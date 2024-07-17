import { Handlers, PageProps } from "$fresh/server.ts";

import Form from "@/islands/Form.tsx";
import Head from "@/components/Head.tsx";
import { StatsView } from "@/components/Stats.tsx";
import { FAQ } from "@/components/FAQ.tsx";
import config from "@/utils/config.ts";
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
      <div class="mt-4" />
      <script
        type="text/javascript"
        src="https://cdnjs.buymeacoffee.com/1.0.0/button.prod.min.js"
        data-name="bmc-button"
        data-slug="dangbinh48a"
        data-color="#FFDD00"
        data-emoji="☕"
        data-font="Cookie"
        data-text="Buy me a coffee"
        data-outline-color="#000000"
        data-font-color="#000000"
        data-coffee-color="#ffffff"
      />
      <Form uploadUrl={uploadUrl} downloadUrl={downloadUrl} />
      <div class="mt-2" />
      <StatsView {...(data as any)} />
      <FAQ />
    </div>
  );
}
