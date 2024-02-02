import { useState } from "preact/hooks";

export type Format = "png" | "jpg" | "webp" | "selectAll";

type ConverterProps = {
  onFormatChange?: (formats: Format[]) => void;
};

const Converter = ({ onFormatChange }: ConverterProps) => {
  const [convertImage, setConvertImage] = useState(false);
  const [selectedFormats, setSelectedFormats] = useState<Format[]>([]);

  const toggleConvertImage = () => {
    setConvertImage(!convertImage);
  };

  const handleFormatToggle = (format: Format) => {
    let newFormats: Format[] = [];
    if (format === "selectAll") {
      newFormats = ["png", "jpg", "webp"];
      setSelectedFormats(newFormats);
    } else {
      const updatedFormats = [...selectedFormats];
      const formatIndex = updatedFormats.indexOf(format);

      if (formatIndex !== -1) {
        updatedFormats.splice(formatIndex, 1);
      } else {
        updatedFormats.push(format);
      }

      newFormats = updatedFormats;
      setSelectedFormats(updatedFormats);
    }
    onFormatChange?.(newFormats);
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md w-80">
      <label className="flex items-center space-x-2 cursor-pointer">
        <label class="switch">
          <input
            type="checkbox"
            checked={convertImage}
            onChange={toggleConvertImage}
          />
          <span class="slider round"></span>
        </label>
        <span className="text-lg font-semibold text-gray-800">
          Convert images
        </span>
      </label>

      {convertImage && (
        <div className="mt-4">
          <p className="text-md font-semibold mb-2 text-gray-800">
            Select image formats:
          </p>
          <div className="flex space-x-2">
            <button
              className={`px-4 py-2 rounded-full text-sm font-semibold ${
                selectedFormats.includes("png")
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("png")}
            >
              PNG
            </button>
            <button
              className={`px-4 py-2 rounded-full text-sm font-semibold ${
                selectedFormats.includes("jpg")
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("jpg")}
            >
              JPG
            </button>
            <button
              className={`px-4 py-2 rounded-full text-sm font-semibold ${
                selectedFormats.includes("webp")
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("webp")}
            >
              WebP
            </button>
            <button
              className={`px-4 py-2 rounded-full text-sm font-semibold ${
                selectedFormats.length === 3
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("selectAll")}
            >
              All
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Converter;
