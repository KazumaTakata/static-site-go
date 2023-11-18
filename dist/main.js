
window.onload = function () {
    var conn;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
        };
        conn.onmessage = function (evt) {
            var messages = evt.data.split('\n');
            console.log(messages)
            location.reload()
        };
    } else {
    }
};

