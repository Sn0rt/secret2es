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
    <div className="flex-grow flex h-full" style={{ minHeight: '100vh' }}>
      {/* Left sidebar - Form controls - Force it to always show */}
      <div
        className="flex-shrink-0 p-4 border-r border-gray-200 dark:border-gray-700 flex flex-col space-y-4 overflow-y-auto"
        style={{
          width: '320px',
          minWidth: '320px',
          maxWidth: '320px',
          position: 'relative',
          zIndex: 10
        }}
      >
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
      <div className="flex-1 flex items-stretch relative p-2" style={{ minWidth: 0 }}>
        {/* Input YAML editor */}
        <div className="flex-1 pr-1" style={{ minWidth: 0 }}>
          <YamlEditor
            value={inputYaml}
            onChange={setInputYaml}
            placeholder={PLACEHOLDERS.input}
            height="calc(100vh - 120px)"
          />
          {yamlError && (
            <p className="text-red-500 mt-1 text-sm">{yamlError}</p>
          )}
        </div>

        {/* Arrow separator */}
        <div className="w-8 flex justify-center items-center flex-shrink-0">
          <ArrowRight className="h-6 w-6 text-blue-500" />
        </div>

        {/* Output YAML editor */}
        <div className="flex-1 pl-1" style={{ minWidth: 0 }}>
          <ConversionResult
            outputYaml={outputYaml}
            contentHeight="calc(100vh - 120px)"
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