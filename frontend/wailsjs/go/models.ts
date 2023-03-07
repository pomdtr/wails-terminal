export namespace main {
	
	export class Theme {
	    cursor: string;
	    selectionBackground: string;
	    brightYellow: string;
	    brightWhite: string;
	    brightRed: string;
	    brightMagenta: string;
	    brightGreen: string;
	    brightCyan: string;
	    brightBlue: string;
	    brightBlack: string;
	    yellow: string;
	    white: string;
	    red: string;
	    magenta: string;
	    green: string;
	    cyan: string;
	    blue: string;
	    black: string;
	    background: string;
	    foreground: string;
	
	    static createFrom(source: any = {}) {
	        return new Theme(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cursor = source["cursor"];
	        this.selectionBackground = source["selectionBackground"];
	        this.brightYellow = source["brightYellow"];
	        this.brightWhite = source["brightWhite"];
	        this.brightRed = source["brightRed"];
	        this.brightMagenta = source["brightMagenta"];
	        this.brightGreen = source["brightGreen"];
	        this.brightCyan = source["brightCyan"];
	        this.brightBlue = source["brightBlue"];
	        this.brightBlack = source["brightBlack"];
	        this.yellow = source["yellow"];
	        this.white = source["white"];
	        this.red = source["red"];
	        this.magenta = source["magenta"];
	        this.green = source["green"];
	        this.cyan = source["cyan"];
	        this.blue = source["blue"];
	        this.black = source["black"];
	        this.background = source["background"];
	        this.foreground = source["foreground"];
	    }
	}

}

