import { 
  ApiError, 
  parseErrorResponse, 
  handleFetchError, 
  getErrorMessage, 
  getUserFriendlyMessage 
} from '../error-handling';

// Mock Response for testing
class MockResponse {
  private _status: number;
  private _data: any;
  private _contentType: string;

  constructor(data: any, status: number = 200, contentType: string = 'application/json') {
    this._status = status;
    this._data = data;
    this._contentType = contentType;
  }

  get status() {
    return this._status;
  }

  get ok() {
    return this._status >= 200 && this._status < 300;
  }

  get statusText() {
    return this._status === 404 ? 'Not Found' : 'Internal Server Error';
  }

  headers = {
    get: (name: string) => {
      if (name === 'content-type') {
        return this._contentType;
      }
      return null;
    }
  };

  async json() {
    if (this._contentType.includes('application/json')) {
      return this._data;
    }
    throw new Error('Not JSON');
  }

  async text() {
    return typeof this._data === 'string' ? this._data : JSON.stringify(this._data);
  }
}

describe('Error Handling Utilities', () => {
  describe('ApiError', () => {
    it('should create an ApiError with all properties', () => {
      const error = new ApiError('Test error', 400, 'Test details', 'TEST_CODE');
      
      expect(error.message).toBe('Test error');
      expect(error.status).toBe(400);
      expect(error.details).toBe('Test details');
      expect(error.code).toBe('TEST_CODE');
      expect(error.name).toBe('ApiError');
    });
  });

  describe('parseErrorResponse', () => {
    it('should parse JSON error response correctly', async () => {
      const mockResponse = new MockResponse(
        { error: 'Conversion failed', details: 'Invalid YAML' },
        400
      ) as any;

      const error = await parseErrorResponse(mockResponse);
      
      expect(error.message).toBe('Conversion failed');
      expect(error.status).toBe(400);
      expect(error.details).toBe('Invalid YAML');
    });

    it('should handle plain text error response', async () => {
      const mockResponse = new MockResponse(
        'Plain text error',
        500,
        'text/plain'
      ) as any;

      const error = await parseErrorResponse(mockResponse);
      
      expect(error.message).toBe('Plain text error');
      expect(error.status).toBe(500);
    });

    it('should handle malformed JSON response', async () => {
      const mockResponse = new MockResponse(
        'Not JSON',
        500,
        'application/json'
      ) as any;

      // Override json method to throw
      mockResponse.json = () => Promise.reject(new Error('Invalid JSON'));

      const error = await parseErrorResponse(mockResponse);
      
      expect(error.message).toBe('Failed to parse error response: Internal Server Error');
      expect(error.status).toBe(500);
    });
  });

  describe('handleFetchError', () => {
    it('should return ApiError as-is', () => {
      const originalError = new ApiError('Test', 400);
      const result = handleFetchError(originalError);
      
      expect(result).toBe(originalError);
    });

    it('should handle network errors', () => {
      const networkError = new TypeError('fetch failed');
      const result = handleFetchError(networkError);
      
      expect(result.code).toBe('NETWORK_ERROR');
      expect(result.message).toBe('Network error: Please check your connection and try again.');
    });

    it('should handle generic errors', () => {
      const genericError = new Error('Something went wrong');
      const result = handleFetchError(genericError);
      
      expect(result.message).toBe('Something went wrong');
      expect(result.code).toBe('UNKNOWN_ERROR');
    });
  });

  describe('getUserFriendlyMessage', () => {
    it('should return user-friendly message for network errors', () => {
      const error = new ApiError('Network failed', 0, undefined, 'NETWORK_ERROR');
      const message = getUserFriendlyMessage(error);
      
      expect(message).toBe('Unable to connect to the server. Please check your internet connection and try again.');
    });

    it('should return user-friendly message for server errors', () => {
      const error = new ApiError('Internal error', 500);
      const message = getUserFriendlyMessage(error);
      
      expect(message).toBe('A server error occurred. Please try again later.');
    });

    it('should return original message for client errors', () => {
      const error = new ApiError('Validation failed', 400);
      const message = getUserFriendlyMessage(error);
      
      expect(message).toBe('Validation failed');
    });
  });
});

// Example of how errors will now be displayed to users
describe('Error Display Examples', () => {
  it('should demonstrate improved error handling', () => {
    // Backend conversion error (now properly formatted as JSON)
    const backendError = {
      error: "Conversion error: secret 'test-secret' contains invalid AVP path"
    };
    
    // This will now be properly caught and displayed in the UI
    expect(backendError.error).toContain('Conversion error');
    
    // Network error
    const networkError = new ApiError(
      'Network error', 
      0, 
      'Failed to fetch', 
      'NETWORK_ERROR'
    );
    
    const userMessage = getUserFriendlyMessage(networkError);
    expect(userMessage).toBe('Unable to connect to the server. Please check your internet connection and try again.');
  });
});