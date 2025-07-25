// TanStack Query hooks for Proompt API

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { 
  promptAPI, 
  snippetAPI, 
  noteAPI, 
  templateAPI, 
  healthAPI,
  type Prompt,
  type Snippet,
  type Note,
  type CreatePromptRequest,
  type UpdatePromptRequest,
  type CreateSnippetRequest,
  type UpdateSnippetRequest,
  type CreateNoteRequest,
  type UpdateNoteRequest,
  type TemplatePreviewRequest,
  type PromptFilters,
  type SnippetFilters,
  isAPIError
} from './api';
import { toast } from '@/hooks/use-toast';

// Query Keys
export const queryKeys = {
  health: ['health'] as const,
  prompts: {
    all: ['prompts'] as const,
    lists: () => [...queryKeys.prompts.all, 'list'] as const,
    list: (filters?: PromptFilters) => [...queryKeys.prompts.lists(), filters] as const,
    details: () => [...queryKeys.prompts.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.prompts.details(), id] as const,
    tags: () => [...queryKeys.prompts.all, 'tags'] as const,
    allTags: () => [...queryKeys.prompts.tags(), 'all'] as const,
  },
  snippets: {
    all: ['snippets'] as const,
    lists: () => [...queryKeys.snippets.all, 'list'] as const,
    list: (filters?: SnippetFilters) => [...queryKeys.snippets.lists(), filters] as const,
    details: () => [...queryKeys.snippets.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.snippets.details(), id] as const,
    tags: () => [...queryKeys.snippets.all, 'tags'] as const,
    allTags: () => [...queryKeys.snippets.tags(), 'all'] as const,
  },
  notes: {
    all: ['notes'] as const,
    byPrompt: (promptId: string) => [...queryKeys.notes.all, 'byPrompt', promptId] as const,
    detail: (id: string) => [...queryKeys.notes.all, 'detail', id] as const,
  },
  template: {
    all: ['template'] as const,
    preview: (request: TemplatePreviewRequest) => [...queryKeys.template.all, 'preview', request] as const,
  },
} as const;

// Health Check
export const useHealthCheck = () => {
  return useQuery({
    queryKey: queryKeys.health,
    queryFn: healthAPI.check,
    staleTime: 10000, // 10 seconds
    refetchInterval: 30000, // Poll every 30 seconds
    refetchIntervalInBackground: true, // Continue polling in background
    retry: (failureCount, error) => {
      // Keep retrying indefinitely for network errors
      if (isAPIError(error) && error.message.includes('NetworkError')) {
        return true;
      }
      // For other errors, stop after 3 attempts
      return failureCount < 3;
    },
    retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000), // Exponential backoff
  });
};

// Prompt Queries
export const usePrompts = (filters?: PromptFilters) => {
  return useQuery({
    queryKey: queryKeys.prompts.list(filters),
    queryFn: () => promptAPI.list(filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const usePrompt = (id: string, enabled = true) => {
  return useQuery({
    queryKey: queryKeys.prompts.detail(id),
    queryFn: () => promptAPI.get(id),
    enabled: enabled && !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const usePromptTags = () => {
  return useQuery({
    queryKey: queryKeys.prompts.allTags(),
    queryFn: promptAPI.getAllTags,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

// Prompt Mutations
export const useCreatePrompt = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreatePromptRequest) => promptAPI.create(data),
    onSuccess: (newPrompt) => {
      // Invalidate and refetch prompt lists
      queryClient.invalidateQueries({ queryKey: queryKeys.prompts.lists() });
      
      // Add the new prompt to cache
      queryClient.setQueryData(
        queryKeys.prompts.detail(newPrompt.id),
        newPrompt
      );
      
      toast({
        title: "Prompt created",
        description: `"${newPrompt.title}" has been created successfully.`,
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to create prompt';
      toast({
        title: "Error creating prompt",
        description: message,
        variant: "destructive",
      });
    },
  });
};

export const useUpdatePrompt = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePromptRequest }) =>
      promptAPI.update(id, data),
    onSuccess: (updatedPrompt) => {
      // Update the prompt in cache
      queryClient.setQueryData(
        queryKeys.prompts.detail(updatedPrompt.id),
        updatedPrompt
      );
      
      // Invalidate lists to reflect changes
      queryClient.invalidateQueries({ queryKey: queryKeys.prompts.lists() });
      
      toast({
        title: "Prompt updated",
        description: `"${updatedPrompt.title}" has been updated successfully.`,
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to update prompt';
      toast({
        title: "Error updating prompt",
        description: message,
        variant: "destructive",
      });
    },
  });
};

export const useDeletePrompt = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => promptAPI.delete(id),
    onSuccess: (_, deletedId) => {
      // Remove from cache
      queryClient.removeQueries({ queryKey: queryKeys.prompts.detail(deletedId) });
      
      // Invalidate lists
      queryClient.invalidateQueries({ queryKey: queryKeys.prompts.lists() });
      
      toast({
        title: "Prompt deleted",
        description: "The prompt has been deleted successfully.",
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to delete prompt';
      toast({
        title: "Error deleting prompt",
        description: message,
        variant: "destructive",
      });
    },
  });
};

// Snippet Queries
export const useSnippets = (filters?: SnippetFilters) => {
  return useQuery({
    queryKey: queryKeys.snippets.list(filters),
    queryFn: () => snippetAPI.list(filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useSnippet = (id: string, enabled = true) => {
  return useQuery({
    queryKey: queryKeys.snippets.detail(id),
    queryFn: () => snippetAPI.get(id),
    enabled: enabled && !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useSnippetTags = () => {
  return useQuery({
    queryKey: queryKeys.snippets.allTags(),
    queryFn: snippetAPI.getAllTags,
    staleTime: 10 * 60 * 1000, // 10 minutes
  });
};

// Snippet Mutations
export const useCreateSnippet = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: CreateSnippetRequest) => snippetAPI.create(data),
    onSuccess: (newSnippet) => {
      queryClient.invalidateQueries({ queryKey: queryKeys.snippets.lists() });
      queryClient.setQueryData(
        queryKeys.snippets.detail(newSnippet.id),
        newSnippet
      );
      
      toast({
        title: "Snippet created",
        description: `"${newSnippet.title}" has been created successfully.`,
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to create snippet';
      toast({
        title: "Error creating snippet",
        description: message,
        variant: "destructive",
      });
    },
  });
};

export const useUpdateSnippet = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateSnippetRequest }) =>
      snippetAPI.update(id, data),
    onSuccess: (updatedSnippet) => {
      queryClient.setQueryData(
        queryKeys.snippets.detail(updatedSnippet.id),
        updatedSnippet
      );
      queryClient.invalidateQueries({ queryKey: queryKeys.snippets.lists() });
      
      toast({
        title: "Snippet updated",
        description: `"${updatedSnippet.title}" has been updated successfully.`,
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to update snippet';
      toast({
        title: "Error updating snippet",
        description: message,
        variant: "destructive",
      });
    },
  });
};

export const useDeleteSnippet = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (id: string) => snippetAPI.delete(id),
    onSuccess: (_, deletedId) => {
      queryClient.removeQueries({ queryKey: queryKeys.snippets.detail(deletedId) });
      queryClient.invalidateQueries({ queryKey: queryKeys.snippets.lists() });
      
      toast({
        title: "Snippet deleted",
        description: "The snippet has been deleted successfully.",
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to delete snippet';
      toast({
        title: "Error deleting snippet",
        description: message,
        variant: "destructive",
      });
    },
  });
};

// Note Queries
export const useNotesForPrompt = (promptId: string) => {
  return useQuery({
    queryKey: queryKeys.notes.byPrompt(promptId),
    queryFn: () => noteAPI.listForPrompt(promptId),
    enabled: !!promptId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

export const useNote = (id: string, enabled = true) => {
  return useQuery({
    queryKey: queryKeys.notes.detail(id),
    queryFn: () => noteAPI.get(id),
    enabled: enabled && !!id,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Note Mutations
export const useCreateNote = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: ({ promptId, data }: { promptId: string; data: CreateNoteRequest }) =>
      noteAPI.create(promptId, data),
    onSuccess: (newNote) => {
      queryClient.invalidateQueries({ 
        queryKey: queryKeys.notes.byPrompt(newNote.prompt_id) 
      });
      queryClient.setQueryData(
        queryKeys.notes.detail(newNote.id),
        newNote
      );
      
      toast({
        title: "Note created",
        description: `"${newNote.title}" has been created successfully.`,
      });
    },
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to create note';
      toast({
        title: "Error creating note",
        description: message,
        variant: "destructive",
      });
    },
  });
};

// Template Queries
export const useTemplatePreview = (request: TemplatePreviewRequest, enabled = true) => {
  return useQuery({
    queryKey: queryKeys.template.preview(request),
    queryFn: () => templateAPI.preview(request),
    enabled: enabled && !!request.content,
    staleTime: 0, // Always fresh for template previews
    gcTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Template Mutations (for one-time previews)
export const useTemplatePreviewMutation = () => {
  return useMutation({
    mutationFn: (data: TemplatePreviewRequest) => templateAPI.preview(data),
    retry: (failureCount, error) => {
      // Retry network errors up to 2 times
      if (isAPIError(error) && (
        error.message.includes('NetworkError') || 
        error.message.includes('fetch') ||
        error.message.includes('Failed to fetch')
      )) {
        return failureCount < 2;
      }
      return false;
    },
    retryDelay: 1000, // 1 second between retries
    onError: (error) => {
      let title = "Preview error";
      let description = 'Failed to generate preview';
      
      if (isAPIError(error)) {
        if (error.message.includes('NetworkError') || 
            error.message.includes('fetch') || 
            error.message.includes('Failed to fetch')) {
          title = "Connection error";
          description = "Unable to connect to backend server. Please check if the server is running.";
        } else {
          description = error.message;
        }
      }
      
      toast({
        title,
        description,
        variant: "destructive",
      });
    },
  });
};

export const useTemplateAnalyzeMutation = () => {
  return useMutation({
    mutationFn: (data: TemplatePreviewRequest) => templateAPI.analyze(data),
    onError: (error) => {
      const message = isAPIError(error) ? error.message : 'Failed to analyze template';
      toast({
        title: "Analysis error",
        description: message,
        variant: "destructive",
      });
    },
  });
}; 