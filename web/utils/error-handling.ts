import { ErrorResponse } from '@/types/conversion';

export class ApiError extends Error {
  public readonly status: number;
  public readonly details?: string;
  public readonly code?: string;

  constructor(message: string, status: number = 500, details?: string, code?: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.details = details;
    this.code = code;
  }
}

export const parseErrorResponse = async (response: Response): Promise<ApiError> => {
  const contentType = response.headers.get('content-type');
  
  try {
    if (contentType && contentType.includes('application/json')) {
      // Try to parse as JSON
      const errorData = await response.json();
      
      if (errorData.error) {
        return new ApiError(
          errorData.error,
          response.status,
          errorData.details,
          errorData.code
        );
      }
      
      // Fallback for unexpected JSON structure
      return new ApiError(
        JSON.stringify(errorData),
        response.status
      );
    } else {
      // Handle plain text error responses
      const textError = await response.text();
      return new ApiError(
        textError || `HTTP Error ${response.status}`,
        response.status
      );
    }
  } catch (parseError) {
    // If parsing fails, return a generic error
    return new ApiError(
      `Failed to parse error response: ${response.statusText}`,
      response.status
    );
  }
};

export const handleFetchError = (error: unknown): ApiError => {
  if (error instanceof ApiError) {
    return error;
  }
  
  if (error instanceof TypeError && error.message.includes('fetch')) {
    // Network error
    return new ApiError(
      'Network error: Please check your connection and try again.',
      0,
      error.message,
      'NETWORK_ERROR'
    );
  }
  
  if (error instanceof Error) {
    return new ApiError(
      error.message,
      500,
      undefined,
      'UNKNOWN_ERROR'
    );
  }
  
  return new ApiError(
    'An unexpected error occurred',
    500,
    String(error),
    'UNKNOWN_ERROR'
  );
};

export const getErrorMessage = (error: unknown): string => {
  if (error instanceof ApiError) {
    return error.message;
  }
  
  if (error instanceof Error) {
    return error.message;
  }
  
  return 'An unexpected error occurred';
};

export const shouldRetry = (error: ApiError): boolean => {
  // Retry on network errors or 5xx server errors
  return error.code === 'NETWORK_ERROR' || (error.status >= 500 && error.status < 600);
};

export const getUserFriendlyMessage = (error: ApiError): string => {
  switch (error.code) {
    case 'NETWORK_ERROR':
      return 'Unable to connect to the server. Please check your internet connection and try again.';
    case 'VALIDATION_ERROR':
      return 'Please check your input and try again.';
    default:
      // For conversion errors, always show the actual error message
      if (error.message && error.message.includes('Conversion error:')) {
        return error.message;
      }
      
      if (error.status >= 500) {
        return 'A server error occurred. Please try again later.';
      }
      if (error.status >= 400 && error.status < 500) {
        return error.message; // Client errors usually have good messages
      }
      return 'An unexpected error occurred. Please try again.';
  }
};