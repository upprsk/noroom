<script lang="ts">
  import { onMount } from 'svelte';

  import * as echarts from 'echarts/core';
  import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components';
  import { PieChart } from 'echarts/charts';
  import { LabelLayout } from 'echarts/features';
  import { CanvasRenderer } from 'echarts/renderers';

  echarts.use([
    TitleComponent,
    TooltipComponent,
    LegendComponent,
    PieChart,
    CanvasRenderer,
    LabelLayout,
  ]);

  export let title: string;
  export let subtitle: string | undefined = undefined;
  export let radius = '50%';

  export let data: { value: number; name: string }[];

  let chartDiv: HTMLElement;
  let chart: echarts.ECharts;

  const updateData = (data: { value: number; name: string }[]) => {
    if (!chart) return;

    chart.setOption({
      series: [{ name: 'results', type: 'pie', data }],
    });
  };

  $: updateData(data);

  onMount(() => {
    chart = echarts.init(chartDiv);
    chart.setOption({
      title: {
        text: title,
        subtext: subtitle,
        left: 'center',
      },
      tooltip: {
        trigger: 'item',
      },
      legend: {
        orient: 'vertical',
        left: 'left',
      },
      series: [{ name: 'results', type: 'pie', radius, data }],
    });
  });
</script>

<div bind:this={chartDiv} class="aspect-video w-full"></div>
