import { create } from 'zustand';
import { persist, createJSONStorage, subscribeWithSelector, devtools } from 'zustand/middleware';
import { ConversionFormData, StoreType, CreationPolicy } from '@/types/conversion';
import { DEFAULT_VALUES } from '@/constants/api';

interface ConversionHistory {
  id: string;
  timestamp: number;
  inputYaml: string;
  outputYaml: string;
  formData: ConversionFormData;
  title?: string;
}

interface SettingsState {
  // User preferences
  theme: 'light' | 'dark' | 'system';
  autoSave: boolean;
  defaultStoreType: StoreType;
  defaultCreationPolicy: CreationPolicy;
  showLineNumbers: boolean;
  autoFormat: boolean;
  
  // Recent conversions history (max 20 items)
  recentConversions: ConversionHistory[];
  
  // Saved form state (for auto-restore)
  lastFormState?: ConversionFormData;
  lastInputYaml?: string;
  
  // UI preferences
  sidebarCollapsed: boolean;
  editorFontSize: number;
  
  // Actions - Settings
  setTheme: (theme: 'light' | 'dark' | 'system') => void;
  setAutoSave: (enabled: boolean) => void;
  setDefaultStoreType: (storeType: StoreType) => void;
  setDefaultCreationPolicy: (policy: CreationPolicy) => void;
  setShowLineNumbers: (show: boolean) => void;
  setAutoFormat: (enabled: boolean) => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  setEditorFontSize: (size: number) => void;
  
  // Actions - History
  addRecentConversion: (conversion: Omit<ConversionHistory, 'id' | 'timestamp'>) => void;
  removeRecentConversion: (id: string) => void;
  clearHistory: () => void;
  getRecentConversion: (id: string) => ConversionHistory | undefined;
  
  // Actions - Auto-save
  saveFormState: (formData: ConversionFormData, inputYaml: string) => void;
  clearSavedState: () => void;
  
  // Actions - Export/Import
  exportSettings: () => string;
  importSettings: (settingsJson: string) => boolean;
  resetToDefaults: () => void;
}

const initialSettings = {
  // User preferences
  theme: 'system' as const,
  autoSave: true,
  defaultStoreType: DEFAULT_VALUES.storeType as StoreType,
  defaultCreationPolicy: DEFAULT_VALUES.creationPolicy as CreationPolicy,
  showLineNumbers: true,
  autoFormat: true,
  
  // History
  recentConversions: [] as ConversionHistory[],
  
  // UI preferences
  sidebarCollapsed: false,
  editorFontSize: 14,
};

export const useSettingsStore = create<SettingsState>()(
  devtools(
    persist(
      subscribeWithSelector((set, get) => ({
        ...initialSettings,
        
        // Settings Actions
        setTheme: (theme: 'light' | 'dark' | 'system') => {
          set({ theme });
        },
        
        setAutoSave: (autoSave: boolean) => {
          set({ autoSave });
        },
        
        setDefaultStoreType: (defaultStoreType: StoreType) => {
          set({ defaultStoreType });
        },
        
        setDefaultCreationPolicy: (defaultCreationPolicy: CreationPolicy) => {
          set({ defaultCreationPolicy });
        },
        
        setShowLineNumbers: (showLineNumbers: boolean) => {
          set({ showLineNumbers });
        },
        
        setAutoFormat: (autoFormat: boolean) => {
          set({ autoFormat });
        },
        
        setSidebarCollapsed: (sidebarCollapsed: boolean) => {
          set({ sidebarCollapsed });
        },
        
        setEditorFontSize: (editorFontSize: number) => {
          if (editorFontSize >= 10 && editorFontSize <= 24) {
            set({ editorFontSize });
          }
        },
        
        // History Actions
        addRecentConversion: (conversion: Omit<ConversionHistory, 'id' | 'timestamp'>) => {
          const id = `conv_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
          const timestamp = Date.now();
          
          set((state) => {
            const newConversion: ConversionHistory = {
              ...conversion,
              id,
              timestamp,
              title: conversion.title || `Conversion ${new Date(timestamp).toLocaleString()}`,
            };
            
            // Keep only the latest 20 conversions
            const updatedHistory = [newConversion, ...state.recentConversions].slice(0, 20);
            
            return {
              recentConversions: updatedHistory,
            };
          });
          
          return id;
        },
        
        removeRecentConversion: (id: string) => {
          set((state) => ({
            recentConversions: state.recentConversions.filter(conv => conv.id !== id),
          }));
        },
        
        clearHistory: () => {
          set({ recentConversions: [] });
        },
        
        getRecentConversion: (id: string): ConversionHistory | undefined => {
          return get().recentConversions.find(conv => conv.id === id);
        },
        
        // Auto-save Actions
        saveFormState: (formData: ConversionFormData, inputYaml: string) => {
          if (get().autoSave) {
            set({
              lastFormState: formData,
              lastInputYaml: inputYaml,
            });
          }
        },
        
        clearSavedState: () => {
          set({
            lastFormState: undefined,
            lastInputYaml: undefined,
          });
        },
        
        // Export/Import Actions
        exportSettings: (): string => {
          const state = get();
          const exportData = {
            version: '1.0',
            settings: {
              theme: state.theme,
              autoSave: state.autoSave,
              defaultStoreType: state.defaultStoreType,
              defaultCreationPolicy: state.defaultCreationPolicy,
              showLineNumbers: state.showLineNumbers,
              autoFormat: state.autoFormat,
              sidebarCollapsed: state.sidebarCollapsed,
              editorFontSize: state.editorFontSize,
            },
            recentConversions: state.recentConversions,
            exportTimestamp: Date.now(),
          };
          
          return JSON.stringify(exportData, null, 2);
        },
        
        importSettings: (settingsJson: string): boolean => {
          try {
            const importData = JSON.parse(settingsJson);
            
            if (importData.version === '1.0' && importData.settings) {
              set((state) => ({
                ...state,
                ...importData.settings,
                recentConversions: importData.recentConversions || state.recentConversions,
              }));
              return true;
            }
            return false;
          } catch (error) {
            console.error('Failed to import settings:', error);
            return false;
          }
        },
        
        resetToDefaults: () => {
          set(initialSettings);
        },
      })),
      {
        name: 'secret2es-settings',
        storage: createJSONStorage(() => localStorage),
        // Only persist certain fields
        partialize: (state) => ({
          theme: state.theme,
          autoSave: state.autoSave,
          defaultStoreType: state.defaultStoreType,
          defaultCreationPolicy: state.defaultCreationPolicy,
          showLineNumbers: state.showLineNumbers,
          autoFormat: state.autoFormat,
          sidebarCollapsed: state.sidebarCollapsed,
          editorFontSize: state.editorFontSize,
          recentConversions: state.recentConversions,
          lastFormState: state.lastFormState,
          lastInputYaml: state.lastInputYaml,
        }),
      }
    ),
    {
      name: 'settings-store',
    }
  )
);

// Selectors for optimal performance
export const useTheme = () => useSettingsStore((state) => state.theme);
export const useAutoSave = () => useSettingsStore((state) => state.autoSave);
export const useDefaultValues = () => useSettingsStore((state) => ({
  defaultStoreType: state.defaultStoreType,
  defaultCreationPolicy: state.defaultCreationPolicy,
}));
export const useRecentConversions = () => useSettingsStore((state) => state.recentConversions);
export const useEditorSettings = () => useSettingsStore((state) => ({
  showLineNumbers: state.showLineNumbers,
  autoFormat: state.autoFormat,
  editorFontSize: state.editorFontSize,
}));
export const useUISettings = () => useSettingsStore((state) => ({
  sidebarCollapsed: state.sidebarCollapsed,
}));
export const useLastSavedState = () => useSettingsStore((state) => ({
  lastFormState: state.lastFormState,
  lastInputYaml: state.lastInputYaml,
}));