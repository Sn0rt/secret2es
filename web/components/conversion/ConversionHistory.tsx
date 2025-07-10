import React from 'react';
import { Button } from '@/components/ui/button';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog';
import { Clock, X, Download, Trash2 } from 'lucide-react';
import { useSettingsStore, useRecentConversions } from '@/stores/settingsStore';
import { useConversionStore } from '@/stores/conversionStore';

const ConversionHistory: React.FC = () => {
  const recentConversions = useRecentConversions();
  const { removeRecentConversion, clearHistory, addRecentConversion } = useSettingsStore();
  const { setFormData, setInputYaml, setOutputYaml } = useConversionStore();
  
  const [deleteAllDialogOpen, setDeleteAllDialogOpen] = React.useState(false);

  const loadConversion = (conversionId: string) => {
    const conversion = recentConversions.find(c => c.id === conversionId);
    if (conversion) {
      setFormData(conversion.formData);
      setInputYaml(conversion.inputYaml);
      setOutputYaml(conversion.outputYaml);
    }
  };

  const formatDateTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString();
  };

  const saveCurrentAsConversion = () => {
    const store = useConversionStore.getState();
    addRecentConversion({
      inputYaml: store.inputYaml,
      outputYaml: store.outputYaml,
      formData: store.formData,
      title: `Manual Save ${new Date().toLocaleString()}`,
    });
  };

  if (recentConversions.length === 0) {
    return (
      <div className="p-4 text-center text-gray-500">
        <Clock className="mx-auto h-8 w-8 mb-2" />
        <p>No conversion history yet</p>
        <p className="text-sm">Your recent conversions will appear here</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <h3 className="font-bold text-base">Conversion History</h3>
        <div className="space-x-2">
          <Button
            onClick={saveCurrentAsConversion}
            size="sm"
            variant="outline"
          >
            <Download className="h-4 w-4 mr-1" />
            Save Current
          </Button>
          <Button
            onClick={() => setDeleteAllDialogOpen(true)}
            size="sm"
            variant="outline"
            className="text-red-600 hover:text-red-700"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      <div className="space-y-2 max-h-96 overflow-y-auto">
        {recentConversions.map((conversion) => (
          <div
            key={conversion.id}
            className="border rounded-lg p-3 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
          >
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="flex items-center justify-between">
                  <h4 className="font-medium text-sm">
                    {conversion.title || 'Untitled Conversion'}
                  </h4>
                  <Button
                    onClick={() => removeRecentConversion(conversion.id)}
                    size="sm"
                    variant="ghost"
                    className="h-6 w-6 p-0 text-gray-400 hover:text-red-600"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
                
                <p className="text-xs text-gray-500 mb-2">
                  {formatDateTime(conversion.timestamp)}
                </p>
                
                <div className="text-xs space-y-1">
                  <p>
                    <span className="font-medium">Store:</span> {conversion.formData.storeType} / {conversion.formData.storeName}
                  </p>
                  <p>
                    <span className="font-medium">Policy:</span> {conversion.formData.creationPolicy}
                  </p>
                  {conversion.formData.resolve && (
                    <p className="text-blue-600">
                      <span className="font-medium">Env Resolution:</span> Enabled
                    </p>
                  )}
                </div>
              </div>
            </div>
            
            <Button
              onClick={() => loadConversion(conversion.id)}
              className="w-full mt-3"
              size="sm"
              variant="outline"
            >
              Load This Conversion
            </Button>
          </div>
        ))}
      </div>

      <AlertDialog open={deleteAllDialogOpen} onOpenChange={setDeleteAllDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Clear All History</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete all conversion history? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => {
                clearHistory();
                setDeleteAllDialogOpen(false);
              }}
              className="bg-red-600 hover:bg-red-700"
            >
              Delete All
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
};

export default ConversionHistory;