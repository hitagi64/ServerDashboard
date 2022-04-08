const cpuhistoryctx = document.getElementById('cpuhistorychart').getContext('2d');
const cpuhistorychart = new Chart(cpuhistoryctx, {
    type: 'line',
    plugins: [
    ],
    data: cpuhistdata,
    options: {
        animation: {duration: 0},
        responsive: true,
        maintainAspectRatio: false,
        elements: {
          point:{
            radius: 0
          }
        },
        plugins: {
            legend: {
                display: false,
                position: 'left',
                labels: {
                    boxWidth: 12,
                    color: 'rgba(255, 255, 255, 1)',
                }
            },
            title: {
                display: false,
                text: ''
            }
        },
        scales: {
            x: {
              display: true,
              title: {
                display: true
              }
            },
            y: {
              display: true,
              title: {
                display: true,
                text: 'Value'
              },
              suggestedMin: 0,
              suggestedMax: 100
            }
        }
    }
});