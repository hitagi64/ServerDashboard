{
    "SystemUptime": "{{.SystemUptime}}",
    "FreeMemory": {{.FreeMemory}},
    "TotalMemory": {{.TotalMemory}},
    "UsedMemory": {{.UsedMemory}},
    "CachedMemory": {{.CachedMemory}},

    "CPUUsage": {{.CPUUsage}},
    "CPUHighestUsage": "{{.CPUHighestUsage}}",
    "CPULogHistory": [
                {{range $i, $a := .CPULogHistory}}{{if ne $i 0}} ,{{end}}[ {{range $ii, $val := $a}}{{if ne $ii 0}} ,{{end}}{{$val}}{{end}} ]
                {{end}}]
}