import { tablePlotState } from '$api/incoming/table-plot';
import { Chart } from 'chart.js';
import { get, writable } from 'svelte/store';

export let chart = writable<Chart | null>(null);

export function createChart() {
    chart.set(new Chart(document.getElementById('gd') as HTMLCanvasElement, {
        type: 'scatter',
        data: {
            datasets: [
                {
                    data: get(tablePlotState).data,
                    backgroundColor: 'rgba(240, 192, 192, 1)',
                    borderColor: 'rgba(75, 192, 192, 1)',
                    borderWidth: 2,
                    pointRadius: 0,
                    showLine: true
                }
            ]
        },
        options: {
            responsive: true,
            scales: {
                x: {
                    beginAtZero: false,
                    grid: {
                        color: '#f8f8f2'
                    },
                    ticks: {
                        color: '#f8f8f2'
                    }
                },
                y: {
                    beginAtZero: false,
                    grid: {
                        color: '#f8f8f2'
                    },
                    ticks: {
                        color: '#f8f8f2'
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                },
                zoom: {
                    pan: {
                        enabled: true
                    },
                    zoom: {
                        wheel: {
                            enabled: true,
                            speed: 0.1
                        }
                    }
                }
            }
        }
    }));
}

export function plotChart() {
    const chartInstance = get(chart);
    if (chartInstance) {
        chartInstance.data.datasets[0].data = get(tablePlotState).data;
        chartInstance.update();
    }
}
export function resetZoom() {
    get(chart)?.resetZoom('active');
    get(chart)?.update();
}