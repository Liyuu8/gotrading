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
  dataTable: {
    index: 0,
    value: null,
  },
  sma: {
    enable: false,
    indexs: [],
    periods: [],
    values: [],
  },
};

function initConfigValues() {
  config.dataTable.index = 0;
  config.sma.indexs = [];
  config.sma.values = [];
}

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

  if (config.sma.enable) {
    config.sma.indexs.forEach((smaIndex) => {
      options.series[smaIndex] = { type: 'line' };
      view.columns.push(config.candlestick.numViews + smaIndex);
    });
  }

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
  if (config.sma.enable) {
    params['sma'] = true;
    params['smaPeriod1'] = config.sma.periods[0];
    params['smaPeriod2'] = config.sma.periods[1];
    params['smaPeriod3'] = config.sma.periods[2];
  }
  $.get('/api/candle/', params).done(function (data) {
    initConfigValues();

    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn('date', 'Date');
    dataTable.addColumn('number', 'Low');
    dataTable.addColumn('number', 'Open');
    dataTable.addColumn('number', 'Close');
    dataTable.addColumn('number', 'High');
    dataTable.addColumn('number', 'Volume');

    if (!!data['smas']) {
      data['smas'].forEach((smaData, index) => {
        if (smaData.length === 0) {
          return;
        }
        config.dataTable.index += 1;
        config.sma.indexs[index] = config.dataTable.index;
        dataTable.addColumn('number', 'SMA' + smaData['period'].toString());
        config.sma.values[index] = smaData['values'];
      });
    }

    var googleChartData = data['candles'].map((candle, index) => {
      var datas = [
        new Date(candle.time),
        candle.low,
        candle.open,
        candle.close,
        candle.high,
        candle.volume,
      ];

      if (!!data['smas']) {
        config.sma.values.forEach((value) =>
          datas.push(value[index] !== 0 ? value[index] : null)
        );
      }
      return datas;
    });
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

  $('#inputSma').on('change', (event) => {
    config.sma.enable = event.currentTarget.checked;
    send();
  });
  $('#inputSmaPeriod1').on('change', (event) => {
    config.sma.periods[0] = event.currentTarget.value;
    send();
  });
  $('#inputSmaPeriod2').on('change', (event) => {
    config.sma.periods[1] = event.currentTarget.value;
    send();
  });
  $('#inputSmaPeriod3').on('change', (event) => {
    config.sma.periods[2] = event.currentTarget.value;
    send();
  });
};
