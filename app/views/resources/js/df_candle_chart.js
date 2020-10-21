google.charts.load('current', { packages: ['corechart', 'controls'] });

var config = {
  api: {
    enable: true,
    interval: 1000 * 3,
  },
  candlestick: {
    product_code: 'BTC_JPY',
    duration: '1m',
    limit: 365,
    numViews: 5,
  },
};

function drawChart(dataTable) {
  var chartDiv = document.getElementById('chart_div');
  var dashboard = new google.visualization.Dashboard(chartDiv);
  var mainChart = new google.visualization.ChartWrapper({
    chartType: 'ComboChart',
    containerId: 'chart_div',
    options: {
      hAxis: { slantedText: false },
      legend: { position: 'none' },
      candlestick: {
        fallingColor: { strokeWidth: 0, fill: '#a52714' },
        rasingColor: { strokeWidth: 0, fill: '#0f9d58' },
      },
      seriesType: 'candlesticks',
      series: {},
    },
    view: {
      columns: [
        {
          cals: (d, rowIndex) => d.getFormattedValue(rowIndex, 0),
          type: 'string',
        },
        1,
        2,
        3,
        4,
      ],
    },
  });
  var charts = [mainChart];

  var options = mainChart.getOptions();
  var view = mainChart.getView();

  var controlWrapper = new google.visualization.ControlWrapper({
    controlType: 'ChartRangeFilter',
    containerId: 'filter_div',
    options: {
      filterColumnIndex: 0,
      ui: {
        chartType: 'LineChart',
        chartView: {
          columns: [0, 4],
        },
      },
    },
  });

  dashboard.bind(controlWrapper, charts);
  dashboard.draw(dataTable);
}

function send() {
  if (!config.api.enable) {
    return;
  }

  var params = {
    product_code: config.candlestick.product_code,
    limit: config.candlestick.limit,
    duration: config.candlestick.duration,
  };
  $.get('/api/candle/', params).done(function (data) {
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn('date', 'Date');
    dataTable.addColumn('number', 'Low');
    dataTable.addColumn('number', 'Open');
    dataTable.addColumn('number', 'Close');
    dataTable.addColumn('number', 'High');
    dataTable.addColumn('number', 'Volume');

    var googleChartData = data['candles'].map((candle) => [
      new Date(candle.time),
      candle.low,
      candle.open,
      candle.close,
      candle.high,
      candle.volume,
    ]);
    dataTable.addRows(googleChartData);
    drawChart(dataTable);
  });
}

function changeDuration(s) {
  config.candlestick.duration = s;
  send();
}

setInterval(send, 1000 * 3);
window.onload = () => {
  send();

  $('#dashboard_div')
    .mouseenter(() => (config.api.enable = false))
    .mouseleave(() => (config.api.enable = true));
};
