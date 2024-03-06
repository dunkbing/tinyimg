import { PageProps } from "$fresh/server.ts";

import Head from "@/components/Head.tsx";
import { About } from "@/components/About.tsx";

export default function (ctx: PageProps) {
  return (
    <div class="flex flex-col justify-center items-center">
      <Head href={ctx.url.href}>
        <link
          as="fetch"
          crossOrigin="anonymous"
          href={ctx.url.href}
          rel="preload"
        />
        <style>
          {`
          pre code {
  background-color: #eee;
  border: 1px solid #999;
  display: block;
  padding: 10px;
}`}
        </style>
      </Head>
      <About />
    </div>
  );
}
