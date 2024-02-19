export function FAQ() {
  return (
    <div
      id="faq"
      class="flex flex-col items-center justify-center mt-2 px-5"
    >
      <hr class="w-48 h-1 mx-auto my-4 border-0 rounded md:my-2 bg-gray-700" />
      <h1 class="text-black text-2xl font-semibold mb-4">FAQ</h1>
      <div class="flex flex-col sm:px-8 lg:px-12 xl:px-32">
        <details>
          <summary class="py-2 outline-none cursor-pointer focus:underline text-center">
            Why should I compress images?
          </summary>
          <div class="px-4 pb-4">
            <p class="text-center">
              Image files from sources like professional DSLR cameras or
              smartphones can be quite large, occupying significant storage
              space. Compressing images helps in reducing file sizes, making
              them more manageable and freeing up storage.
            </p>
          </div>
        </details>
        <details>
          <summary class="py-2 outline-none cursor-pointer focus:underline text-center">
            How does the image compressor work?
          </summary>
          <div class="px-4 pb-4">
            <p class="text-center">
              This tool utilizes lossy compression for PNG, JPG/JPEG, and Webp
              files. You can upload up to 20 images of different types
              simultaneously. The server intelligently analyzes and reduces
              images to the smallest file size without compromising quality.
              Users can adjust compression rates using a quality slider and
              download a ZIP file with all compressed images.
            </p>
          </div>
        </details>
        <details>
          <summary class="py-2 outline-none cursor-pointer focus:underline text-center">
            Is it safe to compress images?
          </summary>
          <div class="px-4 pb-4 space-y-2">
            <p class="text-center">
              Absolutely. Your original files remain untouched, allowing you to
              retry if needed. Additionally, our system automatically purges all
              data after one hour, ensuring the security of your information.
              There's no cost associated with using our service, and you can use
              the tool as many times as necessary.
            </p>
          </div>
        </details>
      </div>
    </div>
  );
}
