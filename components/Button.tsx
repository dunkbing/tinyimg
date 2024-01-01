import { JSX } from "preact";
import { IS_BROWSER } from "$fresh/runtime.ts";

export function Button(
  props: JSX.HTMLAttributes<HTMLButtonElement> & {
    colorMode?: "primary" | "secondary";
  },
) {
  const { colorMode } = props;

  return (
    <button
      {...props}
      disabled={!IS_BROWSER || props.disabled}
      class={`flex items-center space-x-1 px-3 py-2 rounded-md border(gray-500 2) active:bg-gray-300 disabled:(opacity-50 cursor-not-allowed) text-white ${
        props.class ?? ""
      } ${
        colorMode === "primary"
          ? "bg-blue-600 hover:bg-blue-700"
          : "bg-green-600 hover:bg-green-700"
      }`}
    />
  );
}
