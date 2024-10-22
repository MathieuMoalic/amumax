import { writable } from "svelte/store";

export interface Metrics {
    pid: number;
    error: string;
    cpuPercent: number;
    cpuPercentTotal: number;
    ramPercent: number;
    ramPercentTotal: number;
    
    gpuName: string;
    gpuUtilizationPercent: number;
    gpuUUID: string;
    gpuTemperature: number;
    gpuPowerDraw: number;
    gpuPowerLimit: number;
    gpuVramUsed: number;
    gpuVramTotal: number;   
}

export const metricsState = writable<Metrics>({
    pid: 0,
    error: "",
    cpuPercent: 0,
    cpuPercentTotal: 0,
    ramPercent: 0,
    ramPercentTotal: 0,
    gpuName: "",
    gpuUtilizationPercent: 0,
    gpuUUID: "",
    gpuTemperature: 0,
    gpuPowerDraw: 0,
    gpuPowerLimit: 0,
    gpuVramUsed: 0,
    gpuVramTotal: 0
});