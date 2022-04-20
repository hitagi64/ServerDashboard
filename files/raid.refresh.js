var intervalId = window.setInterval(function(){
    fetch('/dashboardRefreshData?req=md').then(function (response) {
        return response.text();
    }).then(function (text) {
        return JSON.parse(text);
    }).then(function (json) {
        for (let i = 0; i < json.MDInfo.length; i++) {
            let name = json.MDInfo[i].Name;

            if (document.getElementById(name + "_title") == null) {
                document.getElementById("dpageupdatemessage").style.visibility = "visible";
                continue;
            }

            document.getElementById(name + "_state").innerHTML = json.MDInfo[i].Array_state;
            document.getElementById(name + "_size").innerHTML = json.MDInfo[i].SizeShortened;
            document.getElementById(name + "_level").innerHTML = json.MDInfo[i].Level;
            document.getElementById(name + "_diskcount").innerHTML = json.MDInfo[i].Raid_disks;
            document.getElementById(name + "_degraded").innerHTML = json.MDInfo[i].Degraded;
            document.getElementById(name + "_syncaction").innerHTML = json.MDInfo[i].Sync_action;
            document.getElementById(name + "_uuid").innerHTML = json.MDInfo[i].UUID;
            if (json.MDInfo[i].ArrayIsDegraded) {
                document.getElementById(name + "_stateindicator").innerHTML = "<i class=\"fa-solid fa-circle-exclamation dorangeindicator\"></i>";
            }else if (json.MDInfo[i].ArrayIsGood) {
                document.getElementById(name + "_stateindicator").innerHTML = "<i class=\"fa-solid fa-circle-check dgreenindicator\"></i>";
            }else {
                document.getElementById(name + "_stateindicator").innerHTML = "<i class=\"fa-solid fa-circle-xmark dredindicator\"></i>";
            }
            if (json.MDInfo[i].Degraded == 0) {
                document.getElementById(name + "_degradedindicator").innerHTML = "<i class=\"fa-solid fa-circle-check dgreenindicator\"></i>";
            }else {
                document.getElementById(name + "_degradedindicator").innerHTML = "<i class=\"fa-solid fa-circle-xmark dredindicator\"></i>";
            }
        }
    }).catch(function (err) {
        console.warn('Something went wrong.', err);
    });
}, 1000);