function makeid(length) {
  var result = "";
  var characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  var charactersLength = characters.length;
  for (var i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
  }
  return result;
}
function make_gen_request(n) {
  for (let i = 0; i < n; ++i) {
    var count = 0;
    $.ajax({
      type: "POST",
      url: "http://localhost:8080/gen-url",
      data: {
        url: "https://www.google.com/" + makeid(10),
      },
      success: function (data, textStatus) {
        console.log(data);
      },
      fail: function (xhr, textStatus, errorThrown) {
        console.log(errorThrown);
      },
    });
    console.log(count);
  }
}
function make_get_request(n) {
  for (let i = 0; i < n; ++i) {
    var count = 0;
    $.ajax({
      type: "GET",
      url: "http://localhost:8080/" + makeid(10),
      success: function (data, textStatus) {
        console.log(data)
      },
      fail: function (xhr, textStatus, errorThrown) {
        console.log(errorThrown);
      },
    });
    console.log(count);
  }
}
function DDos(n) {
  for (let i=0; i<n; i++) {
    task();
 }
}
function task() {
  setTimeout(function() {
      make_request(1000);
  }, 1000);
}