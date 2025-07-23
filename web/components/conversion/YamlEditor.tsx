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
      style={{ height, minHeight: height }}
      className={`relative rounded-md overflow-hidden border border-gray-300 dark:border-gray-600 ${onDoubleClick ? 'cursor-pointer' : ''}`}
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
          padding: '0.75rem',
          resize: 'none',
          fontSize: '14px',
          fontFamily: '"JetBrains Mono", "Fira Code", "SF Mono", Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
          lineHeight: '1.4',
          overflow: 'auto',
          background: 'transparent',
        }}
      >
        {displayValue}
      </SyntaxHighlighter>

      {!readOnly && onChange && (
        <textarea
          className="absolute top-0 left-0 w-full h-full opacity-0 resize-none text-sm caret-black dark:caret-white focus:opacity-100 focus:bg-white dark:focus:bg-gray-800 focus:text-black dark:focus:text-white p-3 font-mono leading-relaxed"
          style={{
            fontFamily: '"JetBrains Mono", "Fira Code", "SF Mono", Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
            fontSize: '14px',
            lineHeight: '1.4',
          }}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder={placeholder}
        />
      )}
    </div>
  );
};

export default YamlEditor;