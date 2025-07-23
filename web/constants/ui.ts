export const MAX_LINES = 4080;
export const LINE_HEIGHT = 20;
export const DEFAULT_HEIGHT = '60vh';

export const TIPS = {
  storeType: "Select the type of secret store you're using",
  storeName: "Enter the name of the SecretStore (namespace-scoped) or\nClusterSecretStore (cluster-scoped) based on the\nStore Type selection above",
  creationPolicy: "Owner: The ExternalSecret owns the created Secret.\nWhen the ExternalSecret is deleted, the Secret is also deleted.\n\nOrphan: The created Secret is independent.\nWhen the ExternalSecret is deleted, the Secret remains.",
  resolve: "Enable to resolve environment variables in the secret",
  refreshPolicy: "Choose when the secret should be refreshed.\n\nLearn more about the 3 different refresh policies:\nhttps://external-secrets.io/latest/api/externalsecret/#update-behavior-with-3-different-refresh-policies",
  refreshInterval: "How often to refresh the secret (only for Periodic policy)",
};

export const STORE_TYPES = [
  { value: 'SecretStore', label: 'SecretStore' },
  { value: 'ClusterSecretStore', label: 'ClusterSecretStore' },
];

export const CREATION_POLICIES = [
  { value: 'Owner', label: 'Owner' },
  { value: 'Orphan', label: 'Orphan' },
];

export const REFRESH_POLICIES = [
  { value: 'OnChange', label: 'OnChange' },
  { value: 'Periodic', label: 'Periodic' },
  { value: 'CreatedOnce', label: 'CreatedOnce' },
];

export const PLACEHOLDERS = {
  input: '# Enter your AVP Secret YAML here...',
  output: '# Converted External Secret YAML will appear here...',
  storeName: 'Enter store name',
  envKey: 'Key',
  envValue: 'Value',
  refreshInterval: '1h',
};