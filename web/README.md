# Secret2ES Web Interface

This is the web interface for the Secret2ES project, a tool designed to convert Kubernetes Secrets to External Secrets.

## Overview

The Secret2ES Web Interface provides a user-friendly way to interact with the Secret2ES conversion tool. It allows users to input Kubernetes Secret YAML and receive the corresponding External Secret YAML without needing to use the command-line interface.

## Features

- Web-based interface for Secret to External Secret conversion
- Real-time conversion without the need for local installation
- Support for all options available in the CLI version

## Getting Started

### Prerequisites

- A modern web browser
- Access to the Secret2ES server (usually running on `http://localhost:8080`)

### Usage

1. Open the web interface in your browser.
2. Paste your Kubernetes Secret YAML into the input field.
3. Configure the conversion options as needed.
4. Click the "Convert" button.
5. The converted External Secret YAML will be displayed in the output field.

## Development

This web interface is built using [insert technology stack here, e.g., React, Vue.js, etc.]. To set up the development environment:

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
