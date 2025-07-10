import React from 'react';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';

interface AlertMessagesProps {
  warning?: string | null;
  error?: string | null;
  onClearWarning?: () => void;
  onClearError?: () => void;
}

const AlertMessages: React.FC<AlertMessagesProps> = ({
  warning,
  error,
  onClearWarning,
  onClearError,
}) => {
  if (!warning && !error) {
    return null;
  }

  return (
    <div className="absolute top-0 right-0 max-w-md space-y-2">
      {warning && (
        <Alert variant="warning" onClose={onClearWarning}>
          <AlertTitle>Warning</AlertTitle>
          <AlertDescription>{warning}</AlertDescription>
        </Alert>
      )}
      
      {error && (
        <Alert variant="destructive" onClose={onClearError}>
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
    </div>
  );
};

export default AlertMessages;