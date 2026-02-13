import type {
  ProcessEvent,
  NetworkConnection,
  DiskLatency,
  CPUProfile,
  TCPLifeEvent,
  SyscallStat,
  ApiResponse,
} from '../types/types';

const BASE_URL = '/api/metrics';

async function fetchApi<T>(endpoint: string, limit: number = 50): Promise<ApiResponse<T>> {
  const res = await fetch(`${BASE_URL}${endpoint}?limit=${limit}`);
  if (!res.ok) {
    throw new Error(`API error: ${res.status} ${res.statusText}`);
  }
  return res.json();
}

export const api = {
  getProcesses: (limit = 50) => fetchApi<ProcessEvent>('/processes', limit),
  getNetwork: (limit = 50) => fetchApi<NetworkConnection>('/network', limit),
  getDisk: (limit = 20) => fetchApi<DiskLatency>('/disk', limit),
  getCPUProfile: (limit = 50) => fetchApi<CPUProfile>('/cpuprofile', limit),
  getTCPLife: (limit = 50) => fetchApi<TCPLifeEvent>('/tcplife', limit),
  getSyscalls: (limit = 50) => fetchApi<SyscallStat>('/syscalls', limit),
};

export async function checkHealth(): Promise<boolean> {
  try {
    const res = await fetch('/health');
    return res.ok;
  } catch {
    return false;
  }
}
