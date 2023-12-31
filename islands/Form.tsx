import { JSX } from "preact";
import { useState } from "preact/hooks";
import { useToaster } from "fresh_toaster/hooks/index.tsx";

import FileItem from "@/islands/FileItem.tsx";
import Converter, { Format } from "@/islands/Converter.tsx";

type FormProps = {
  uploadUrl: string;
};

export default function Form(props: FormProps) {
  const [toasts, toaster] = useToaster();
  const [files, setFiles] = useState<FileList | null>(null);
  const [formats, setFormats] = useState<Format[]>([]);

  const handleFileChange: JSX.GenericEventHandler<HTMLInputElement> = (
    event,
  ) => {
    const target = event.target as HTMLInputElement;
    const selectedFile = target?.files;
    if (selectedFile) {
      setFiles(selectedFile);
    }
  };

  const handleDragOver: JSX.DragEventHandler<HTMLLabelElement> = (event) => {
    event.preventDefault();
    event.currentTarget.classList.add("bg-gray-300");
  };

  const handleDragLeave: JSX.DragEventHandler<HTMLLabelElement> = (event) => {
    event.preventDefault();
    event.currentTarget.classList.remove("bg-gray-300");
  };

  const handleDrop: JSX.DragEventHandler<HTMLLabelElement> = (event) => {
    event.preventDefault();
    event.currentTarget.classList.remove("bg-gray-300");

    const droppedFile = event.dataTransfer?.files;
    if (droppedFile) {
      setFiles(droppedFile);
    }
  };

  return (
    <div class="w-full">
      <div className="flex items-center justify-center w-full">
        <label
          htmlFor="dropzone-file"
          className="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-300 border-dashed rounded-lg cursor-pointer bg-gray-200 hover:bg-gray-300"
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <div className="flex flex-col items-center justify-center pt-5 pb-6">
            <svg
              className="w-8 h-8 mb-4 text-gray-500"
              aria-hidden="true"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 20 16"
            >
            </svg>
            <p className="mb-2 text-sm text-gray-500">
              <span className="font-semibold">Click to upload</span>{" "}
              or drag and drop
            </p>
            <p className="text-xs text-gray-500">
              PNG, JPG, or WEBP
            </p>
          </div>
          <input
            id="dropzone-file"
            type="file"
            className="hidden"
            onChange={handleFileChange}
            multiple
          />
        </label>
      </div>
      <Converter onFormatChange={(nf) => setFormats(nf)} />
      {files && (
        [...files].map((file: File) => (
          <FileItem file={file} uploadUrl={props.uploadUrl} formats={formats} />
        ))
      )}
    </div>
  );
}
