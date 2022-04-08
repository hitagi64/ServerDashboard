<link rel="stylesheet" href="files/dashboard.css"/>
<script type="text/javascript" src="files/chart.min.js"></script>

<div class="container" style="max-width: 100%">
    <div class="row">
        <div class="col-xl-4 col-lg-6 col-md-12" style="padding: 0px">
            <div class="dcard">
                <div class="dcardtitle">System information</div>
                <div class="dcardcontent">
                    <span class="dinfoline"><b>System Name: </b> {{.SystemName}}</span>
                    <span class="dinfoline"><b>Kernel: </b> {{.Kernel}} {{.KernelVersion}}</span>
                    <span class="dinfoline"><b>Distro: </b> {{.DistroName}}</span>
                    <span class="dinfoline"><b>Architecture: </b> {{.SystemArchitecture}}</span>
                    <span class="dinfoline"><b>Uptime: </b> <span id="systemUptimeID">{{.SystemUptime}}</span></span>
                </div>
            </div>
        </div>

        <div class="col-xl-4 col-lg-6 col-md-12" style="padding: 0px">
            <div class="dcard">
                <div class="dcardtitle">Memory</div>
                <div class="dcardcontent">
                    <div style="width: 100%; height: 100%">
                        <canvas id="memchart"></canvas>
                        <script>
                            var freeMemory = {{.FreeMemory}};
                            var usedMemory = {{.UsedMemory}};
                            var cachedMemory = {{.CachedMemory}};
                            var totalMemory = {{.TotalMemory}};
                        </script>
                        <script type="text/javascript" src="files/mem.chart.js"></script>
                    </div>
                </div>
            </div>
        </div>

        <div class="col-xl-4 col-lg-6 col-md-12" style="padding: 0px">
            <div class="dcard">
                <div class="dcardtitle">CPU</div>
                <div class="dcardcontent">
                    <div class="dcpuname">
                        <span>{{.CPUName}}</span>
                    </div>
                    <div class="dcpuflex">
                        <div style="width: 40%; height: 100%">
                            <canvas id="cpuchart"></canvas>
                            <script>
                                var cpuUsage = {{.CPUUsage}};
                            </script>
                            <script type="text/javascript" src="files/cpu.chart.js"></script>
                        </div>
                        <div style="flex-grow: 2">
                            <span class="dinfoline"><b>Cores: </b> {{.CPUCoreCount}}</span>
                            <span class="dinfoline"><b>Threads: </b> {{.CPUThreadCount}}</span>
                            <span class="dinfoline"><b>Highest usage: </b> <span id="cpuHighestUsageID">{{.CPUHighestUsage}}</span></span>
                            <span class="dinfoline"><b>Virtualization: </b> {{if .CPUHasVirtualization}}Enabled{{else}}Disabled{{end}}</span>
                        </div>
                    </div>
                    <div style="height: 45%; width: 100%">
                        <canvas id="cpuhistorychart"></canvas>
                        <script>
                            const labels = [];
                            for (let i = {{.CPUMaxHistoryLength}}-1; i >= 0; i--) {
                                labels.push(i.toString());
                            }
                            function getRandomInt(max) {
                                return Math.floor(Math.random() * max);
                            }
                            function getNewColor() {
                                return "rgb(" + 
                                    getRandomInt(255) + ", " + 
                                    getRandomInt(255) + ", " + 
                                    getRandomInt(255) + ")";
                            }
                            const cpuhistdata = {
                                labels: labels,
                                datasets: [
                                    {{range $i, $a := .CPULogHistory}}
                                    {
                                        label: 'Core {{$i}}',
                                        data: [
                                            {{range $ii, $val := $a}}
                                                {{$val}},
                                            {{end}}
                                        ],
                                        borderColor: getNewColor(),
                                        fill: false,
                                        cubicInterpolationMode: 'monotone',
                                        tension: 0.4
                                    },
                                    {{end}}
                                ]
                            }
                        </script>
                        <script type="text/javascript" src="files/cpuhistory.chart.js"></script>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<script type="text/javascript" src="files/mem.cpu.refresh.js"></script>