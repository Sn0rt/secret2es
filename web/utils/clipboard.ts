import { toast } from 'react-hot-toast';

export const copyToClipboard = async (text: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard!');
  } catch (error) {
    console.error('Failed to copy text: ', error);
    toast.error('Failed to copy to clipboard');
  }
};