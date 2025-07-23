export const API_ENDPOINTS = {
  convert: '/api/convert',
};

export const HTTP_METHODS = {
  POST: 'POST',
  GET: 'GET',
} as const;

export const CONTENT_TYPES = {
  JSON: 'application/json',
} as const;

export const DEFAULT_VALUES = {
  storeType: 'SecretStore',
  creationPolicy: 'Owner',
  resolve: false,
  refreshPolicy: 'OnChange',
  refreshInterval: '1h',
} as const;