#!/usr/bin/env node

const degit = require("degit");
const path = require("path");
const fs = require("fs/promises");
const url = require("url");
const os = require("os");
const _ = require("lodash");

const replacementMap = {
  "terminal.foreground": "foreground",
  "terminal.background": "background",
  "terminal.ansiBlack": "black",
  "terminal.ansiBlue": "blue",
  "terminal.ansiCyan": "cyan",
  "terminal.ansiGreen": "green",
  "terminal.ansiMagenta": "magenta",
  "terminal.ansiRed": "red",
  "terminal.ansiWhite": "white",
  "terminal.ansiYellow": "yellow",
  "terminal.ansiBrightBlack": "brightBlack",
  "terminal.ansiBrightBlue": "brightBlue",
  "terminal.ansiBrightCyan": "brightCyan",
  "terminal.ansiBrightGreen": "brightGreen",
  "terminal.ansiBrightMagenta": "brightMagenta",
  "terminal.ansiBrightRed": "brightRed",
  "terminal.ansiBrightWhite": "brightWhite",
  "terminal.ansiBrightYellow": "brightYellow",
  "terminal.selectionBackground": "selectionBackground",
  "terminalCursor.foreground": "cursor",
};

async function main() {
  const vscodeDir = await fs.mkdtemp(path.join(os.tmpdir(), "vscode-themes"));
  console.log("Downloading themes...");
  const downloader = degit("mbadolato/iTerm2-Color-Schemes/vscode");
  await downloader.clone(vscodeDir);

  const targetDir = path.join(__dirname, "..", "themes");
  await fs.rm(targetDir, { recursive: true, force: true });
  await fs.mkdir(targetDir, { recursive: true });

  const entries = await fs.readdir(vscodeDir);
  console.log(`Converting ${entries.length} themes...`);
  const promises = entries.map(async (theme) => {
    const filepath = path.join(vscodeDir, theme);
    const content = await fs.readFile(filepath, "utf-8");
    const vscodeTheme = JSON.parse(content);
    const xtermTheme = {};

    for (const [key, color] of Object.entries(
      vscodeTheme["workbench.colorCustomizations"]
    )) {
      xtermTheme[replacementMap[key]] = color;
    }

    const { name, ext } = path.parse(theme);
    const clean = _.kebabCase(name);
    await fs.writeFile(
      path.join(targetDir, `${clean}${ext}`),
      JSON.stringify(xtermTheme, null, 2)
    );
  });

  await Promise.all(promises);
  await fs.rm(vscodeDir, { recursive: true, force: true });
  console.log("Done!");
}

main();
