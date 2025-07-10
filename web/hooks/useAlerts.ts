import { useState, useCallback } from 'react';

export const useAlerts = () => {
  const [warning, setWarning] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const showWarning = useCallback((message: string) => {
    setWarning(message);
  }, []);

  const showError = useCallback((message: string) => {
    setErrorMessage(message);
  }, []);

  const clearWarning = useCallback(() => {
    setWarning(null);
  }, []);

  const clearError = useCallback(() => {
    setErrorMessage(null);
  }, []);

  const clearAll = useCallback(() => {
    setWarning(null);
    setErrorMessage(null);
  }, []);

  return {
    warning,
    errorMessage,
    showWarning,
    showError,
    clearWarning,
    clearError,
    clearAll,
  };
};