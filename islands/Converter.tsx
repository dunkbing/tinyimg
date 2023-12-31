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
    // setSelectedFormats([]);
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
    <div className="bg-white p-6 rounded-lg shadow-md">
      <label className="flex items-center space-x-2 cursor-pointer">
        <input
          type="checkbox"
          checked={convertImage}
          onChange={toggleConvertImage}
          className="hidden"
        />
        <div className="relative flex">
          <div
            className={`w-12 h-6 bg-gray-300 rounded-full transition-colors duration-300 ease-in-out flex items-center ${
              convertImage ? "bg-green-500" : "bg-gray-300"
            }`}
          >
            <div
              className={`bg-white w-5 h-5 rounded-full shadow-md transform transition-transform duration-300 ease-in-out ${
                convertImage ? "translate-x-6" : "translate-x-1"
              }`}
            />
          </div>
        </div>
        <span className="text-lg font-semibold text-gray-800">
          Convert my image
        </span>
      </label>

      {convertImage && (
        <div className="mt-4">
          <p className="text-lg font-semibold mb-2 text-gray-800">
            Select image formats:
          </p>
          <div className="flex space-x-4 divide-x">
            <button
              className={`px-4 py-2 rounded-full ${
                selectedFormats.includes("png")
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("png")}
            >
              PNG
            </button>
            <button
              className={`px-4 py-2 rounded-full ${
                selectedFormats.includes("jpg")
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("jpg")}
            >
              JPG
            </button>
            <button
              className={`px-4 py-2 rounded-full ${
                selectedFormats.includes("webp")
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("webp")}
            >
              WebP
            </button>
            <button
              className={`px-4 py-2 rounded-full ${
                selectedFormats.length === 3
                  ? "bg-green-500 text-white"
                  : "bg-gray-200 text-gray-700"
              } hover:bg-green-600 hover:text-white focus:outline-none focus:ring focus:border-blue-300`}
              onClick={() => handleFormatToggle("selectAll")}
            >
              Select All
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default Converter;