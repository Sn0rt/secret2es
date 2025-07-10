import { useState, useCallback } from 'react';
import { validateYamlInput, calculateContentHeight } from '@/utils/yaml';
import { DEFAULT_HEIGHT } from '@/constants/ui';

export const useYamlEditor = () => {
  const [inputYaml, setInputYaml] = useState('');
  const [outputYaml, setOutputYaml] = useState('');
  const [contentHeight, setContentHeight] = useState(DEFAULT_HEIGHT);
  const [error, setError] = useState<string | null>(null);

  const handleInputChange = useCallback((value: string) => {
    const validationError = validateYamlInput(value);
    
    if (validationError) {
      setError(validationError);
    } else {
      setError(null);
      setInputYaml(value);
      const newHeight = calculateContentHeight(value);
      setContentHeight(newHeight);
    }
  }, []);

  const setOutput = useCallback((output: string) => {
    setOutputYaml(output);
    const newHeight = calculateContentHeight(output);
    setContentHeight(newHeight);
  }, []);

  const clearOutput = useCallback(() => {
    setOutputYaml('');
  }, []);

  return {
    inputYaml,
    outputYaml,
    contentHeight,
    error,
    handleInputChange,
    setOutput,
    clearOutput,
  };
};