export const isArrayWithElements = (data: any): boolean => {
  if (Array.isArray(data)) {
    return data.length > 0;
  }
  return false;
};
