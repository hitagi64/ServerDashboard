// Update as much as possible every second
var intervalId = window.setInterval(function(){
    fetch('/dashboardRefreshData').then(function (response) {
        return response.text();
    }).then(function (text) {
        return JSON.parse(text);
    }).then(function (json) {
        document.getElementById("cpuHighestUsageID").innerHTML = json.CPUHighestUsage;

        document.getElementById("systemUptimeID").innerHTML = json.SystemUptime;

        memchart.data.datasets[0].data = [json.FreeMemory, json.CachedMemory, json.UsedMemory];
        memchart.data.labels[0] = (json.FreeMemory/1024.0).toFixed(1) + 'GB Free';
        memchart.data.labels[1] = (json.CachedMemory/1024.0).toFixed(1) + 'GB Cache/Free';
        memchart.data.labels[2] = (json.UsedMemory/1024.0).toFixed(1) + 'GB Used';

        cpuchart.data.datasets[0].data = [json.CPUUsage, 100-json.CPUUsage];
        cpuchart.options.elements.center.text = json.CPUUsage.toFixed(0) + "%";
        
        for (let i = 0; i < json.CPULogHistory.length; i++) {
            if (i < cpuhistorychart.data.datasets.length) {
                cpuhistorychart.data.datasets[i].data = json.CPULogHistory[i];
            }else{
                // Didn't test this really, but should work
                // I doubt this would ever be used anyways
                // Unless you somehow can add a core to your system while its running
                console.warn("Core count in graph does not match core count in data");
                console.warn("Updating graph to match data");
                cpuhistorychart.data.datasets.push({
                    label: 'Core ' + i,
                    data: json.CPULogHistory[i],
                    borderColor: getNewColor(),
                    borderWidth: 1
                });
            }
        }
        
        memchart.update();
        cpuchart.update();
        cpuhistorychart.update();
    }).catch(function (err) {
        console.warn('Something went wrong.', err);
    });
}, 1000);