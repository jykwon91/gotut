package main

import (
	"fmt"
	"log"
	"io"
	"html/template"
	"net/http"
	"golang.org/x/net/websocket"
)

const listenAddr = "localhost:4000"

type socket struct {
	io.ReadWriter
	done chan bool
	name string
}

func (s socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(ws *websocket.Conn) {
	s := socket{ws, make(chan bool), ""}
	go match(s)
	<-s.done
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.Handle("/socket", websocket.Handler(socketHandler))
	err := http.ListenAndServe("localhost:4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type PageData struct {
	ListenAddr string
	Name string
}

var partner = make(chan io.ReadWriteCloser)
func match(s socket) {
	fmt.Fprint(s, "Waiting for a partner...")
	select {
	case partner <- s:
	// now handled by the other goroutine
	case p := <-partner:
		chat(p, s)
	}
}
func chat(a, b io.ReadWriteCloser) {
	fmt.Fprintln(a, "Found one! Say Hi.")
	fmt.Fprintln(b, "Found one! Say Hi.")
	errc := make(chan error, 1)
	go cp(a, b, errc, "Jason")
	go cp(b, a, errc, "Kwon")
	if err := <-errc; err != nil {
		fmt.Println(err)
	}
	a.Close()
	b.Close()
}
func cp(w io.Writer, r io.Reader, errc chan<- error, name string) {
	_, err := io.Copy(w, r)
	if err != nil {
		fmt.Println(err.Error())
	}
	errc <- err
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		ListenAddr: listenAddr,
		Name: "Received",
	}
	err := rootTemplate.Execute(w, data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

var rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8" />
</head>
<body>
<div>
    <input id="message" type="text" placeholder="message" />
    <button id="send">Send</button>
</div>
<div id="output">
</div>
</body>
<script>
	const output = document.getElementById("output");
	const button = document.querySelector("#send");
	
	websocket = new WebSocket("ws://{{.ListenAddr}}/socket");
	
	websocket.onmessage = function (event) {
		writeToScreen('<span style="color: blue;">{{.Name}}: ' + event.data+'</span>');
	};
  	
	websocket.onclose = function onClose(evt) {
    		writeToScreen("DISCONNECTED");
  	}
	
	function writeToScreen(message) {
		var pre = document.createElement("p");
		pre.style.wordWrap = "break-word";
		pre.innerHTML = message;
		output.appendChild(pre);
	}

	button.addEventListener("click", () => {
		const message = document.querySelector("#message");

		// Send composed message to the server
		websocket.send(message.value);
		writeToScreen(message.value);

		// clear input fields
		message.value = "";
	});
</script>
</html>
`))
