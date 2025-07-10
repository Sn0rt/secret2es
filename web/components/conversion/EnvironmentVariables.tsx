import React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { X } from 'lucide-react';
import { PLACEHOLDERS } from '@/constants/ui';
import { useConversionStore, useEnvVars } from '@/stores/conversionStore';

interface EnvironmentVariablesProps {}

const EnvironmentVariables: React.FC<EnvironmentVariablesProps> = () => {
  // Use Zustand store directly
  const envVars = useEnvVars();
  const { addEnvVar, updateEnvVar, removeEnvVar } = useConversionStore();
  return (
    <div className="space-y-4">
      <h3 className="font-bold text-base">Environment Variables</h3>
      
      {envVars.map((envVar, index) => (
        <div key={index} className="space-y-2">
          <div className="flex items-center space-x-2">
            <Input
              type="text"
              className="flex-1 text-base"
              placeholder={PLACEHOLDERS.envKey}
              value={envVar.key}
              onChange={(e) => updateEnvVar(index, 'key', e.target.value)}
            />
            {envVars.length > 1 && (
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => removeEnvVar(index)}
                className="p-2"
              >
                <X className="h-4 w-4" />
              </Button>
            )}
          </div>
          
          <Input
            type="text"
            className="w-full text-base"
            placeholder={PLACEHOLDERS.envValue}
            value={envVar.value}
            onChange={(e) => updateEnvVar(index, 'value', e.target.value)}
          />
        </div>
      ))}
      
      <Button
        type="button"
        onClick={addEnvVar}
        className="w-full text-base"
        variant="outline"
      >
        Add Environment Variable
      </Button>
    </div>
  );
};

export default EnvironmentVariables;