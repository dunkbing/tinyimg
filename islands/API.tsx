import Code from "@/islands/Code.tsx";

export default function API(props: { baseUrl: string }) {
  const sampleResponse = {
    data: [
      {
        savedBytes: 0,
        newSize: 0,
        imageUrl: "",
        format: "",
      },
    ],
    errors: [],
  };

  return (
    <div class="flex flex-col items-center justify-center px-5 bg-white rounded shadow-lg p-8">
      <div>
        <h1 class="text-2xl font-semibold mb-4">TinyIMG API Instructions</h1>

        <p class="mb-4">
          Welcome to the TinyIMG API! Follow the instructions below to compress
          your images using our API.
        </p>

        <h2 class="text-xl font-semibold mb-2">1. Make a Request</h2>
        <p class="mb-2">
          Use the following <b>curl</b>{" "}
          command to make a sample request to convert text to audio:
        </p>

        <Code
          code={`curl --location '${props.baseUrl}/upload' --form 'file=@"/path/to/image"' --form 'formats="png,jpg,webp"'`}
        />

        <h2 className="text-xl font-semibold mb-2">2. Response</h2>
        <p className="mb-2">
          A successful response will be an array of objects like this:
        </p>

        <Code code={JSON.stringify(sampleResponse, null, 2)} />

        <p className="mt-4">
          Each object in the array contains a URL and the corresponding text.
        </p>
      </div>
    </div>
  );
}
