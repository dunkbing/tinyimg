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

const FileItem = ({ file, uploadUrl, formats, filesSig }: FileItemProps) => {
  const imageType = file.type?.split("/")?.[1];
  const [state, setState] = useState<FileItemState[]>([]);
  const [compressing, setCompressing] = useState(false);

  useEffect(() => {
    const compressFile = async () => {
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
        setState(newState);
      } catch (error) {
        console.error(error);
      } finally {
        setCompressing(false);
      }
    };

    compressFile();
  }, [file, formats, uploadUrl]);

  return (
    <div className="flex flex-col items-center bg-white p-3 rounded-lg shadow-md my-2">
      <div className="flex flex-col lg:flex-row items-center space-y-4 sm:space-y-0 sm:space-x-4">
        <div className="flex-shrink-0">
          <img
            src={URL.createObjectURL(file)}
            alt={file.name}
            className="w-32 h-32 object-cover rounded-lg"
          />
        </div>
        <div className="flex flex-col space-y-1">
          <p className="text-sm font-semibold line-clamp-1">
            {truncateString(file.name, 30)}
          </p>
          <div className="flex flex-row text-gray-500">
            <p>{imageType}</p>
            <span className="mx-2">•</span>
            <p>{(file.size / 1024).toFixed(2)} KB</p>
          </div>
          {compressing
            ? <Loader />
            : (
              <div class="flex flex-col lg:flex-row items-center space-y-1 md:space-x-1">
                <div className="flex flex-col py-2 text-sm hidden lg:block">
                  <p className="text-gray-500">new size</p>
                  <p className={`font-bold`}>saved</p>
                </div>
                {state.map((s, index) => (
                  <div class="flex flex-row" key={index}>
                    <div className="flex flex-col items-center bg-blue-100 rounded-l-md py-2 text-sm ml-1 w-24">
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
