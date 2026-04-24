export namespace main {
	
	export class DorkEntry {
	    id: string;
	    title: string;
	    query: string;
	    category: string;
	    tags: string;
	
	    static createFrom(source: any = {}) {
	        return new DorkEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.query = source["query"];
	        this.category = source["category"];
	        this.tags = source["tags"];
	    }
	}
	export class HistoryEntry {
	    id: string;
	    domain: string;
	    dorkTitle: string;
	    query: string;
	    fullQuery: string;
	    timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new HistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.domain = source["domain"];
	        this.dorkTitle = source["dorkTitle"];
	        this.query = source["query"];
	        this.fullQuery = source["fullQuery"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class SearchResult {
	    success: boolean;
	    message: string;
	    fullUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.fullUrl = source["fullUrl"];
	    }
	}

}

