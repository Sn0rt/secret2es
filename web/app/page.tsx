"use client"

import * as React from "react"
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { ArrowRight, Info } from "lucide-react"
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter'
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

type EnvVar = { key: string; value: string };

const MAX_LINES = 4080;
const LINE_HEIGHT = 20;
const DEFAULT_HEIGHT = '60vh';

const TIPS = {
  storeType: "Select the type of secret store you're using",
  storeName: "Enter the name of your secret store",
  creationPolicy: "Choose how the ExternalSecret should be created",
  resolve: "Enable to resolve environment variables in the secret",
};

export default function Home() {
  const [inputYaml, setInputYaml] = React.useState("")
  const [outputYaml, setOutputYaml] = React.useState("")
  const [storeType, setStoreType] = React.useState("SecretStore")
  const [storeName, setStoreName] = React.useState("")
  const [creationPolicy, setCreationPolicy] = React.useState("Owner")
  const [resolve, setResolve] = React.useState(false)
  const [envVars, setEnvVars] = React.useState<EnvVar[]>([{ key: '', value: '' }])
  const [error, setError] = React.useState<string | null>(null)
  const [contentHeight, setContentHeight] = React.useState(DEFAULT_HEIGHT)

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const newValue = e.target.value;
    const lineCount = newValue.split('\n').length;
    if (lineCount > MAX_LINES) {
      setError(`Input exceeds maximum limit of ${MAX_LINES} lines.`);
    } else {
      setError(null);
      setInputYaml(newValue);
      const newHeight = `${Math.max(parseInt(DEFAULT_HEIGHT), lineCount * LINE_HEIGHT)}px`;
      setContentHeight(newHeight);
    }
  }

  const handleConvert = async () => {
    if (error) return;
    try {
      const envVarsObject = envVars.reduce((acc, { key, value }) => {
        if (key) acc[key] = value;
        return acc;
      }, {} as Record<string, string>);

      const response = await fetch('http://localhost:8080/api/convert', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          content: inputYaml,
          storeType,
          storeName,
          creationPolicy,
          resolve,
          envVars: resolve ? envVarsObject : undefined,
        }),
      })
      const data = await response.json()
      setOutputYaml(data.result)
      const outputLineCount = data.result.split('\n').length;
      const newHeight = `${Math.max(parseInt(contentHeight), outputLineCount * LINE_HEIGHT)}px`;
      setContentHeight(newHeight);
    } catch (error) {
      console.error('Error converting YAML:', error)
      setOutputYaml('Error converting YAML. Please try again.')
    }
  }

  const handleAddEnvVar = () => {
    setEnvVars([...envVars, { key: '', value: '' }])
  }

  const handleEnvVarChange = (index: number, field: 'key' | 'value', value: string) => {
    const newEnvVars = [...envVars]
    newEnvVars[index][field] = value
    setEnvVars(newEnvVars)
  }

  return (
    <div className="container mx-auto p-4 flex-grow flex">
      <div className="w-1/6 pr-4 flex flex-col">
        <Button
          onClick={handleConvert}
          className="w-full px-4 py-2 text-base bg-blue-500 hover:bg-blue-600 text-white mb-4"
          disabled={!!error}
        >
          Convert
        </Button>
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
          <Select value={storeType} onValueChange={setStoreType}>
            <SelectTrigger className="text-base">
              <SelectValue placeholder="Select store type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="SecretStore">SecretStore</SelectItem>
              <SelectItem value="ClusterSecretStore">ClusterSecretStore</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div>
          <label className="block mb-1 text-base">
            Store Name
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-base">{TIPS.storeName}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </label>
          <input
            type="text"
            className="w-full p-1 border rounded text-base"
            value={storeName}
            onChange={(e) => setStoreName(e.target.value)}
            placeholder="Enter store name"
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
                  <p className="text-base">{TIPS.creationPolicy}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </label>
          <Select value={creationPolicy} onValueChange={setCreationPolicy}>
            <SelectTrigger className="text-base">
              <SelectValue placeholder="Select creation policy" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Owner">Owner</SelectItem>
              <SelectItem value="Orphan">Orphan</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-center">
          <input
            type="checkbox"
            id="resolve"
            checked={resolve}
            onChange={(e) => setResolve(e.target.checked)}
            className="mr-2"
          />
          <label htmlFor="resolve" className="text-base">
            Resolve ENV variables
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Info className="inline-block ml-1 h-5 w-5 text-gray-500 cursor-help" />
                </TooltipTrigger>
                <TooltipContent>
                  <p className="text-base">{TIPS.resolve}</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
          </label>
        </div>
        {resolve && (
          <div>
            <h3 className="font-bold mb-2 text-base">Environment Variables</h3>
            {envVars.map((envVar, index) => (
              <div key={index} className="flex flex-col mb-2">
                <input
                  type="text"
                  className="w-full p-1 border rounded mb-1 text-base"
                  placeholder="Key"
                  value={envVar.key}
                  onChange={(e) => handleEnvVarChange(index, 'key', e.target.value)}
                />
                <input
                  type="text"
                  className="w-full p-1 border rounded text-base"
                  placeholder="Value"
                  value={envVar.value}
                  onChange={(e) => handleEnvVarChange(index, 'value', e.target.value)}
                />
              </div>
            ))}
            <Button onClick={handleAddEnvVar} className="w-full text-base">Add Environment Variable</Button>
          </div>
        )}
      </div>
      <div className="w-5/6 flex items-start">
        <div className="w-[48%]">
          <div style={{ height: contentHeight, minHeight: DEFAULT_HEIGHT }} className="relative rounded-md overflow-hidden">
            <SyntaxHighlighter
              language="yaml"
              style={tomorrow}
              customStyle={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                margin: 0,
                padding: '1rem',
                resize: 'none',
                fontSize: '16px',
                overflow: 'auto',
              }}
            >
              {inputYaml || '# Enter your AVP Secret YAML here...'}
            </SyntaxHighlighter>
            <textarea
              className="absolute top-0 left-0 w-full h-full opacity-0 resize-none text-base caret-black dark:caret-white focus:opacity-100 focus:bg-white dark:focus:bg-gray-800 focus:text-black dark:focus:text-white p-4"
              value={inputYaml}
              onChange={handleInputChange}
              placeholder="Enter your AVP Secret YAML here..."
            />
          </div>
          {error && <p className="text-red-500 mt-2 text-base">{error}</p>}
        </div>
        <div className="w-[4%] flex justify-center items-start pt-2">
          <ArrowRight className="h-10 w-10 text-blue-500" />
        </div>
        <div className="w-[48%]">
          <div style={{ height: contentHeight, minHeight: DEFAULT_HEIGHT }} className="rounded-md overflow-hidden">
            <SyntaxHighlighter
              language="yaml"
              style={tomorrow}
              customStyle={{
                height: '100%',
                width: '100%',
                margin: 0,
                padding: '1rem',
                resize: 'none',
                fontSize: '16px',
                overflow: 'auto',
              }}
            >
              {outputYaml || '# Converted External Secret YAML will appear here...'}
            </SyntaxHighlighter>
          </div>
        </div>
      </div>
    </div>
  )
}
