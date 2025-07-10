import React from 'react';
import YamlEditor from './YamlEditor';
import { copyToClipboard } from '@/utils/clipboard';
import { PLACEHOLDERS } from '@/constants/ui';

interface ConversionResultProps {
  outputYaml: string;
  contentHeight: string;
}

const ConversionResult: React.FC<ConversionResultProps> = ({
  outputYaml,
  contentHeight,
}) => {
  const handleDoubleClick = () => {
    if (outputYaml) {
      copyToClipboard(outputYaml);
    }
  };

  return (
    <YamlEditor
      value={outputYaml}
      placeholder={PLACEHOLDERS.output}
      height={contentHeight}
      readOnly
      onDoubleClick={handleDoubleClick}
    />
  );
};

export default ConversionResult;