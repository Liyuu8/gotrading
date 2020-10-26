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
  ema: {
    enable: false,
    indexs: [],
    periods: [],
    values: [],
  },
  bbands: {
    enable: false,
    indexs: [],
    n: 20,
    k: 2,
    up: [],
    mid: [],
    down: [],
  },
  ichimoku: {
    enable: false,
    indexs: [],
    tenkan: [],
    kijun: [],
    senkouA: [],
    senkouB: [],
    chikou: [],
  },
  volume: {
    enable: false,
    indexs: [],
    values: [],
  },
  rsi: {
    enable: false,
    indexs: { up: 0, value: 0, down: 0 },
    period: 14,
    up: 70,
    values: [],
    down: 30,
  },
  macd: {
    enable: false,
    indexs: [],
    periods: [],
    values: [],
  },
  hv: {
    enable: false,
    indexes: [],
    periods: [],
    values: [],
  },
};

function initConfigValues() {
  config.dataTable.index = 0;

  config.sma.indexs = [];
  config.sma.values = [];
  config.ema.indexs = [];
  config.ema.values = [];

  config.bbands.indexs = [];
  config.bbands.up = [];
  config.bbands.mid = [];
  config.bbands.down = [];

  config.ichimoku.indexs = [];
  config.ichimoku.tenkan = [];
  config.ichimoku.kijun = [];
  config.ichimoku.senkouA = [];
  config.ichimoku.senkouB = [];
  config.ichimoku.chikou = [];

  config.volume.indexs = [];

  config.rsi.indexs = [];

  config.macd.indexs = [];
  config.macd.periods = [];
  config.macd.values = [];

  config.hv.indexes = [];
  config.hv.periods = [];
  config.hv.values = [];
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
  if (config.ema.enable) {
    config.ema.indexs.forEach((emaIndex) => {
      options.series[emaIndex] = { type: 'line' };
      view.columns.push(config.candlestick.numViews + emaIndex);
    });
  }
  if (config.bbands.enable) {
    config.bbands.indexs.forEach((bbandIndex) => {
      options.series[bbandIndex] = {
        type: 'line',
        color: 'blue',
        lineWidth: 1,
      };
      view.columns.push(config.candlestick.numViews + bbandIndex);
    });
  }
  if (config.ichimoku.enable) {
    config.ichimoku.indexs.forEach((ichimokuIndex) => {
      options.series[ichimokuIndex] = {
        type: 'line',
        lineWidth: 1,
      };
      view.columns.push(config.candlestick.numViews + ichimokuIndex);
    });
  }

  if (config.volume.enable) {
    if ($('#volume_div').length === 0) {
      $('#technical_div').append(
        '<div id="volume_div" class="bottom_chart">' +
          '<span class="technical_title">Volume</span>' +
          '<div id="volume_chart"></div>' +
          '</div>'
      );
    }
    var volumeChart = new google.visualization.ChartWrapper({
      chartType: 'ColumnChart',
      containerId: 'volume_chart',
      options: {
        hAxis: { slantedText: false },
        legend: { position: 'none' },
        series: {},
      },
      view: {
        columns: [
          {
            type: 'string',
          },
          5,
        ],
      },
    });
    charts.push(volumeChart);
  }

  if (config.rsi.enable) {
    if ($('#rsi_div').length === 0) {
      $('#technical_div').append(
        '<div id="rsi_div" class="bottom_chart">' +
          '<span class="technical_title">RSI</span>' +
          '<div id="rsi_chart"></div>' +
          '</div>'
      );
    }
    var rsiChart = new google.visualization.ChartWrapper({
      chartType: 'LineChart',
      containerId: 'rsi_chart',
      options: {
        hAxis: { slantedText: false },
        legend: { position: 'none' },
        series: {
          0: { color: 'black', lineWidth: 1 },
          1: { color: '#e2431e' },
          2: { color: 'black', lineWidth: 1 },
        },
      },
      view: {
        columns: [
          {
            type: 'string',
          },
          config.candlestick.numViews + config.rsi.indexs.up,
          config.candlestick.numViews + config.rsi.indexs.value,
          config.candlestick.numViews + config.rsi.indexs.down,
        ],
      },
    });
    charts.push(rsiChart);
  }

  if (config.macd.enable) {
    if (config.macd.indexs === 0) {
      return;
    }
    if ($('#macd_div').length === 0) {
      $('#technical_div').append(
        '<div id="macd_div" class="bottom_chart">' +
          '<span class="technical_title">MACD</span>' +
          '<div id="macd_chart"></div>' +
          '</div>'
      );
    }
    var macdChart = new google.visualization.ChartWrapper({
      chartType: 'ComboChart',
      containerId: 'macd_chart',
      options: {
        legend: { position: 'none' },
        seriesType: 'bars',
        series: {
          1: { type: 'line', lineWidth: 1 },
          2: { type: 'line', lineWidth: 1 },
        },
      },
      view: {
        columns: [
          {
            type: 'string',
          },
          config.candlestick.numViews + config.macd.indexs[2],
          config.candlestick.numViews + config.macd.indexs[0],
          config.candlestick.numViews + config.macd.indexs[1],
        ],
      },
    });
    charts.push(macdChart);
  }

  if (config.hv.enable) {
    if (config.hv.indexes.length == 0) {
      return;
    }
    if ($('#hv_div').length == 0) {
      $('#technical_div').append(
        '<div id="hv_div">' +
          '<span class="technical_title">Hv</span>' +
          '<div id="hv_chart"></div>' +
          '</div>'
      );
    }

    var hvSeries = {};
    var hvColumns = [{ type: 'string' }];
    config.hv.indexes.forEach((index) => {
      hvSeries[index] = { lineWidth: 1 };
      hvColumns.push(config.candlestick.numViews + index);
    });
    var hvChart = new google.visualization.ChartWrapper({
      chartType: 'LineChart',
      containerId: 'hv_chart',
      options: {
        legend: { position: 'none' },
        series: hvSeries,
      },
      view: {
        columns: hvColumns,
      },
    });
    charts.push(hvChart);
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
  if (config.ema.enable) {
    params['ema'] = true;
    params['emaPeriod1'] = config.ema.periods[0];
    params['emaPeriod2'] = config.ema.periods[1];
    params['emaPeriod3'] = config.ema.periods[2];
  }
  if (config.bbands.enable) {
    params['bband'] = true;
    params['bbandsN'] = config.bbands.n;
    params['bbandsK'] = config.bbands.k;
  }
  if (config.ichimoku.enable) {
    params['ichimoku'] = true;
  }
  if (config.rsi.enable) {
    params['rsi'] = true;
    params['rsiPeriod'] = config.rsi.period;
  }
  if (config.macd.enable) {
    params['macd'] = true;
    params['macdPeriod1'] = config.macd.periods[0];
    params['macdPeriod2'] = config.macd.periods[1];
    params['macdPeriod3'] = config.macd.periods[2];
  }
  if (config.hv.enable) {
    params['hv'] = true;
    params['hvPeriod1'] = config.hv.periods[0];
    params['hvPeriod2'] = config.hv.periods[1];
    params['hvPeriod3'] = config.hv.periods[2];
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
    if (!!data['emas']) {
      data['emas'].forEach((emaData, index) => {
        if (emaData.length === 0) {
          return;
        }
        config.dataTable.index += 1;
        config.ema.indexs[index] = config.dataTable.index;
        dataTable.addColumn('number', 'EMA' + emaData['period'].toString());
        config.ema.values[index] = emaData['values'];
      });
    }
    if (!!data['bbands']) {
      config.dataTable.index += 1;
      config.bbands.indexs[0] = config.dataTable.index;
      config.dataTable.index += 1;
      config.bbands.indexs[1] = config.dataTable.index;
      config.dataTable.index += 1;
      config.bbands.indexs[2] = config.dataTable.index;

      var bbandsData = data['bbands'];
      dataTable.addColumn(
        'number',
        `BBands Up (${bbandsData['n']},${bbandsData['k']})`
      );
      dataTable.addColumn(
        'number',
        `BBands Mid (${bbandsData['n']},${bbandsData['k']})`
      );
      dataTable.addColumn(
        'number',
        `BBands Down (${bbandsData['n']},${bbandsData['k']})`
      );
      config.bbands.up = bbandsData['up'];
      config.bbands.mid = bbandsData['mid'];
      config.bbands.down = bbandsData['down'];
    }
    if (!!data['ichimoku']) {
      config.dataTable.index += 1;
      config.ichimoku.indexs[0] = config.dataTable.index;
      config.dataTable.index += 1;
      config.ichimoku.indexs[1] = config.dataTable.index;
      config.dataTable.index += 1;
      config.ichimoku.indexs[2] = config.dataTable.index;
      config.dataTable.index += 1;
      config.ichimoku.indexs[3] = config.dataTable.index;
      config.dataTable.index += 1;
      config.ichimoku.indexs[4] = config.dataTable.index;

      var ichimokuData = data['ichimoku'];
      config.ichimoku.tenkan = ichimokuData['tenkan'];
      config.ichimoku.kijun = ichimokuData['kijun'];
      config.ichimoku.senkouA = ichimokuData['senkoua'];
      config.ichimoku.senkouB = ichimokuData['senkoub'];
      config.ichimoku.chikou = ichimokuData['chikou'];

      dataTable.addColumn('number', 'Tenkan');
      dataTable.addColumn('number', 'Kijun');
      dataTable.addColumn('number', 'SenkouA');
      dataTable.addColumn('number', 'SenkouB');
      dataTable.addColumn('number', 'Chikou');
    }
    if (!!data['rsi']) {
      config.dataTable.index += 1;
      config.rsi.indexs.up = config.dataTable.index;
      config.dataTable.index += 1;
      config.rsi.indexs.value = config.dataTable.index;
      config.dataTable.index += 1;
      config.rsi.indexs.down = config.dataTable.index;

      var rsiData = data['rsi'];
      config.rsi.period = rsiData['period'];
      config.rsi.values = rsiData['values'];

      dataTable.addColumn('number', 'RSI Thread');
      dataTable.addColumn('number', `RSI (${config.rsi.period})`);
      dataTable.addColumn('number', 'RSI Thread');
    }
    if (!!data['macd']) {
      config.dataTable.index += 1;
      config.macd.indexs[0] = config.dataTable.index;
      config.dataTable.index += 1;
      config.macd.indexs[1] = config.dataTable.index;
      config.dataTable.index += 1;
      config.macd.indexs[2] = config.dataTable.index;

      var macdData = data['macd'];
      var fastPeriod = macdData['fast_period'].toString();
      var slowPeriod = macdData['slow_period'].toString();
      var signalPeriod = macdData['signal_period'].toString();
      var macd = macdData['macd'];
      var macdSignal = macdData['macd_signal'];
      var macdHist = macdData['macd_hist'];

      var periods = `(${fastPeriod}, ${slowPeriod}, ${signalPeriod})`;
      dataTable.addColumn('number', 'MD' + periods);
      dataTable.addColumn('number', 'MS' + periods);
      dataTable.addColumn('number', 'HT' + periods);

      config.macd.periods[0] = fastPeriod;
      config.macd.periods[1] = slowPeriod;
      config.macd.periods[2] = signalPeriod;
      config.macd.values[0] = macd;
      config.macd.values[1] = macdSignal;
      config.macd.values[2] = macdHist;
    }
    if (!!data['hvs']) {
      data['hvs'].forEach((hvData, index) => {
        if (hvData.length == 0) {
          return;
        }
        config.dataTable.index += 1;
        config.hv.indexes[index] = config.dataTable.index;

        var period = hvData['period'].toString();
        var value = hvData['values'];
        dataTable.addColumn('number', `HV(${period})`);
        config.hv.periods[index] = period;
        config.hv.values[index] = value;
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
      if (!!data['emas']) {
        config.ema.values.forEach((value) =>
          datas.push(value[index] !== 0 ? value[index] : null)
        );
      }
      if (!!data['bbands']) {
        datas.push(
          config.bbands.up[index] !== 0 ? config.bbands.up[index] : null
        );
        datas.push(
          config.bbands.mid[index] !== 0 ? config.bbands.mid[index] : null
        );
        datas.push(
          config.bbands.down[index] !== 0 ? config.bbands.down[index] : null
        );
      }
      if (!!data['ichimoku']) {
        datas.push(
          config.ichimoku.tenkan[index] !== 0
            ? config.ichimoku.tenkan[index]
            : null
        );
        datas.push(
          config.ichimoku.kijun[index] !== 0
            ? config.ichimoku.kijun[index]
            : null
        );
        datas.push(
          config.ichimoku.senkouA[index] !== 0
            ? config.ichimoku.senkouA[index]
            : null
        );
        datas.push(
          config.ichimoku.senkouB[index] !== 0
            ? config.ichimoku.senkouB[index]
            : null
        );
        datas.push(
          config.ichimoku.chikou[index] !== 0
            ? config.ichimoku.chikou[index]
            : null
        );
      }
      if (!!data['rsi']) {
        datas.push(config.rsi.up);
        datas.push(
          config.rsi.values[index] !== 0 ? config.rsi.values[index] : null
        );
        datas.push(config.rsi.down);
      }
      if (!!data['macd']) {
        config.macd.values.forEach((value) =>
          datas.push(value[index] !== 0 ? value[index] : null)
        );
      }
      if (!!data['hvs']) {
        config.hv.values.forEach((value) =>
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

  $('#inputEma').on('change', (event) => {
    config.ema.enable = event.currentTarget.checked;
    send();
  });
  $('#inputEmaPeriod1').on('change', (event) => {
    config.ema.periods[0] = event.currentTarget.value;
    send();
  });
  $('#inputEmaPeriod2').on('change', (event) => {
    config.ema.periods[1] = event.currentTarget.value;
    send();
  });
  $('#inputEmaPeriod3').on('change', (event) => {
    config.ema.periods[2] = event.currentTarget.value;
    send();
  });

  $('#inputBBands').on('change', (event) => {
    config.bbands.enable = event.currentTarget.checked;
    send();
  });
  $('#inputBBandsN').on('change', (event) => {
    config.bbands.n = event.currentTarget.value;
    send();
  });
  $('#inputBBandsK').on('change', (event) => {
    config.bbands.k = event.currentTarget.value;
    send();
  });

  $('#inputIchimoku').on('change', (event) => {
    config.ichimoku.enable = event.currentTarget.checked;
    send();
  });

  $('#inputVolume').on('change', (event) => {
    config.volume.enable = event.currentTarget.checked;
    if (!event.currentTarget.checked) {
      $('#volume_div').remove();
    }
    send();
  });

  $('#inputRsi').on('change', (event) => {
    config.rsi.enable = event.currentTarget.checked;
    if (!event.currentTarget.checked) {
      $('#rsi_div').remove();
    }
    send();
  });
  $('#inputRsiPeriod').on('change', (event) => {
    config.rsi.period = event.currentTarget.value;
    send();
  });

  $('#inputMacd').on('change', (event) => {
    config.macd.enable = event.currentTarget.checked;
    if (!event.currentTarget.checked) {
      $('#macd_div').remove();
    }
    send();
  });
  $('#inputMacdPeriod1').on('change', (event) => {
    config.macd.periods[0] = event.currentTarget.value;
    send();
  });
  $('#inputMacdPeriod2').on('change', (event) => {
    config.macd.periods[1] = event.currentTarget.value;
    send();
  });
  $('#inputMacdPeriod3').on('change', (event) => {
    config.macd.periods[2] = event.currentTarget.value;
    send();
  });

  $('#inputHv').on('change', (event) => {
    config.hv.enable = event.currentTarget.checked;
    if (!event.currentTarget.checked) {
      $('#hv_div').remove();
    }
    send();
  });
  $('#inputHvPeriod1').on('change', (event) => {
    config.hv.periods[0] = event.currentTarget.value;
    send();
  });
  $('#inputHvPeriod2').on('change', (event) => {
    config.hv.periods[1] = event.currentTarget.value;
    send();
  });
  $('#inputHvPeriod3').on('change', (event) => {
    config.hv.periods[2] = event.currentTarget.value;
    send();
  });
};
