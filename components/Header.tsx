import { SITE_BAR_STYLES, SITE_NAME } from "@/utils/constants.ts";

export interface HeaderProps {
  /**
   * URL of the current page. This is used for highlighting the currently
   * active page in navigation.
   */
  url: URL;
}

export default function Header(props: HeaderProps) {
  const NAV_ITEM = "text-green-900 px-3 py-4 sm:py-2";

  return (
    <header
      class={`${SITE_BAR_STYLES} my-2 flex-col sm:flex-row`}
    >
      <div class="flex justify-between items-center">
        <a href="/" class="shrink-0">
          <img
            src="/cover.png"
            alt={SITE_NAME + " logo"}
            width={190}
            height={100}
          />
        </a>
      </div>
      <nav
        class={"font-semibold flex flex-col gap-x-4 divide-y divide-solid sm:flex sm:items-center sm:flex-row sm:divide-y-0"}
      >
        <a
          href="/about"
          className={NAV_ITEM}
        >
          About
        </a>
        <a
          href="mailto:bing@db99.dev"
          className={NAV_ITEM}
        >
          Feedback
        </a>
      </nav>
    </header>
  );
}
