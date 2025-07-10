"use client"

import React from "react";
import { ArrowRight } from "lucide-react";
import { PLACEHOLDERS } from '@/constants/ui';

// Components
import ConversionForm from '@/components/conversion/ConversionForm';
import EnvironmentVariables from '@/components/conversion/EnvironmentVariables';
import YamlEditor from '@/components/conversion/YamlEditor';
import ConversionResult from '@/components/conversion/ConversionResult';
import AlertMessages from '@/components/conversion/AlertMessages';
import ConversionButton from '@/components/conversion/ConversionButton';

// Zustand Stores
import { 
  useConversionStore,
  useFormData,
  useIsFormValid,
  useInputYaml,
  useOutputYaml,
  useContentHeight,
  useYamlError,
  useIsLoading,
  useWarning,
  useErrorMessage,
} from '@/stores/conversionStore';
// import { useSettingsStore, useLastSavedState, useAutoSave } from '@/stores/settingsStore';

export default function Home() {
  // Zustand store selectors (optimized for performance)
  const formData = useFormData();
  const isFormValid = useIsFormValid();
  const inputYaml = useInputYaml();
  const outputYaml = useOutputYaml();
  const contentHeight = useContentHeight();
  const yamlError = useYamlError();
  const isLoading = useIsLoading();
  const warning = useWarning();
  const errorMessage = useErrorMessage();
  
  // Zustand store actions
  const {
    setInputYaml,
    clearWarning,
    clearError,
    convertYaml,
  } = useConversionStore();
  
  // Temporarily disable auto-save to isolate the infinite loop issue
  // const autoSave = useAutoSave();
  // const { lastFormState, lastInputYaml } = useLastSavedState();
  // const { saveFormState } = useSettingsStore();

  // Conversion handler (now simplified)
  const handleConvert = async () => {
    await convertYaml();
  };

  return (
    <div className="container mx-auto p-4 flex-grow flex">
      {/* Left sidebar - Form controls */}
      <div className="w-1/6 pr-4 flex flex-col space-y-4">
        <ConversionForm />
        
        <ConversionButton
          onClick={handleConvert}
          isLoading={isLoading}
          disabled={!isFormValid}
        />
        
        {formData.resolve && (
          <EnvironmentVariables />
        )}
      </div>

      {/* Main content area - YAML editors */}
      <div className="w-5/6 flex items-start relative">
        {/* Input YAML editor */}
        <div className="w-[48%]">
          <YamlEditor
            value={inputYaml}
            onChange={setInputYaml}
            placeholder={PLACEHOLDERS.input}
            height={contentHeight}
          />
          {yamlError && (
            <p className="text-red-500 mt-2 text-base">{yamlError}</p>
          )}
        </div>

        {/* Arrow separator */}
        <div className="w-[4%] flex justify-center items-start pt-2">
          <ArrowRight className="h-10 w-10 text-blue-500" />
        </div>

        {/* Output YAML editor */}
        <div className="w-[48%]">
          <ConversionResult
            outputYaml={outputYaml}
            contentHeight={contentHeight}
          />
        </div>

        {/* Alert messages */}
        <AlertMessages
          warning={warning}
          error={errorMessage}
          onClearWarning={clearWarning}
          onClearError={clearError}
        />
      </div>
    </div>
  );
}