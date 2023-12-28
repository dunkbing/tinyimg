import { useMemo, useState } from "preact/hooks";
import { Button } from "@/components/Button.tsx";
import { downloadFile } from "@/utils/http.ts";
import { truncateString } from "@/utils/strings.ts";
import IconFileDownload from "tabler_icons_tsx/file-download.tsx";

type FileItemProps = { file: File };

const FileItem = ({ file }: FileItemProps) => {
  return (
    <div className="flex flex-row items-center justify-between bg-white p-4 rounded-lg shadow-md mb-4 w-full">
      <div className="flex flex-row items-center space-x-4">
        <div className="flex-shrink-0">
          <img
            src={URL.createObjectURL(file)}
            alt={file.name}
            className="w-12 h-12 object-cover rounded-md"
          />
        </div>
        <div className="flex flex-col">
          <p className="text-lg font-semibold">
            {truncateString(file.name, 25)}
          </p>
          <div className="flex flex-row text-gray-500">
            <p>{file.type}</p>
            <span className="mx-2">•</span>
            <p>{(file.size / 1024).toFixed(2)} KB</p>
          </div>
        </div>
      </div>
      <div class="flex flex-row items-center space-x-4 bg-blue-100 rounded-md">
        <div className="flex flex-col pl-3">
          <p className="text-red-500">-7%</p>
          <p className="text-gray-500">64 KB</p>
        </div>
        <button className="flex flex-col items-center justify-center bg-blue-200 text-blue-500 px-3 py-1.5 rounded-md hover:bg-blue-300 focus:outline-none focus:ring focus:border-blue-300">
          <IconFileDownload className="w-6 h-6" />
          {/* Adjust the size as needed */}
          <p className="text-sm font-bold">jpeg</p>
        </button>
      </div>
    </div>
  );
};

export default FileItem;
