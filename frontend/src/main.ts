import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { CanvasAddon } from "xterm-addon-canvas";
import * as App from "../wailsjs/go/main/consoleService";
import * as runtime from "../wailsjs/runtime/runtime.js";
import { Base64 } from "js-base64";

async function main() {
  const terminal = new Terminal({
    cursorBlink: true,
    allowProposedApi: true,
    allowTransparency: true,
    macOptionIsMeta: true,
    macOptionClickForcesSelection: true,
    scrollback: 0,
    fontSize: 18,
    fontFamily: "Microsoft YaHei Mono,Consolas,Liberation Mono,Menlo,Courier,monospace",
  });

  const fitAddon = new FitAddon();
  const canvasAddon = new CanvasAddon();

  terminal.open(document.getElementById("terminal")!);

  terminal.loadAddon(fitAddon);
  terminal.loadAddon(canvasAddon);

  terminal.onResize((event) => {
    var rows = event.rows;
    var cols = event.cols;
    App.Resize(rows, cols);
  });

  terminal.onData(function (data) {
    App.SendText(data);
  });

  window.onresize = () => {
    fitAddon.fit();
  };

  runtime.EventsOn("tty-data", (data: string) => {
    terminal.write(Base64.toUint8Array(data));
  });

  runtime.EventsOn("clear-terminal", () => {
    terminal.clear();
  });

  terminal.focus();
  fitAddon.fit();

  App.LoopRead();
}

main();
