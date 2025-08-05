/**
 * Removes all search parameters from an image URL
 * @param imageUrl - The image URL that may contain search parameters
 * @returns The clean image URL without search parameters
 */
export function removeSearchParams(imageUrl: string): string {
  try {
    const url = new URL(imageUrl);
    url.search = "";
    return url.toString();
  } catch (error) {
    const questionMarkIndex = imageUrl.indexOf("?");
    if (questionMarkIndex !== -1) {
      return imageUrl.substring(0, questionMarkIndex);
    }
    return imageUrl;
  }
}
