import { useEffect, useState } from "preact/hooks";
import IconFileDownload from "tabler_icons_tsx/file-download.tsx";

import { formatPercentage, truncateString } from "@/utils/strings.ts";
import { Loader } from "@/components/Loader.tsx";
import { downloadFile } from "@/utils/http.ts";
import { Signal } from "@preact/signals";

type FileItemProps = {
  file: File;
  uploadUrl: string;
  formats: string[];
  filesSig: Signal<string[]>;
};

type FileResponse = {
  savedBytes: number;
  newSize: number;
  time: number;
  imageUrl: string;
  format: string;
};

type FileItemState = {
  savedPercentage: number;
  newSize: number;
  statusColor: string;
  prefixSign: string;
  format: string;
  imageUrl: string;
};

const updateStats = async (stats: {
  totalFiles: number;
  totalSize: number;
}) => {
  await fetch("/api/stats", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(stats),
  });
};

const maxFileSize = 10 * 1024 * 1024;

const FileItem = ({ file, uploadUrl, formats, filesSig }: FileItemProps) => {
  const imageType = file.type?.split("/")?.[1];
  const [state, setState] = useState<FileItemState[]>([]);
  const [compressing, setCompressing] = useState(false);

  useEffect(() => {
    const compressFile = async () => {
      if (file.size > maxFileSize) {
        return;
      }

      setCompressing(true);
      const formData = new FormData();
      formData.append("file", file);
      formData.append("formats", formats.join(","));
      try {
        const res = await fetch(uploadUrl, {
          method: "POST",
          body: formData,
        });
        const { data, files } = await res.json() as {
          data: FileResponse[];
          errors: string[];
          files: string[];
        };
        filesSig.value = [...filesSig.value, ...files];
        const newState: FileItemState[] = data.map((fr) => {
          const savedPercentage = (fr.savedBytes / file.size) *
            100;
          const statusColor = fr.newSize > file.size
            ? "text-red-500"
            : fr.newSize < file.size
            ? "text-green-500"
            : "text-gray-500";
          const prefixSign = fr.newSize > file.size
            ? "+"
            : fr.newSize < file.size
            ? "-"
            : "";
          return {
            savedPercentage: Number(savedPercentage),
            newSize: fr.newSize,
            statusColor,
            prefixSign,
            format: fr.format,
            imageUrl: fr.imageUrl,
          };
        });
        const savedBytes = data.reduce((acc, fr) => acc + fr.savedBytes, 0);
        void updateStats({
          totalFiles: files.length,
          totalSize: savedBytes,
        });

        setState(newState);
      } catch (error) {
        console.error(error);
      } finally {
        setCompressing(false);
      }
    };

    void compressFile();
  }, [file, formats, uploadUrl]);

  return (
    <div className="flex flex-col items-center px-4 my-2">
      <div className="flex flex-col lg:flex-row items-center justify-center bg-white p-4 md:p-2 rounded-lg shadow-md">
        <div className="hidden lg:flex mx-1.5">
          <img
            src={URL.createObjectURL(file)}
            alt={file.name}
            className="w-24 h-24 object-cover rounded-lg"
          />
        </div>
        <div className="flex flex-col space-y-0.5 ml-1">
          <p className="text-sm font-semibold">
            {truncateString(file.name, 30)}
          </p>
          <div className="flex flex-row text-gray-500 mb-0.5">
            <p>{imageType}</p>
            <span className="mx-2">•</span>
            <p>{(file.size / 1024).toFixed(2)} KB</p>
          </div>
          {file.size > maxFileSize && (
            <span className="text-red-500">
              The file is too large (max 10MB)
            </span>
          )}
          {compressing
            ? <Loader />
            : (
              <div class="flex flex-row flex-wrap items-center">
                {state.map((s, index) => (
                  <div class="flex flex-row m-0.5" key={index}>
                    <div className="flex flex-col items-center bg-blue-100 rounded-l-md py-2 text-sm w-24">
                      <p className="text-gray-500">
                        {((s.newSize || 0) / 1024).toFixed(2)} KB
                      </p>
                      <p className={`font-bold ${s.statusColor}`}>
                        {formatPercentage(s.savedPercentage)}
                      </p>
                    </div>
                    <button
                      className="flex flex-col items-center justify-center w-16 bg-blue-200 text-blue-500 px-3 py-1.5 rounded-r-md hover:bg-blue-300 focus:outline-none focus:ring focus:border-blue-300"
                      onClick={() => {
                        const params = new URL(s.imageUrl).searchParams;
                        const filename = params.get("f");
                        downloadFile(s.imageUrl, filename as string);
                      }}
                    >
                      <IconFileDownload className="w-6 h-6" />
                      <p className="text-sm font-bold">{s.format || "jpg"}</p>
                    </button>
                  </div>
                ))}
              </div>
            )}
        </div>
      </div>
    </div>
  );
};

export default FileItem;
