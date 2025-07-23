import React from 'react';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Info } from "lucide-react";
import { TIPS, STORE_TYPES, CREATION_POLICIES, REFRESH_POLICIES, PLACEHOLDERS } from '@/constants/ui';
import { useConversionStore, useFormData } from '@/stores/conversionStore';

interface ConversionFormProps { }

const ConversionForm: React.FC<ConversionFormProps> = () => {
  // Use Zustand store directly for better performance
  const formData = useFormData();
  const {
    updateStoreType,
    updateStoreName,
    updateCreationPolicy,
    updateResolve,
    updateRefreshPolicy,
    updateRefreshInterval,
  } = useConversionStore();

  const handleStoreNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateStoreName(e.target.value);
  };

  const handleRefreshIntervalChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateRefreshInterval(e.target.value);
  };

  return (
    <div className="space-y-4">
      <div>
        <label className="block mb-1 text-base">
          Store Type
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
              </TooltipTrigger>
              <TooltipContent>
                <p className="text-base">{TIPS.storeType}</p>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </label>
        <Select value={formData.storeType} onValueChange={updateStoreType}>
          <SelectTrigger className="text-base">
            <SelectValue placeholder="Select store type" />
          </SelectTrigger>
          <SelectContent>
            {STORE_TYPES.map(({ value, label }) => (
              <SelectItem key={value} value={value}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div>
        <label className="block mb-1 text-base">
          {formData.storeType} Name
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
              </TooltipTrigger>
              <TooltipContent>
                <div className="text-base whitespace-pre-line max-w-sm">{TIPS.storeName}</div>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </label>
        <Input
          type="text"
          className="text-base"
          value={formData.storeName}
          onChange={handleStoreNameChange}
          placeholder={PLACEHOLDERS.storeName}
        />
      </div>

      <div>
        <label className="block mb-1 text-base">
          Creation Policy
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
              </TooltipTrigger>
              <TooltipContent>
                <div className="text-base whitespace-pre-line max-w-sm">{TIPS.creationPolicy}</div>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </label>
        <Select value={formData.creationPolicy} onValueChange={updateCreationPolicy}>
          <SelectTrigger className="text-base">
            <SelectValue placeholder="Select creation policy" />
          </SelectTrigger>
          <SelectContent>
            {CREATION_POLICIES.map(({ value, label }) => (
              <SelectItem key={value} value={value}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div>
        <label className="block mb-1 text-base">
          Refresh Policy
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
              </TooltipTrigger>
              <TooltipContent>
                <div className="text-base whitespace-pre-line max-w-sm">{TIPS.refreshPolicy}</div>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </label>
        <Select value={formData.refreshPolicy} onValueChange={updateRefreshPolicy}>
          <SelectTrigger className="text-base">
            <SelectValue placeholder="Select refresh policy" />
          </SelectTrigger>
          <SelectContent>
            {REFRESH_POLICIES.map(({ value, label }) => (
              <SelectItem key={value} value={value}>
                {label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {formData.refreshPolicy === 'Periodic' && (
        <div>
          <label className="block mb-1 text-base">
            Refresh Interval
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
                </TooltipTrigger>
                <TooltipContent>
                  <div className="text-base whitespace-pre-line max-w-sm">{TIPS.refreshInterval}</div>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </label>
          <Input
            type="text"
            className="text-base"
            value={formData.refreshInterval}
            onChange={handleRefreshIntervalChange}
            placeholder={PLACEHOLDERS.refreshInterval}
          />
        </div>
      )}

      <div className="flex items-center space-x-2">
        <Checkbox
          id="resolve"
          checked={formData.resolve}
          onCheckedChange={updateResolve}
        />
        <label htmlFor="resolve" className="text-base">
          Resolve ENV variables
          <TooltipProvider>
            <Tooltip>
              <TooltipTrigger asChild>
                <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
              </TooltipTrigger>
              <TooltipContent>
                <div className="text-base whitespace-pre-line max-w-sm">{TIPS.resolve}</div>
              </TooltipContent>
            </Tooltip>
          </TooltipProvider>
        </label>
      </div>
    </div>
  );
};

export default ConversionForm;