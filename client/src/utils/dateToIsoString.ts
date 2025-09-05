export const dateToIsoString = (dateString: string): string => {
  const date = new Date(dateString);

  if (isNaN(date.getTime())) {
    throw new Error("Invalid date string provided");
  }

  return date.toISOString();
};
