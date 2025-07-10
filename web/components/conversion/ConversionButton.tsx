import React from 'react';
import { Button } from '@/components/ui/button';
import { Loading } from '@/components/ui/loading';

interface ConversionButtonProps {
  onClick: () => void;
  isLoading: boolean;
  disabled: boolean;
}

const ConversionButton: React.FC<ConversionButtonProps> = ({
  onClick,
  isLoading,
  disabled,
}) => {
  return (
    <Button
      onClick={onClick}
      className="w-full px-4 py-2 text-base bg-blue-500 hover:bg-blue-600 text-white my-4"
      disabled={disabled || isLoading}
    >
      {isLoading ? (
        <>
          <Loading size="sm" className="mr-2" />
          Converting...
        </>
      ) : (
        'Convert'
      )}
    </Button>
  );
};

export default ConversionButton;