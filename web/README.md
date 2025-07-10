# Secret2ES Web UI

A modern, componentized React application for converting ArgoCD Vault Plugin secrets to External Secrets format.

## Architecture Overview

This application has been fully refactored with a clean, componentized architecture featuring:

- **Modular Components**: Separated UI logic into reusable components
- **Custom Hooks**: Encapsulated business logic in custom React hooks  
- **Type Safety**: Full TypeScript integration with proper type definitions
- **Clean Separation**: Clear separation between UI, business logic, and utilities

## 📁 Directory Structure

```
web/
├── app/                        # Next.js app directory
│   ├── page.tsx               # Main page (componentized)
│   └── layout.tsx             # Layout wrapper
├── components/                 # React components
│   ├── ui/                    # Base UI components (Shadcn/ui)
│   │   ├── button.tsx
│   │   ├── input.tsx
│   │   ├── checkbox.tsx
│   │   ├── select.tsx
│   │   ├── alert.tsx
│   │   └── loading.tsx
│   └── conversion/            # Conversion-specific components
│       ├── ConversionForm.tsx
│       ├── EnvironmentVariables.tsx
│       ├── YamlEditor.tsx
│       ├── ConversionResult.tsx
│       ├── AlertMessages.tsx
│       └── ConversionButton.tsx
├── hooks/                     # Custom React hooks
│   ├── useConversion.ts
│   ├── useEnvironmentVariables.ts
│   ├── useYamlEditor.ts
│   └── useAlerts.ts
├── types/                     # TypeScript type definitions
│   └── conversion.ts
├── constants/                 # Application constants
│   ├── ui.ts
│   └── api.ts
└── utils/                     # Utility functions
    ├── yaml.ts
    ├── clipboard.ts
    └── validation.ts
```

## 🎯 Key Components

### Core Components
- **ConversionForm**: Handles form inputs (store type, name, creation policy)
- **EnvironmentVariables**: Manages environment variable inputs
- **YamlEditor**: Syntax-highlighted YAML editor with overlay input
- **ConversionResult**: Output display with copy-to-clipboard functionality
- **AlertMessages**: Error and warning notifications
- **ConversionButton**: Convert button with loading states

### Custom Hooks
- **useConversion**: Handles API calls and conversion logic
- **useEnvironmentVariables**: Manages environment variable state
- **useYamlEditor**: Handles YAML editing, validation, and height management
- **useAlerts**: Manages error and warning messages

## 🚀 Features

- **Modular Architecture**: Clean separation of concerns
- **Type Safety**: Full TypeScript integration
- **Reusable Components**: Components can be easily reused
- **Custom Hooks**: Encapsulated business logic
- **Error Handling**: Comprehensive error and warning system
- **Loading States**: Proper loading indicators
- **Accessibility**: ARIA-compliant components
- **Performance**: Optimized with React best practices

## 🛠 Technical Stack

- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Library**: Radix UI primitives
- **Icons**: Lucide React
- **Syntax Highlighting**: Prism.js via react-syntax-highlighter
- **Notifications**: React Hot Toast

## Development

1. Clone the repository:
   ```
   git clone https://github.com/Sn0rt/sercert2extsecret.git
   ```
2. Navigate to the web directory:
   ```
   cd sercert2extsecret/web
   ```
3. Install dependencies:
   ```
   npm install
   ```
4. Start the development server:
   ```
   npm run dev
   ```

## Contributing

Contributions to the Secret2ES Web Interface are welcome! Please refer to the main project's contributing guidelines for more information.

## License

This project is licensed under the same terms as the main Secret2ES project. Please refer to the LICENSE file in the root directory for more information.

## Additional Information

For more details about the Secret2ES project, including its CLI usage and backend implementation, please refer to the [main README](../README.md) in the project root.
