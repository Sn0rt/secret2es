import { MAX_LINES, LINE_HEIGHT, DEFAULT_HEIGHT } from '@/constants/ui';

export const validateYamlInput = (input: string): string | null => {
  const lineCount = input.split('\n').length;
  if (lineCount > MAX_LINES) {
    return `Input exceeds maximum limit of ${MAX_LINES} lines.`;
  }
  return null;
};

export const calculateContentHeight = (content: string): string => {
  const lineCount = content.split('\n').length;
  const newHeight = Math.max(parseInt(DEFAULT_HEIGHT), lineCount * LINE_HEIGHT);
  return `${newHeight}px`;
};

export const formatYamlForDisplay = (yaml: string, placeholder: string): string => {
  return yaml || placeholder;
};