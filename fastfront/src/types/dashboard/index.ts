export type TSystemInfo = {
  cpu: TCpu;
  memory: TMemory;
  disk: TDisk;
  os: TOs;
  appVersion: string;
};

export type TCpu = {
  count: number;
  cpu_percent: number[];
  states: TCpuState[];
};
export type TMemory = {
  available: number;
  total: number;
  used: number;
  used_percent?: number;
};
export type TOs = {
  compiler: string;
  go_version: string;
  num_cpu: number;
  num_goroutine: number;
  os: string;
};
export type TDisk = {
  free: number;
  partition: string;
  total: number;
  used: number;
  used_percent?: number;
};
export type TCpuState = {
  cpu: number;
  vendorId: string;
  family: string;
  model: string;
  stepping: number;
  physicalId: string;
  coreId: string;
  cores: number;
  modelName: string;
  mhz: number;
  cacheSize: number;
  flags: string[];
  microcode: string;
};
