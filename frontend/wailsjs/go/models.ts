export namespace main {
	
	export class Config {
	    excelPath: string;
	    imageDir: string;
	    codeCol: string;
	    imageCol: string;
	    sheetName: string;
	    rowHeight: number;
	    colWidth: number;
	    workerCount: number;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.excelPath = source["excelPath"];
	        this.imageDir = source["imageDir"];
	        this.codeCol = source["codeCol"];
	        this.imageCol = source["imageCol"];
	        this.sheetName = source["sheetName"];
	        this.rowHeight = source["rowHeight"];
	        this.colWidth = source["colWidth"];
	        this.workerCount = source["workerCount"];
	    }
	}
	export class ProcessResult {
	    success: boolean;
	    message: string;
	    missingCodes: string[];
	    outputPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.missingCodes = source["missingCodes"];
	        this.outputPath = source["outputPath"];
	    }
	}

}

