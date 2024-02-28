import { JSX } from "preact";
import { useState } from "preact/hooks";
import IconPhoto from "tabler_icons_tsx/photo.tsx"

import FileItem from "@/islands/FileItem.tsx";
import Converter, { Format } from "@/islands/Converter.tsx";
import { signal } from "@preact/signals";
import { Button } from "@/components/Button.tsx";
import { Loader } from "@/components/Loader.tsx";

type FormProps = {
  uploadUrl: string;
  downloadUrl: string;
};

const filesSig = signal<string[]>([]);
const downloadSig = signal<boolean>(false);

export default function Form(props: FormProps) {
  const [files, setFiles] = useState<FileList | null>(null);
  const [formats, setFormats] = useState<Format[]>([]);

  const handleFileChange: JSX.GenericEventHandler<HTMLInputElement> = (
    event,
  ) => {
    const target = event.target as HTMLInputElement;
    const selectedFile = target?.files;
    if (selectedFile && selectedFile.length > 20) {
      alert("You can only upload a maximum of 20 files");
      return;
    }
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

  const downloadAll = async () => {
    downloadSig.value = true;
    const myHeaders = new Headers();
    myHeaders.append("Content-Type", "application/json");
    const raw = JSON.stringify({
      "files": filesSig.value,
    });

    const requestOptions: RequestInit = {
      method: "POST",
      body: raw,
    };

    try {
      const response = await fetch(
        props.downloadUrl,
        requestOptions,
      );
      const blob = await response.blob();
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = "images.zip";
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      console.error(error);
    } finally {
      downloadSig.value = false;
    }
  };

  return (
    <div class="w-full flex flex-col items-center space-y-4 mt-8">
      <h1 class="text-green-700 text-base text-center w-3/4">
        Efficient WebP, PNG, and JPEG Compression for Faster Websites
      </h1>
      <label
        for="dropzone-file"
        class="block text-gray-700 font-bold mb-1 text-xl text-center"
      >
        Files
      </label>
      <div className="flex items-center justify-center w-full px-2">
        <label
          htmlFor="dropzone-file"
          className="flex flex-col items-center justify-center w-full h-64 border-2 border-gray-400 border-dashed rounded-lg cursor-pointer bg-gray-200 hover:bg-gray-300 shadow-lg"
          onDragOver={handleDragOver}
          onDragLeave={handleDragLeave}
          onDrop={handleDrop}
        >
          <div className="flex flex-col items-center justify-center space-y-1">
            <IconPhoto size={55} />
            <p className="text-sm text-gray-500">
              <span className="font-semibold">Click to upload</span>{" "}
              or drop your images here
            </p>
            <p className="mb-2 text-xs text-gray-500">
              Upto 20 images at a time
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
            accept="image/png, image/jpeg, image/webp"
          />
        </label>
      </div>
      <Converter
        onFormatChange={(nf) => {
          filesSig.value = [];
          setFormats(nf);
        }}
      />
      {(!!files?.length &&
        (filesSig.value.length === files.length * (formats.length || 1))) && (
          <Button colorMode="secondary" onClick={downloadAll}>
            {downloadSig.value && <Loader />}
            Download all
          </Button>
        )}
      {files && (
        [...files].map((file: File) => (
          <FileItem
            file={file}
            uploadUrl={props.uploadUrl}
            formats={formats}
            filesSig={filesSig}
          />
        ))
      )}
    </div>
  );
}
