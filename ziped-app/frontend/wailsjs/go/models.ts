export namespace models {
	
	export class ArchiveResult {
	    success: boolean;
	    error: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ArchiveResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.error = source["error"];
	        this.message = source["message"];
	    }
	}
	export class FileInfo {
	    name: string;
	    path: string;
	    originalSize: number;
	    compressedSize: number;
	    type: string;
	    modified: string;
	    isDir: boolean;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.originalSize = source["originalSize"];
	        this.compressedSize = source["compressedSize"];
	        this.type = source["type"];
	        this.modified = source["modified"];
	        this.isDir = source["isDir"];
	    }
	}

}

