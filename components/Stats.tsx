function formatNumber(number: number) {
  if (number >= 1000000) {
    return (number >= 2000000) ? Math.floor(number / 1000000) + "M+" : "1M+";
  }
  return number.toLocaleString("en-US");
}

function formatBytes(bytes: number) {
  if (bytes >= 1000000) {
    return (bytes >= 2000000) ? Math.floor(bytes / 1000000).toLocaleString() + "Mb" : "1Mb";
  }
  if (bytes >= 1000) {
    return (bytes >= 2000) ? Math.floor(bytes / 1000).toLocaleString() + "Kb" : "1Kb";
  }
  return bytes.toLocaleString() + "Bytes";
}

export function StatsView(props: { totalFiles: number; totalSize: number }) {
  const { totalFiles = 10000, totalSize = 1000000000 } = props;

  return (
    <div
      id="about"
      class="flex flex-col items-center justify-center mt-2 px-5"
    >
      <p class="text-green-800 text-center text-xl mb-4">
        <b>{formatNumber(totalFiles)}</b> files converted <br />
        <b>{formatBytes(totalSize)}</b> saved
      </p>
    </div>
  );
}
