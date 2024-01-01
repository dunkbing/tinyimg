import { getFileName } from "@/utils/strings.ts";

export function downloadFile(url: string, name?: string) {
  const filename = name || getFileName(url);
  const a = document.createElement("a");
  a.style.display = "none";
  document.body.appendChild(a);
  a.href = url;
  a.download = filename as string;
  a.click();
  document.body.removeChild(a);
}
