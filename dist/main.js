
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
    }


    let currentColorMode = "dark"

    let toggle = document.getElementsByClassName("theme-toggler")
    if (toggle.length > 0) {
        toggle[0].addEventListener("click", () => {
            if (currentColorMode == "dark") {
                if (document.body.classList.contains("dark-mode")) {
                    document.body.classList.remove("dark-mode")
                }

                document.body.classList.add("light-mode")
                currentColorMode = "light"
                return
            }

            if (currentColorMode == "light") {
                if (document.body.classList.contains("light-mode")) {
                    document.body.classList.remove("light-mode")
                }

                document.body.classList.add("dark-mode")
                currentColorMode = "dark"
                return
            }

        })
    }

    
};

