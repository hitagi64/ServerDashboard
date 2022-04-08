const memctx = document.getElementById('memchart').getContext('2d');
const memchart = new Chart(memctx, {
    type: 'doughnut',
    plugins: [{
        afterDraw: chart => {
            let ctx = memctx;
            ctx.save();
            ctx.textAlign = 'center';
            ctx.font = '26px Arial';
            ctx.fillStyle = 'white';
            ctx.textAlign = 'left';
            ctx.fillText((totalMemory/1024.0).toFixed(1) + " GB", 10, 22);
            ctx.font = '12px Arial';
            ctx.fillText("Total system memory", 10, 40);
            ctx.restore();
        }
    }],
    data: {
        labels: [(freeMemory/1024.0).toFixed(1) + 'GB Free', 
                (cachedMemory/1024.0).toFixed(1) + 'GB Cache/Free', 
                (usedMemory/1024.0).toFixed(1) + 'GB Used'],
        datasets: [
            {
                label: 'Usage',
                data: [freeMemory, cachedMemory, usedMemory],
                backgroundColor: [
                    'rgba(54, 162, 235, 0.5)',
                    'rgba(79, 135, 180, 0.5)',
                    'rgba(255, 90, 86, 0.5)'
                ],
                borderColor: [
                    'rgba(54, 162, 235, 1)',
                    'rgba(79, 135, 180, 1)',
                    'rgba(255, 90, 86, 1)'
                ],
                borderWidth: 1
            }
        ]
    },
    options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'left',
                align: 'end',
                labels: {
                    boxWidth: 12,
                    color: 'rgba(255, 255, 255, 1)',
                }
            },
            title: {
                display: false,
                text: ''
            }
        }
    }
});