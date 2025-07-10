import { useState, useCallback } from 'react';
import { EnvVar } from '@/types/conversion';

export const useEnvironmentVariables = (initialEnvVars: EnvVar[] = [{ key: '', value: '' }]) => {
  const [envVars, setEnvVars] = useState<EnvVar[]>(initialEnvVars);

  const handleAddEnvVar = useCallback(() => {
    setEnvVars(prev => [...prev, { key: '', value: '' }]);
  }, []);

  const handleEnvVarChange = useCallback((index: number, field: 'key' | 'value', value: string) => {
    setEnvVars(prev => {
      const newEnvVars = [...prev];
      newEnvVars[index][field] = value;
      return newEnvVars;
    });
  }, []);

  const handleRemoveEnvVar = useCallback((index: number) => {
    setEnvVars(prev => prev.filter((_, i) => i !== index));
  }, []);

  const getEnvVarsAsRecord = useCallback((): Record<string, string> => {
    return envVars.reduce((acc, { key, value }) => {
      if (key.trim()) {
        acc[key] = value;
      }
      return acc;
    }, {} as Record<string, string>);
  }, [envVars]);

  return {
    envVars,
    handleAddEnvVar,
    handleEnvVarChange,
    handleRemoveEnvVar,
    getEnvVarsAsRecord,
  };
};