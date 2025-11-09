"use client";

import Image, { StaticImageData } from "next/image";

interface ThemedImageProps {
  lightSrc: StaticImageData;
  darkSrc: StaticImageData;
  alt: string;
  className?: string;
}

export function ThemedImage({
  lightSrc,
  darkSrc,
  alt,
  className = "",
}: ThemedImageProps) {
  return (
    <>
      <Image
        src={lightSrc}
        alt={alt}
        className={`dark:hidden ${className}`}
        style={{ width: "100%", height: "auto" }}
        loading="eager"
      />
      <Image
        src={darkSrc}
        alt={alt}
        className={`hidden dark:block ${className}`}
        style={{ width: "100%", height: "auto" }}
        loading="eager"
      />
    </>
  );
}
