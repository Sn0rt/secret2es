import { create } from 'zustand';
import { subscribeWithSelector, devtools } from 'zustand/middleware';
import { ConversionFormData, EnvVar, ConversionRequest } from '@/types/conversion';
import { DEFAULT_VALUES } from '@/constants/api';
import { DEFAULT_HEIGHT } from '@/constants/ui';
import { validateFormData, validateEnvVars } from '@/utils/validation';
import { validateYamlInput, calculateContentHeight } from '@/utils/yaml';
import { API_ENDPOINTS, HTTP_METHODS, CONTENT_TYPES } from '@/constants/api';
import { parseErrorResponse, handleFetchError, ApiError, getUserFriendlyMessage } from '@/utils/error-handling';
// import { useSettingsStore } from './settingsStore';

interface ConversionState {
  // Form state
  formData: ConversionFormData;
  isFormValid: boolean;
  
  // YAML state
  inputYaml: string;
  outputYaml: string;
  contentHeight: string;
  yamlError: string | null;
  
  // Environment variables
  envVars: EnvVar[];
  
  // Conversion state
  isLoading: boolean;
  
  // Alert state
  warning: string | null;
  errorMessage: string | null;
  
  // Actions - Form
  setFormData: (data: Partial<ConversionFormData>) => void;
  updateStoreType: (storeType: string) => void;
  updateStoreName: (storeName: string) => void;
  updateCreationPolicy: (creationPolicy: string) => void;
  updateResolve: (resolve: boolean) => void;
  
  // Actions - YAML
  setInputYaml: (yaml: string) => void;
  setOutputYaml: (yaml: string) => void;
  clearOutput: () => void;
  
  // Actions - Environment Variables
  addEnvVar: () => void;
  updateEnvVar: (index: number, field: 'key' | 'value', value: string) => void;
  removeEnvVar: (index: number) => void;
  getEnvVarsAsRecord: () => Record<string, string>;
  
  // Actions - Alerts
  showWarning: (message: string) => void;
  showError: (message: string) => void;
  clearWarning: () => void;
  clearError: () => void;
  clearAllAlerts: () => void;
  
  // Actions - Conversion
  convertYaml: () => Promise<void>;
  
  // Actions - Utilities
  validateForm: () => void;
  resetState: () => void;
}

const initialState = {
  // Form state
  formData: {
    storeType: DEFAULT_VALUES.storeType,
    storeName: '',
    creationPolicy: DEFAULT_VALUES.creationPolicy,
    resolve: DEFAULT_VALUES.resolve,
  },
  isFormValid: false,
  
  // YAML state
  inputYaml: '',
  outputYaml: '',
  contentHeight: DEFAULT_HEIGHT,
  yamlError: null,
  
  // Environment variables
  envVars: [{ key: '', value: '' }] as EnvVar[],
  
  // Conversion state
  isLoading: false,
  
  // Alert state
  warning: null,
  errorMessage: null,
};

export const useConversionStore = create<ConversionState>()(
  devtools(
    subscribeWithSelector((set, get) => ({
      ...initialState,
      
      // Form Actions
      setFormData: (data: Partial<ConversionFormData>) => {
        set((state) => {
          const newFormData = { ...state.formData, ...data };
          return {
            formData: newFormData,
            isFormValid: validateFormData(newFormData.storeName, state.yamlError),
          };
        });
      },
      
      updateStoreType: (storeType: string) => {
        set((state) => {
          const newFormData = { ...state.formData, storeType };
          return {
            formData: newFormData,
            isFormValid: validateFormData(newFormData.storeName, state.yamlError),
          };
        });
      },
      
      updateStoreName: (storeName: string) => {
        set((state) => {
          const newFormData = { ...state.formData, storeName };
          return {
            formData: newFormData,
            isFormValid: validateFormData(newFormData.storeName, state.yamlError),
          };
        });
      },
      
      updateCreationPolicy: (creationPolicy: string) => {
        set((state) => {
          const newFormData = { ...state.formData, creationPolicy };
          return {
            formData: newFormData,
            isFormValid: validateFormData(newFormData.storeName, state.yamlError),
          };
        });
      },
      
      updateResolve: (resolve: boolean) => {
        set((state) => {
          const newFormData = { ...state.formData, resolve };
          return {
            formData: newFormData,
            isFormValid: validateFormData(newFormData.storeName, state.yamlError),
          };
        });
      },
      
      // YAML Actions
      setInputYaml: (yaml: string) => {
        const validationError = validateYamlInput(yaml);
        
        if (validationError) {
          set({ yamlError: validationError });
        } else {
          const newHeight = calculateContentHeight(yaml);
          set((state) => ({
            inputYaml: yaml,
            yamlError: null,
            contentHeight: newHeight,
            isFormValid: validateFormData(state.formData.storeName, null),
          }));
        }
      },
      
      setOutputYaml: (yaml: string) => {
        const newHeight = calculateContentHeight(yaml);
        set({
          outputYaml: yaml,
          contentHeight: newHeight,
        });
      },
      
      clearOutput: () => {
        set({ outputYaml: '' });
      },
      
      // Environment Variable Actions
      addEnvVar: () => {
        set((state) => ({
          envVars: [...state.envVars, { key: '', value: '' }],
        }));
      },
      
      updateEnvVar: (index: number, field: 'key' | 'value', value: string) => {
        set((state) => {
          const newEnvVars = [...state.envVars];
          newEnvVars[index][field] = value;
          return { envVars: newEnvVars };
        });
      },
      
      removeEnvVar: (index: number) => {
        set((state) => ({
          envVars: state.envVars.filter((_, i) => i !== index),
        }));
      },
      
      getEnvVarsAsRecord: (): Record<string, string> => {
        return validateEnvVars(get().envVars);
      },
      
      // Alert Actions
      showWarning: (message: string) => {
        set({ warning: message });
      },
      
      showError: (message: string) => {
        set({ errorMessage: message });
      },
      
      clearWarning: () => {
        set({ warning: null });
      },
      
      clearError: () => {
        set({ errorMessage: null });
      },
      
      clearAllAlerts: () => {
        set({ warning: null, errorMessage: null });
      },
      
      // Conversion Action
      convertYaml: async () => {
        const state = get();
        
        if (!state.isFormValid) {
          return;
        }
        
        set({ isLoading: true });
        get().clearAllAlerts();
        
        try {
          const envVarsObject = state.formData.resolve ? state.getEnvVarsAsRecord() : undefined;
          
          const request: ConversionRequest = {
            content: state.inputYaml,
            storeType: state.formData.storeType,
            storeName: state.formData.storeName,
            creationPolicy: state.formData.creationPolicy,
            resolve: state.formData.resolve,
            envVars: envVarsObject,
          };
          
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
          
          get().setOutputYaml(data.result);
          
          if (data.warnings) {
            get().showWarning(data.warnings);
          }
          
          // Temporarily disable cross-store communication to isolate infinite loop
          // const settingsStore = useSettingsStore.getState();
          // settingsStore.addRecentConversion({
          //   inputYaml: state.inputYaml,
          //   outputYaml: data.result,
          //   formData: state.formData,
          // });
        } catch (error) {
          let errorMsg: string;
          
          if (error instanceof ApiError) {
            errorMsg = getUserFriendlyMessage(error);
            console.error('Conversion failed:', {
              message: error.message,
              status: error.status,
              details: error.details,
              code: error.code,
            });
          } else {
            const apiError = handleFetchError(error);
            errorMsg = getUserFriendlyMessage(apiError);
            console.error('Conversion Error:', {
              message: apiError.message,
              status: apiError.status,
              details: apiError.details,
              code: apiError.code,
              originalError: error,
            });
          }
          
          get().showError(errorMsg);
          get().clearOutput();
        } finally {
          set({ isLoading: false });
        }
      },
      
      // Utility Actions
      validateForm: () => {
        const state = get();
        set({
          isFormValid: validateFormData(state.formData.storeName, state.yamlError),
        });
      },
      
      resetState: () => {
        set(initialState);
      },
    })),
    {
      name: 'conversion-store',
    }
  )
);

// Selectors for optimal performance with SSR support
export const useFormData = () => useConversionStore((state) => state.formData);
export const useIsFormValid = () => useConversionStore((state) => state.isFormValid);

// Individual selectors to avoid object creation on every render
export const useInputYaml = () => useConversionStore((state) => state.inputYaml);
export const useOutputYaml = () => useConversionStore((state) => state.outputYaml);
export const useContentHeight = () => useConversionStore((state) => state.contentHeight);
export const useYamlError = () => useConversionStore((state) => state.yamlError);

export const useEnvVars = () => useConversionStore((state) => state.envVars);
export const useIsLoading = () => useConversionStore((state) => state.isLoading);

// Individual alert selectors
export const useWarning = () => useConversionStore((state) => state.warning);
export const useErrorMessage = () => useConversionStore((state) => state.errorMessage);