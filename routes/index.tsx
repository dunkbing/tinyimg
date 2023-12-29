import { PageProps } from "$fresh/server.ts";

import Form from "@/islands/Form.tsx";
import Head from "@/components/Head.tsx";
import config from "@/utils/config.ts";

const uploadUrl = `${config.apiUrl}/upload`;

export default function Home(ctx: PageProps) {
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
      <Form uploadUrl={uploadUrl} />
    </div>
  );
}
