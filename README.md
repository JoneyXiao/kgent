# KGent - Kubernetes AI Assistant

KGent is a smart Kubernetes CLI assistant powered by AI that helps you create, list, and delete Kubernetes resources using natural language. It acts as a bridge between your natural language input and the Kubernetes API.

## Features

- **Natural Language Interface**: Interact with your Kubernetes cluster using everyday language
- **Resource Creation**: Generate YAML files for Kubernetes resources based on your description
- **Resource Management**: List and delete resources through conversation
- **AI-Powered**: Uses large language models to understand requests and generate responses

## Prerequisites

- Go 1.21.3 or later
- Access to a Kubernetes cluster
- DashScope API token (or other compatible API token)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/kgent.git
   cd kgent
   ```

2. Build the project:
   ```bash
   go build -o kgent
   ```

3. Set up your environment:
   
   **Option 1: Using .env file (recommended)**
   
   Copy the example environment file and edit it with your credentials:
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and configuration
   ```
   
   **Option 2: Using environment variables directly**
   ```bash
   export DASH_SCOPE_API_KEY=your_api_token_here # required
   export DASH_SCOPE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
   export DASH_SCOPE_MODEL=qwen-turbo
   ```

## Usage

### Basic Chat

Start a conversation with the Kubernetes assistant:

```bash
./kgent chat
```

### Using a Specific Namespace

You can specify a default namespace for all operations:

```bash
./kgent chat --namespace default
```

### Example Conversations

- Creating a pod:
  ```
  > Create a pod named nginx-pod with nginx image
  ```

- Listing resources:
  ```
  > List all pods in the default namespace
  ```

- Deleting a resource:
  ```
  > Delete the pod named nginx-pod
  ```

## Configuration

You can configure KGent through either environment variables or a .env file:

| Environment Variable | Description | Default Value |
|----------------------|-------------|---------------|
| DASH_SCOPE_API_KEY   | DashScope API Key | (required) |
| DASH_SCOPE_URL       | DashScope API URL | https://dashscope.aliyuncs.com/compatible-mode/v1 |
| DASH_SCOPE_MODEL     | AI Model to use    | qwen-max |

## License

This project is licensed under the MIT License - see the LICENSE file for details. 