export namespace models {
	
	export class Forum {
	    forum_id: string;
	    forum_url: string;
	    forum_name: string;
	    forum_description: string;
	    last_scaned: string;
	    forum_html: string;
	    forum_screenshot: string;
	
	    static createFrom(source: any = {}) {
	        return new Forum(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.forum_id = source["forum_id"];
	        this.forum_url = source["forum_url"];
	        this.forum_name = source["forum_name"];
	        this.forum_description = source["forum_description"];
	        this.last_scaned = source["last_scaned"];
	        this.forum_html = source["forum_html"];
	        this.forum_screenshot = source["forum_screenshot"];
	    }
	}

}

