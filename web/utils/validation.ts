export const validateStoreName = (storeName: string): boolean => {
  return storeName.trim() !== '';
};

export const validateFormData = (storeName: string, error: string | null): boolean => {
  return validateStoreName(storeName) && !error;
};

export const validateEnvVars = (envVars: Array<{ key: string; value: string }>): Record<string, string> => {
  return envVars.reduce((acc, { key, value }) => {
    if (key.trim()) {
      acc[key] = value;
    }
    return acc;
  }, {} as Record<string, string>);
};