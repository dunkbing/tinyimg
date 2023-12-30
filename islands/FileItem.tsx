import { useEffect, useState } from "preact/hooks";
import IconFileDownload from "tabler_icons_tsx/file-download.tsx";

import { truncateString } from "@/utils/strings.ts";
import { Loader } from "@/components/Loader.tsx";

type FileItemProps = { file: File; uploadUrl: string };

type FileResponse = {
  newSize: number;
  imageUrl: string;
  savedBytes: number;
  time: number;
};

type FileItemState = {
  savedPercentage: number;
  newSize: number;
  status: "gt" | "lt" | "eq";
};

const FileItem = ({ file, uploadUrl }: FileItemProps) => {
  const imageType = file.type?.split("/")?.[1];
  const [state, setState] = useState<FileItemState | null>(null);
  const [compressing, setCompressing] = useState(false);
  const statusColor = {
    gt: "text-red-500",
    lt: "text-green-500",
    eq: "text-gray-500",
  }[state?.status || "eq"];
  const prefixSign = {
    gt: "+",
    lt: "-",
    eq: "",
  }[state?.status || "eq"];

  useEffect(() => {
    const compressFile = async () => {
      setCompressing(true);
      const formData = new FormData();
      formData.append("file", file);
      try {
        const res = await fetch(uploadUrl, {
          method: "POST",
          body: formData,
        });
        const fr = await res.json() as FileResponse;
        const savedPercentage = (
          (fr.savedBytes / file.size) *
          100
        ).toFixed(2);
        const status = fr.newSize > file.size
          ? "gt"
          : fr.newSize < file.size
          ? "lt"
          : "eq";
        setState({
          savedPercentage: Number(savedPercentage),
          newSize: fr.newSize,
          status,
        });
      } catch (error) {
        console.error(error);
      } finally {
        setCompressing(false);
      }
    };

    compressFile();
  }, []);

  return (
    <div className="flex flex-row items-center justify-between bg-white p-3 rounded-lg shadow-md my-2 w-full">
      <div className="flex flex-row items-center space-x-4">
        <div className="flex-shrink-0">
          <img
            src={URL.createObjectURL(file)}
            alt={file.name}
            className="w-20 h-20 object-cover rounded-md"
          />
        </div>
        <div className="flex flex-col">
          <p className="text-lg font-semibold">
            {truncateString(file.name, 25)}
          </p>
          <div className="flex flex-row text-gray-500">
            <p>{imageType}</p>
            <span className="mx-2">•</span>
            <p>{(file.size / 1024).toFixed(2)} KB</p>
          </div>
        </div>
      </div>
      {compressing ? <Loader /> : (
        <div class="flex flex-row items-center">
          <div className="flex flex-col p-2 text-sm">
            <p className="text-gray-500">new size</p>
            <p className={`font-bold ${statusColor}`}>saved</p>
          </div>
          <div className="flex flex-col bg-blue-100 rounded-l-md p-2 text-sm">
            <p className="text-gray-500">
              {prefixSign}
              {((state?.newSize || 0) / 1024).toFixed(2)} KB
            </p>
            <p
              className={`font-bold ${statusColor}`}
            >
              {state?.savedPercentage}%
            </p>
          </div>
          <button className="flex flex-col items-center justify-center bg-blue-200 text-blue-500 px-3 py-1.5 rounded-r-md hover:bg-blue-300 focus:outline-none focus:ring focus:border-blue-300">
            <IconFileDownload className="w-6 h-6" />
            <p className="text-sm font-bold">jpeg</p>
          </button>
        </div>
      )}
    </div>
  );
};

export default FileItem;
