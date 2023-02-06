import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { CanvasAddon } from "xterm-addon-canvas";
import * as App from "../wailsjs/go/main/App.js";
import * as runtime from "../wailsjs/runtime/runtime.js";
import { Base64 } from "js-base64";

const terminal = new Terminal({
  cursorBlink: true,
  allowProposedApi: true,
  allowTransparency: true,
  macOptionIsMeta: true,
  macOptionClickForcesSelection: true,
  scrollback: 0,
  fontSize: 13,
  fontFamily: "Consolas,Liberation Mono,Menlo,Courier,monospace",
});

const fitAddon = new FitAddon();
const canvasAddon = new CanvasAddon();

terminal.open(document.getElementById("terminal")!);

terminal.loadAddon(fitAddon);
terminal.loadAddon(canvasAddon);

terminal.focus();

terminal.onResize((event) => {
  var rows = event.rows;
  var cols = event.cols;
  App.SetTTYSize(rows, cols);
});

terminal.onData(function (data) {
  App.SendText(data);
});

window.onresize = () => {
  fitAddon.fit();
};

runtime.EventsOn("ttyData", (data: string) => {
  terminal.write(Base64.toUint8Array(data));
});

runtime.EventsOn("clear-terminal", () => {
  terminal.clear();
});

window.onblur = () => {
  App.HideWindow();
};

fitAddon.fit();
