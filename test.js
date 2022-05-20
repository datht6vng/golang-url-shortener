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
for (let i = 0; i < 1000; ++i) {
  var count = 0;
  $.ajax({
    type: "POST",
    url: "http://localhost:8080/gen-url",
    data: {
      url: "https://www.google.com/" + makeid(10),
    },
    success: function (data, textStatus) {},
    fail: function (xhr, textStatus, errorThrown) {
      ++count;
    },
  });
  console.log(count);
}
