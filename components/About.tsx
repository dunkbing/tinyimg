export function About() {
  return (
    <div
      id="about"
      class="flex flex-col items-center justify-center mt-5 px-5"
    >
      <h1 class="text-black text-2xl font-semibold mb-4">
        Optimising image losslessly
      </h1>
      <p>
        In the digital age, image optimization is a crucial aspect of web
        development. It not only enhances the user experience but also
        contributes to a faster loading time, which is a significant factor in
        SEO. In this blog post, we will explore various tools that can be used
        to losslessly optimize images, thereby reducing their file sizes without
        compromising their quality.
      </p>
      <h2 class="text-black text-2xl font-semibold mt-5 mb-2">
        Tools for Lossless Image Optimization
      </h2>
      <ol class="space-y-2">
        <li>
          <p>
            <strong>pngquant</strong>{" "}
            by Kornelski: This tool is designed to reduce the file size of PNG
            images without losing any quality. It uses a unique conversion
            algorithm that allows it to achieve significant file size
            reductions.
          </p>
          <pre><code class="language-bash">pngquant --quality=0-80 --speed=1 $IN.png --output $OUT.png --force --strip</code></pre>
        </li>
        <li>
          <p>
            <strong>jpegoptim</strong>{" "}
            by Timo Kokkonen: This tool is specifically designed for JPEG
            images. It uses a variety of optimization techniques to reduce the
            file size of JPEG images while maintaining their quality.
          </p>
          <pre><code class="language-bash">jpegoptim --all-normal --verbose $IN.jpg $OUT.jpg</code></pre>
        </li>
        <li>
          <p>
            <strong>cwebp</strong>{" "}
            by Google: This tool is used to convert images to the WebP format, a
            modern image format that provides superior lossy and lossless
            compression for images on the web.
          </p>
          <pre><code class="language-bash">cwebp -q 80 $IN.png -o $OUT.webp</code></pre>
        </li>
        <li>
          <p>
            <strong>gifsicle</strong>{" "}
            by Eddie Kohler: This tool is used to optimize GIF images. It uses a
            combination of lossless and lossy compression techniques to reduce
            the file size of GIF images.
          </p>
          <pre><code class="language-bash">gifsicle -O3 --verbose -i $IN.gif -o $OUT.gif</code></pre>
        </li>
        <li>
          <p>
            <strong>scour</strong>{" "}
            by Jeff Schiller and Louis Simard: This tool is specifically
            designed for SVG images. It uses a combination of optimization
            techniques to reduce the file size of SVG images while maintaining
            their quality.
          </p>
          <pre><code class="language-bash">scour -i $IN.svg -o $OUT.svg</code></pre>
        </li>
      </ol>
      <h2 class="text-black text-2xl font-semibold mt-5 mb-2">
        How This Website Optimizes Images
      </h2>
      <p>
        This website employs a strategy to invoke the right optimizer based on
        the extension of the input file. Every file that ends up here goes
        through one of these tools. This ensures that the right optimization
        technique is applied to each image type, thereby achieving optimal
        results.
      </p>
      <p>
        In conclusion, lossless image optimization is a powerful technique that
        can significantly enhance the performance of your website. By using the
        right tools and techniques, you can reduce the file size of your images
        without compromising their quality.
      </p>
    </div>
  );
}
