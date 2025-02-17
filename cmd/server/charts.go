package main

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/channingko-madden/pi-vitrine/internal"
	"github.com/channingko-madden/pi-vitrine/internal/cher"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
)

// Copied from go-echarts/templates/base.tpl
var baseTpl = `
<div class="container">
    <div class="item" id="{{ .ChartID }}" style="width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>

{{- range .JSAssets.Values }}
    <script src="{{ . }}"></script>
{{- end }}

<script type="text/javascript">
    "use strict";
    let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}", { renderer: "{{  .Initialization.Renderer }}" });
    let option_{{ .ChartID | safeJS }} = {{- .JSONNotEscaped | safeJS }}
    goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});

  {{- range  $listener := .EventListeners }}
    {{if .Query  }}
    goecharts_{{ $.ChartID | safeJS }}.on({{ $listener.EventName }}, {{ $listener.Query | safeJS }}, {{ injectInstance $listener.Handler "%MY_ECHARTS%"  $.ChartID | safeJS }});
    {{ else }}
    goecharts_{{ $.ChartID | safeJS }}.on({{ $listener.EventName }}, {{ injectInstance $listener.Handler "%MY_ECHARTS%"  $.ChartID | safeJS }})
    {{ end }}
  {{- end }}

    {{- range .JSFunctions.Fns }}
    {{ injectInstance . "%MY_ECHARTS%"  $.ChartID  | safeJS }}
    {{- end }}
</script>
`

type chartRender struct {
	render.BaseRender
	c      interface{}
	before []func()
}

func newChartRender(c interface{}, before ...func()) render.Renderer {
	return &chartRender{c: c, before: before}
}

func (c *chartRender) Render(w io.Writer) error {
	for _, fn := range c.before {
		fn()
	}

	contents := []string{baseTpl}
	tpl := render.MustTemplate("chart", contents)
	return tpl.ExecuteTemplate(w, "chart", c.c)
}

func chartIndoorClimate(w http.ResponseWriter, data []cher.IndoorClimate) error {
	// time
	xdata := make([]string, len(data))
	// temperature, RH, Pressure
	tempData := make([]opts.LineData, len(data))
	pressureData := make([]opts.LineData, len(data))
	rhData := make([]opts.LineData, len(data))
	for i, data := range data {
		xdata[i] = data.CreatedAt.Format(time.DateTime)
		tempData[i] = opts.LineData{Value: internal.KelvinToCelcius(data.AirTemp)}
		pressureData[i] = opts.LineData{Value: data.Pressure}
		rhData[i] = opts.LineData{Value: data.RelativeHumidity}
	}

	var buf bytes.Buffer

	tempLine := charts.NewLine()
	tempLine.Renderer = newChartRender(tempLine, tempLine.Validate)
	tempLine.SetXAxis(xdata).AddSeries("Temperature (C)", tempData)

	err := tempLine.Render(&buf)
	if err != nil {
		return err
	}

	pressureLine := charts.NewLine()
	pressureLine.Renderer = newChartRender(pressureLine, pressureLine.Validate)
	pressureLine.SetXAxis(xdata).AddSeries("Pressure (bar)", pressureData)

	err = pressureLine.Render(&buf)
	if err != nil {
		return err
	}

	rhLine := charts.NewLine()
	rhLine.Renderer = newChartRender(rhLine, rhLine.Validate)
	rhLine.SetXAxis(xdata).AddSeries("Relative Humidity (%)", rhData)

	err = rhLine.Render(&buf)
	if err != nil {
		return err
	}

	w.Write(buf.Bytes())

	return nil
}

func chartSystemData(w http.ResponseWriter, data []cher.System) error {
	// time
	xdata := make([]string, len(data))
	// system
	cpuTempData := make([]opts.LineData, len(data))
	gpuTempData := make([]opts.LineData, len(data))
	for i, data := range data {
		xdata[i] = data.CreatedAt.Format(time.DateTime)
		cpuTempData[i] = opts.LineData{Value: internal.KelvinToCelcius(data.CPUTemp)}
		gpuTempData[i] = opts.LineData{Value: internal.KelvinToCelcius(data.GPUTemp)}
	}

	var buf bytes.Buffer
	cpuTempLine := charts.NewLine()
	cpuTempLine.Renderer = newChartRender(cpuTempLine, cpuTempLine.Validate)
	cpuTempLine.SetXAxis(xdata).AddSeries("CPU Temperature (C)", cpuTempData)

	err := cpuTempLine.Render(&buf)
	if err != nil {
		return err
	}

	gpuTempLine := charts.NewLine()
	gpuTempLine.Renderer = newChartRender(gpuTempLine, gpuTempLine.Validate)
	gpuTempLine.SetXAxis(xdata).AddSeries("GPU Temperature (C)", gpuTempData)

	err = gpuTempLine.Render(&buf)
	if err != nil {
		return err
	}

	w.Write(buf.Bytes())

	return nil
}
