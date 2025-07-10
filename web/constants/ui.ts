export const MAX_LINES = 4080;
export const LINE_HEIGHT = 20;
export const DEFAULT_HEIGHT = '60vh';

export const TIPS = {
  storeType: "Select the type of secret store you're using",
  storeName: "Enter the name of your secret store",
  creationPolicy: "Choose how the ExternalSecret should be created",
  resolve: "Enable to resolve environment variables in the secret",
};

export const STORE_TYPES = [
  { value: 'SecretStore', label: 'SecretStore' },
  { value: 'ClusterSecretStore', label: 'ClusterSecretStore' },
];

export const CREATION_POLICIES = [
  { value: 'Owner', label: 'Owner' },
  { value: 'Orphan', label: 'Orphan' },
];

export const PLACEHOLDERS = {
  input: '# Enter your AVP Secret YAML here...',
  output: '# Converted External Secret YAML will appear here...',
  storeName: 'Enter store name',
  envKey: 'Key',
  envValue: 'Value',
};