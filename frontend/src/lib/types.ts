export interface MemStatus {
	alloc: string;
	total_alloc: string;
	sys: string;
	num_gc: number;
}

export interface ServerInfo {
	os: string;
	architecture: string;
	go_version: string;
	num_cpu: number;
	num_goroutine: number;
	uptime: string;
	memory_usage: MemStatus;
}

export interface HealthResponse {
	status: string;
	database_status: string;
	server_info: ServerInfo;
	app_logs: string[];
}

export interface HistoryPoint {
	time: string;
	alloc: number;
	sys: number;
}
