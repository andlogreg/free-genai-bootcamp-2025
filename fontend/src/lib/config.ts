// API Configuration

// Define the API configuration interface
interface ApiConfig {
  // Whether to use mock API implementation
  useMockApi: boolean;
  
  // Backend API base URL
  baseUrl: string;
}

// Get values from environment variables if available, otherwise use defaults
const config: ApiConfig = {
  // Set this to true to use mock API, false to use real backend API
  // In a real application, this would come from an environment variable
  useMockApi: import.meta.env.VITE_USE_MOCK_API === 'true' || false,
  
  // Backend API base URL
  baseUrl: import.meta.env.VITE_API_BASE_URL || "http://localhost:3000/api",
};

export default config; 