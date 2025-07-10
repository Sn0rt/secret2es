import { useState, useCallback } from 'react';
import { ConversionRequest, ConversionResponse } from '@/types/conversion';
import { API_ENDPOINTS, HTTP_METHODS, CONTENT_TYPES } from '@/constants/api';
import { parseErrorResponse, handleFetchError, ApiError } from '@/utils/error-handling';

export const useConversion = () => {
  const [isLoading, setIsLoading] = useState(false);

  const convertYaml = useCallback(async (request: ConversionRequest): Promise<ConversionResponse> => {
    setIsLoading(true);
    
    try {
      const response = await fetch(API_ENDPOINTS.convert, {
        method: HTTP_METHODS.POST,
        headers: {
          'Content-Type': CONTENT_TYPES.JSON,
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        const apiError = await parseErrorResponse(response);
        console.error('API Error:', {
          message: apiError.message,
          status: apiError.status,
          details: apiError.details,
          code: apiError.code,
        });
        throw apiError;
      }

      const data = await response.json();
      return data;
    } catch (error) {
      const apiError = handleFetchError(error);
      console.error('Conversion Error:', {
        message: apiError.message,
        status: apiError.status,
        details: apiError.details,
        code: apiError.code,
        originalError: error,
      });
      throw apiError;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    isLoading,
    convertYaml,
  };
};