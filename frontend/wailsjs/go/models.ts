export namespace registry {
	
	export class File {
	    Name: string;
	    Path: string;
	    Size: number;
	    Uploaded_by: string;
	    Date: number;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Path = source["Path"];
	        this.Size = source["Size"];
	        this.Uploaded_by = source["Uploaded_by"];
	        this.Date = source["Date"];
	    }
	}

}

