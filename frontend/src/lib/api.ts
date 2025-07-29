// API Configuration and Types for Proompt Backend

const API_BASE_URL = 'http://localhost:8080/api';

// Core Types matching the Go backend models
export interface Prompt {
  id: string;
  title: string;
  content: string;
  type: 'system' | 'user' | 'image' | 'video';
  use_case?: string;
  model_compatibility_tags: string[];
  temperature_suggestion?: number;
  other_parameters?: Record<string, any>;
  created_at: string;
  updated_at: string;
  git_ref?: string;
}

export interface Snippet {
  id: string;
  title: string;
  content: string;
  description?: string;
  created_at: string;
  updated_at: string;
  git_ref?: string;
  tags?: string[]; // Added missing tags property
  category?: string; // Added missing category property
}

export interface Note {
  id: string;
  prompt_id: string;
  title: string;
  body?: string;
  created_at: string;
  updated_at: string;
}

export interface PromptLink {
  from_prompt_id: string;
  to_prompt_id: string;
  link_type: string;
  created_at: string;
}

// API Request/Response Types
export interface CreatePromptRequest {
  title: string;
  content: string;
  type: 'system' | 'user' | 'image' | 'video';
  use_case: string;
  model_compatibility_tags?: string[];
  temperature_suggestion?: number;
  other_parameters?: Record<string, any>;
  notes?: string;
}

export interface UpdatePromptRequest {
  title?: string;
  content?: string;
  type?: 'system' | 'user' | 'image' | 'video';
  use_case?: string;
  model_compatibility_tags?: string[];
  temperature_suggestion?: number;
  other_parameters?: Record<string, any>;
  notes?: string;
}

export interface CreateSnippetRequest {
  title: string;
  content: string;
  description?: string;
}

export interface UpdateSnippetRequest {
  title?: string;
  content?: string;
  description?: string;
}

export interface CreateNoteRequest {
  title: string;
  body: string;
}

export interface UpdateNoteRequest {
  title?: string;
  body?: string;
}

export interface TemplatePreviewRequest {
  content: string;
  variables?: Record<string, string>;
}

export interface TemplateVariable {
  name: string;
  default_value?: string;
  has_default: boolean;
  status: 'provided' | 'default' | 'missing';
}

export interface TemplatePreviewResponse {
  resolved_content: string;
  variables: TemplateVariable[];
  warnings: string[];
}

export interface ListResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface PromptFilters {
  type?: string;
  use_case?: string;
  limit?: number;
  offset?: number;
}

export interface SnippetFilters {
  limit?: number;
  offset?: number;
}

// API Error Class
export class APIError extends Error {
  public error: string;
  public code: number;
  public details?: Record<string, string>;

  constructor(error: string, message: string, code: number, details?: Record<string, string>) {
    super(message);
    this.name = 'APIError';
    this.error = error;
    this.code = code;
    this.details = details;
  }
}

export interface ValidationError {
  error: string;
  message: string;
  code: number;
  fields: Record<string, string>;
}

// HTTP Client with error handling
class APIClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new APIError(
          errorData.error || 'Request failed',
          errorData.message || `HTTP ${response.status}`,
          response.status,
          errorData.details
        );
      }

      return await response.json();
    } catch (error) {
      if (error instanceof APIError) {
        throw error;
      }
      throw new APIError(
        'Network Error',
        'Failed to connect to the server',
        0
      );
    }
  }

  async get<T>(endpoint: string, params?: Record<string, string>): Promise<T> {
    const url = params 
      ? `${endpoint}?${new URLSearchParams(params).toString()}`
      : endpoint;
    
    return this.request<T>(url, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }
}

// Create API client instance
const apiClient = new APIClient(API_BASE_URL);

// Prompt API
export const promptAPI = {
  list: (filters?: PromptFilters) => {
    const params: Record<string, string> = {};
    if (filters?.type) params.type = filters.type;
    if (filters?.use_case) params.use_case = filters.use_case;
    if (filters?.limit) params.limit = filters.limit.toString();
    if (filters?.offset) params.offset = filters.offset.toString();
    
    return apiClient.get<ListResponse<Prompt>>('/prompts', params);
  },

  create: (data: CreatePromptRequest) =>
    apiClient.post<Prompt>('/prompts', data),

  get: (id: string) =>
    apiClient.get<Prompt>(`/prompts/${id}`),

  update: (id: string, data: UpdatePromptRequest) =>
    apiClient.put<Prompt>(`/prompts/${id}`, data),

  delete: (id: string) =>
    apiClient.delete<void>(`/prompts/${id}`),

  // Links
  createLink: (fromId: string, toId: string, linkType = 'followup') =>
    apiClient.post<void>(`/prompts/${fromId}/links`, { 
      to_prompt_id: toId, 
      link_type: linkType 
    }),

  deleteLink: (fromId: string, toId: string) =>
    apiClient.delete<void>(`/prompts/${fromId}/links/${toId}`),

  getLinks: (id: string) =>
    apiClient.get<PromptLink[]>(`/prompts/${id}/links`),

  getBacklinks: (id: string) =>
    apiClient.get<PromptLink[]>(`/prompts/${id}/backlinks`),

  // Tags
  addTag: (id: string, tagName: string) =>
    apiClient.post<void>(`/prompts/${id}/tags`, { tag_name: tagName }),

  removeTag: (id: string, tagName: string) =>
    apiClient.delete<void>(`/prompts/${id}/tags/${tagName}`),

  getTags: (id: string) =>
    apiClient.get<{ tags: string[] }>(`/prompts/${id}/tags`),

  getAllTags: () =>
    apiClient.get<{ tags: string[] }>('/prompts/tags'),
};

// Snippet API
export const snippetAPI = {
  list: (filters?: SnippetFilters) => {
    const params: Record<string, string> = {};
    if (filters?.limit) params.limit = filters.limit.toString();
    if (filters?.offset) params.offset = filters.offset.toString();
    
    return apiClient.get<ListResponse<Snippet>>('/snippets', params);
  },

  create: (data: CreateSnippetRequest) =>
    apiClient.post<Snippet>('/snippets', data),

  get: (id: string) =>
    apiClient.get<Snippet>(`/snippets/${id}`),

  update: (id: string, data: UpdateSnippetRequest) =>
    apiClient.put<Snippet>(`/snippets/${id}`, data),

  delete: (id: string) =>
    apiClient.delete<void>(`/snippets/${id}`),

  // Tags
  addTag: (id: string, tagName: string) =>
    apiClient.post<void>(`/snippets/${id}/tags`, { tag_name: tagName }),

  removeTag: (id: string, tagName: string) =>
    apiClient.delete<void>(`/snippets/${id}/tags/${tagName}`),

  getTags: (id: string) =>
    apiClient.get<{ tags: string[] }>(`/snippets/${id}/tags`),

  getAllTags: () =>
    apiClient.get<{ tags: string[] }>('/snippets/tags'),
};

// Notes API
export const noteAPI = {
  listForPrompt: (promptId: string) =>
    apiClient.get<ListResponse<Note>>(`/prompts/${promptId}/notes`),

  create: (promptId: string, data: CreateNoteRequest) =>
    apiClient.post<Note>(`/prompts/${promptId}/notes`, data),

  get: (id: string) =>
    apiClient.get<Note>(`/notes/${id}`),

  update: (id: string, data: UpdateNoteRequest) =>
    apiClient.put<Note>(`/notes/${id}`, data),

  delete: (id: string) =>
    apiClient.delete<void>(`/notes/${id}`),
};

// Template API
export const templateAPI = {
  preview: (data: TemplatePreviewRequest) =>
    apiClient.post<TemplatePreviewResponse>('/template/preview', data),

  analyze: (data: TemplatePreviewRequest) =>
    apiClient.post<TemplatePreviewResponse>('/template/analyze', data),
};

// Health API
export const healthAPI = {
  check: () =>
    apiClient.get<{ status: string; timestamp: string; version: string }>('/health'),
};

// Utility function to check if error is an API error
export const isAPIError = (error: unknown): error is APIError => {
  return error instanceof APIError;
}; 