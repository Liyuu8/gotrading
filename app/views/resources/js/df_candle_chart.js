function send() {
  var params = {
    product_code: "BTC_JPY",
    limit: 10,
    duration: "1s",
  };
  $.get("/api/candle/", params).done(function (data) {
    console.log(data);
    var candles = data["candles"];
    $("#dashboard_div").append(candles[0].open + "<br>");
  });
}

setInterval(send, 1000 * 3);
window.onload = () => send();