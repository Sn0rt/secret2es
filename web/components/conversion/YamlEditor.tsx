import React from 'react';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism';
import { formatYamlForDisplay } from '@/utils/yaml';
import { DEFAULT_HEIGHT } from '@/constants/ui';

interface YamlEditorProps {
  value: string;
  onChange?: (value: string) => void;
  placeholder: string;
  height?: string;
  readOnly?: boolean;
  onDoubleClick?: () => void;
}

const YamlEditor: React.FC<YamlEditorProps> = ({
  value,
  onChange,
  placeholder,
  height = DEFAULT_HEIGHT,
  readOnly = false,
  onDoubleClick,
}) => {
  const displayValue = formatYamlForDisplay(value, placeholder);

  return (
    <div
      style={{ height, minHeight: DEFAULT_HEIGHT }}
      className={`relative rounded-md overflow-hidden ${onDoubleClick ? 'cursor-pointer' : ''}`}
      onDoubleClick={onDoubleClick}
    >
      <SyntaxHighlighter
        language="yaml"
        style={tomorrow}
        customStyle={{
          position: readOnly ? 'relative' : 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          margin: 0,
          padding: '1rem',
          resize: 'none',
          fontSize: '16px',
          overflow: 'auto',
        }}
      >
        {displayValue}
      </SyntaxHighlighter>
      
      {!readOnly && onChange && (
        <textarea
          className="absolute top-0 left-0 w-full h-full opacity-0 resize-none text-base caret-black dark:caret-white focus:opacity-100 focus:bg-white dark:focus:bg-gray-800 focus:text-black dark:focus:text-white p-4"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
        />
      )}
    </div>
  );
};

export default YamlEditor;