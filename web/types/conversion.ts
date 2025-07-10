export interface EnvVar {
  key: string;
  value: string;
}

export interface ConversionRequest {
  content: string;
  storeType: string;
  storeName: string;
  creationPolicy: string;
  resolve: boolean;
  envVars?: Record<string, string>;
}

export interface ConversionResponse {
  result: string;
  warnings?: string;
  error?: string;
}

export interface ErrorResponse {
  error: string;
  details?: string;
  code?: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: ErrorResponse;
}

export interface ConversionFormData {
  storeType: string;
  storeName: string;
  creationPolicy: string;
  resolve: boolean;
}

export interface YamlEditorState {
  inputYaml: string;
  outputYaml: string;
  contentHeight: string;
  error: string | null;
}

export interface AlertState {
  warning: string | null;
  errorMessage: string | null;
}

export type StoreType = 'SecretStore' | 'ClusterSecretStore';
export type CreationPolicy = 'Owner' | 'Orphan';