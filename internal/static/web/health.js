// health.js
let startTime;
let realtimeChart;

document.addEventListener('DOMContentLoaded', () => {
    const serverData = document.getElementById('dynamic');
    const diskUse = parseFloat(serverData.dataset.diskUse);
    const cpuUse = parseFloat(serverData.dataset.cpuUse);
    const memoryUse = parseFloat(serverData.dataset.memoryUse);
    startTime = new Date(serverData.dataset.startTime);

    createDiskChart(diskUse);
    createRealtimeChart();
    initWebSocket();
    updateUptime();
    updateStatusValues({ CPUUse: cpuUse, MemoryUse: memoryUse });
});

function createDiskChart(diskUse) {
    const canvas = document.getElementById('diskChart');
    if (!canvas) {
        console.error('Canvas element not found');
        return;
    }
    const ctx = canvas.getContext('2d');
    const color = getComputedStyle(document.documentElement).getPropertyValue(diskUse > 90 ? '--lv4' : diskUse > 70 ? '--lv3' : diskUse > 50 ? '--lv2' : '--lv1').trim();
    const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-color').trim();

    new Chart(ctx, {
        type: 'doughnut',
        data: {
            datasets: [{
                data: [diskUse, 100 - diskUse],
                backgroundColor: [color, 'rgba(255, 255, 255, 0.1)'],
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            cutout: '70%',
            plugins: {
                legend: { display: false },
                tooltip: {
                    callbacks: {
                        label: function(context) {
                            return `Disk Usage: ${context.parsed}%`;
                        }
                    },
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    titleColor: textColor,
                    bodyColor: textColor
                }
            }
        }
    });
}

function createRealtimeChart() {
    const ctx = document.getElementById('realtimeChart').getContext('2d');
    const cpuColor = getComputedStyle(document.documentElement).getPropertyValue('--lv4').trim();
    const memoryColor = getComputedStyle(document.documentElement).getPropertyValue('--lv1').trim();
    const textColor = getComputedStyle(document.documentElement).getPropertyValue('--text-color').trim();
    const gridColor = getComputedStyle(document.documentElement).getPropertyValue('--border-color').trim();

    realtimeChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [
                {
                    label: 'CPU',
                    data: [],
                    borderColor: cpuColor,
                    backgroundColor: `${cpuColor}33`,
                    borderWidth: 2,
                    tension: 0.4,
                    fill: true
                },
                {
                    label: 'Memory',
                    data: [],
                    borderColor: memoryColor,
                    backgroundColor: `${memoryColor}33`,
                    borderWidth: 2,
                    tension: 0.4,
                    fill: true
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true,
                    max: 100,
                    ticks: {
                        color: textColor,
                        font: {
                            weight: 'bold'
                        }
                    },
                    grid: {
                        color: gridColor,
                        lineWidth: 0.5
                    }
                },
                x: {
                    ticks: {
                        color: textColor,
                        font: {
                            weight: 'bold'
                        }
                    },
                    grid: {
                        color: gridColor,
                        lineWidth: 0.5
                    }
                }
            },
            plugins: {
                legend: {
                    labels: {
                        color: textColor,
                        font: {
                            weight: 'bold'
                        }
                    }
                },
                tooltip: {
                    mode: 'index',
                    intersect: false,
                    backgroundColor: 'rgba(0, 0, 0, 0.8)',
                    titleColor: textColor,
                    bodyColor: textColor
                }
            },
            interaction: {
                mode: 'nearest',
                axis: 'x',
                intersect: false
            }
        }
    });
}

function initWebSocket() {
    const socket = new WebSocket('ws://' + window.location.host + '/ws');
    socket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        updateChartData(data);
        updateStatusValues(data);
    };
}

function updateChartData(data) {
    const now = new Date();
    realtimeChart.data.labels.push(now.toLocaleTimeString());
    realtimeChart.data.datasets[0].data.push(data.CPUUse);
    realtimeChart.data.datasets[1].data.push(data.MemoryUse);

    if (realtimeChart.data.labels.length > 20) {
        realtimeChart.data.labels.shift();
        realtimeChart.data.datasets[0].data.shift();
        realtimeChart.data.datasets[1].data.shift();
    }

    realtimeChart.update();
}

function updateStatusValues(data) {
    document.getElementById('cpuValue').textContent = `${data.CPUUse.toFixed(1)}%`;
    document.getElementById('memoryValue').textContent = `${data.MemoryUse.toFixed(1)}%`;
}

function updateUptime() {
    const uptimeElement = document.getElementById('uptime');
    setInterval(() => {
        const now = new Date();
        const diff = now - startTime;
        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);
        uptimeElement.textContent = `${days}d ${hours}h ${minutes}m ${seconds}s`;
    }, 1000);
}

// 既存のコードの最後に以下を追加
window.addEventListener('resize', () => {
    if (realtimeChart) {
        realtimeChart.resize();
    }
});