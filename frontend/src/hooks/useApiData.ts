import { useState, useEffect, useCallback } from 'react';
import type { ApiResponse } from '../types/types';

interface UseApiDataOptions<T> {
  fetchFn: () => Promise<ApiResponse<T>>;
  pollingInterval?: number; // ms, 0 = no polling
}

interface UseApiDataResult<T> {
  data: T[];
  loading: boolean;
  error: string | null;
  refresh: () => void;
}

export function useApiData<T>({
  fetchFn,
  pollingInterval = 5000,
}: UseApiDataOptions<T>): UseApiDataResult<T> {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      const response = await fetchFn();
      setData(response.data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'API connection lost');
    } finally {
      setLoading(false);
    }
  }, [fetchFn]);

  useEffect(() => {
    fetchData();

    if (pollingInterval > 0) {
      const interval = setInterval(fetchData, pollingInterval);
      return () => clearInterval(interval);
    }
  }, [fetchData, pollingInterval]);

  return { data, loading, error, refresh: fetchData };
}
