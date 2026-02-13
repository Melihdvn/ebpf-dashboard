export interface ProcessEvent {
  id: number;
  timestamp: string;
  time: string;
  pid: string;
  comm: string;
  args: string;
}

export interface NetworkConnection {
  id: number;
  timestamp: string;
  pid: string;
  comm: string;
  ip_version: string;
  source_addr: string;
  source_port: string;
  dest_addr: string;
  dest_port: string;
}

export interface DiskLatency {
  id: number;
  timestamp: string;
  range_min: number;
  range_max: number;
  count: number;
}

export interface CPUProfile {
  id: number;
  timestamp: string;
  process_name: string;
  stack_trace: string;
  sample_count: number;
}

export interface TCPLifeEvent {
  id: number;
  timestamp: string;
  pid: number;
  comm: string;
  local_addr: string;
  local_port: number;
  remote_addr: string;
  remote_port: number;
  tx_kb: number;
  rx_kb: number;
  duration_ms: number;
}

export interface SyscallStat {
  id: number;
  timestamp: string;
  syscall_name: string;
  count: number;
}

export interface ApiResponse<T> {
  count: number;
  data: T[];
}
