import { signal } from "@preact/signals";
import { JSX } from "preact";
import { useMemo, useRef, useState } from "preact/hooks";
import { useToaster } from "fresh_toaster/hooks/index.tsx";

import { Button } from "@/components/Button.tsx";
import { Loader } from "@/components/Loader.tsx";
import { downloadFile } from "@/utils/http.ts";

type Audio = { url: string; text: string; index: number };

const AudioCard = (props: Audio) => {
  const audio = useMemo(() => new Audio(props.url), [props.url]);
  const [playing, setPlaying] = useState(false);

  audio.onpause = () => setPlaying(false);

  return (
    <div class="w-2/3 mx-auto bg-white rounded-lg overflow-hidden shadow-lg mb-2">
      <div class="p-4">
        <p class="text-gray-700 text-center">
          {props.index}. {props.text}
        </p>
      </div>
      <div class="px-4 pb-4 flex justify-center">
        <Button
          class="text-white font-semibold"
          onClick={() => {
            if (playing) {
              setPlaying(false);
              audio.pause();
            } else {
              setPlaying(true);
              audio.play();
            }
          }}
        >
          {playing ? "Pause" : "Play"}
        </Button>
        <div class="w-5" />
        <Button
          class="text-white font-semibold"
          colorMode="secondary"
          onClick={() => downloadFile(props.url, String(props.index))}
        >
          Download
        </Button>
      </div>
    </div>
  );
};

const audios = signal<Audio[]>([]);
const converting = signal(false);

export default function Form() {
  const [toasts, toaster] = useToaster();
  const [files, setFiles] = useState<FileList | null>(null);

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
    <>
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
              {/* SVG path data */}
            </svg>
            <p className="mb-2 text-sm text-gray-500">
              <span className="font-semibold">Click to upload</span>{" "}
              or drag and drop
            </p>
            <p className="text-xs text-gray-500">
              SVG, PNG, JPG, or GIF (MAX. 800x400px)
            </p>
          </div>
          <input
            id="dropzone-file"
            type="file"
            className="hidden"
            onChange={handleFileChange}
          />
        </label>

        {/* Display the selected file name */}
        {files && <p>{[...files].map((f) => f.name)}</p>}
      </div>

      {converting.value ? <Loader /> : (
        <>
          {audios.value.length > 1
            ? (
              <Button
                class="text-white font-semibold mb-2"
                colorMode="secondary"
                onClick={() =>
                  audios.value.map((v, i) =>
                    downloadFile(v.url, String(i + 1)).catch((err) =>
                      toaster.error(err)
                    )
                  )}
              >
                Download All
              </Button>
            )
            : null}
          {audios.value.map((v, i) => (
            <AudioCard key={i} text={v.text} url={v.url} index={i + 1} />
          ))}
        </>
      )}
    </>
  );
}
