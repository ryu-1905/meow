<script lang="ts">
	import { PUBLIC_BACKEND_URL } from '$env/static/public';

	// Định nghĩa kiểu dữ liệu khớp chính xác với backend Go
	interface MemStatus {
		alloc: string;
		total_alloc: string;
		sys: string;
		num_gc: number;
	}

	interface ServerInfo {
		os: string;
		architecture: string;
		go_version: string;
		num_cpu: number;
		num_goroutine: number;
		uptime: string;
		memory_usage: MemStatus;
	}

	interface HealthResponse {
		status: string;
		database_status: string;
		server_info: ServerInfo;
		app_logs: string[];
	}

	interface HistoryPoint {
		time: string;
		alloc: number;
		sys: number;
	}

	// Svelte 5 Runes
	let healthData = $state<HealthResponse | null>(null);
	let error = $state<string | null>(null);
	let lastRefreshed = $state<string>('');
	let memoryHistory = $state<HistoryPoint[]>([]);
	let logsContainer = $state<HTMLDivElement | null>(null);

	// Parse RAM string (ví dụ "2.50 MB") -> số thực
	function parseMemory(memStr: string): number {
		if (!memStr) return 0;
		return parseFloat(memStr.replace(/[^\d.]/g, '')) || 0;
	}

	const fetchHealth = async () => {
		try {
			error = null;
			const targetUrl = `${PUBLIC_BACKEND_URL}/api/health`;
			const response = await fetch(targetUrl);

			if (!response.ok) {
				throw new Error(`HTTP error! status: ${response.status}`);
			}

			const data: HealthResponse = await response.json();
			healthData = data;
			lastRefreshed = new Date().toLocaleTimeString();

			// Thêm điểm lấy mẫu mới vào lịch sử (giữ tối đa 15 điểm)
			if (data.server_info && data.server_info.memory_usage) {
				const allocVal = parseMemory(data.server_info.memory_usage.alloc);
				const sysVal = parseMemory(data.server_info.memory_usage.sys);
				const now = new Date();
				const timeStr = now.toLocaleTimeString([], {
					hour: '2-digit',
					minute: '2-digit',
					second: '2-digit'
				});

				memoryHistory = [...memoryHistory, { time: timeStr, alloc: allocVal, sys: sysVal }].slice(
					-15
				);
			}
		} catch (err) {
			console.error('Fetch health failed:', err);
			error = (err as Error).message || 'Không thể kết nối đến Backend Server';
			healthData = null;
		}
	};

	// Tự động fetch mỗi 10 giây
	$effect(() => {
		fetchHealth();
		const interval = setInterval(fetchHealth, 10000);
		return () => clearInterval(interval);
	});

	// Tự động cuộn logs container xuống cuối
	$effect(() => {
		if (healthData && healthData.app_logs && logsContainer) {
			const container = logsContainer;
			setTimeout(() => {
				container.scrollTop = container.scrollHeight;
			}, 50);
		}
	});

	// Các giá trị tính toán phái sinh ($derived) cho biểu đồ lịch sử Y-axis
	let maxMem = $derived(
		memoryHistory.length > 0
			? Math.max(10, Math.max(...memoryHistory.map((p) => Math.max(p.alloc, p.sys))) * 1.15)
			: 30
	);

	let allocPoints = $derived(
		memoryHistory.map((p, idx) => {
			const x = 50 + (idx / Math.max(1, memoryHistory.length - 1)) * 430;
			const y = 170 - (p.alloc / maxMem) * 140;
			return { x, y, point: p };
		})
	);

	let sysPoints = $derived(
		memoryHistory.map((p, idx) => {
			const x = 50 + (idx / Math.max(1, memoryHistory.length - 1)) * 430;
			const y = 170 - (p.sys / maxMem) * 140;
			return { x, y, point: p };
		})
	);

	let allocPath = $derived(
		allocPoints.length > 0 ? 'M ' + allocPoints.map((pt) => `${pt.x} ${pt.y}`).join(' L ') : ''
	);

	let allocAreaPath = $derived(
		allocPoints.length > 0
			? `M ${allocPoints[0].x} 170 L ` +
					allocPoints.map((pt) => `${pt.x} ${pt.y}`).join(' L ') +
					` L ${allocPoints[allocPoints.length - 1].x} 170 Z`
			: ''
	);

	let sysPath = $derived(
		sysPoints.length > 0 ? 'M ' + sysPoints.map((pt) => `${pt.x} ${pt.y}`).join(' L ') : ''
	);

	let sysAreaPath = $derived(
		sysPoints.length > 0
			? `M ${sysPoints[0].x} 170 L ` +
					sysPoints.map((pt) => `${pt.x} ${pt.y}`).join(' L ') +
					` L ${sysPoints[sysPoints.length - 1].x} 170 Z`
			: ''
	);

	let yGridValues = $derived([0, maxMem * 0.25, maxMem * 0.5, maxMem * 0.75, maxMem]);

	// Tỷ lệ phần trăm bộ nhớ RAM thực tế đang dùng
	let ramPercentage = $derived(
		healthData && healthData.server_info
			? Math.min(
					100,
					Math.round(
						(parseMemory(healthData.server_info.memory_usage.alloc) /
							Math.max(1, parseMemory(healthData.server_info.memory_usage.sys))) *
							100
					)
				)
			: 0
	);

	// Định dạng màu sắc log trong console
	function formatLogLine(log: string): string {
		if (!log) return '';
		let escaped = log.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
		return escaped
			.replace(/\b(GET)\b/g, '<span class="log-method get">$1</span>')
			.replace(/\b(POST)\b/g, '<span class="log-method post">$1</span>')
			.replace(/\b(PUT)\b/g, '<span class="log-method put">$1</span>')
			.replace(/\b(DELETE)\b/g, '<span class="log-method delete">$1</span>')
			.replace(/\b(INFO)\b/g, '<span class="log-level info">$1</span>')
			.replace(/\b(WARN(?:ING)?)\b/gi, '<span class="log-level warn">$1</span>')
			.replace(/\b(ERROR|FAIL(?:ED)?|ERR)\b/gi, '<span class="log-level error">$1</span>')
			.replace(/\b(200|201)\b/g, '<span class="log-status success">$1</span>')
			.replace(/\b(400|401|403|404)\b/g, '<span class="log-status warning">$1</span>')
			.replace(/\b(500|502|503|504)\b/g, '<span class="log-status danger">$1</span>')
			.replace(/(\/api\/[a-zA-Z0-9_/]+)/g, '<span class="log-url">$1</span>');
	}
</script>

<div class="dashboard">
	<!-- Hero Section -->
	<div class="hero-section">
		<div class="title-container">
			<h1 class="gradient-text dashboard-title animate-title">Meow Backend Monitor</h1>
			{#if error}
				<span class="badge badge-error pulse-error shadow-error font-semibold">OFFLINE</span>
			{:else if healthData}
				<span class="badge badge-success pulse-success shadow-success font-semibold">ONLINE</span>
				{#if lastRefreshed}
					<span class="text-xs text-muted font-mono" style="margin-left: 10px; opacity: 0.8;"
						>Cập nhật lúc: {lastRefreshed}</span
					>
				{/if}
			{/if}
		</div>
		<p class="dashboard-subtitle">
			Hệ thống giám sát sức khỏe dịch vụ và hiệu suất tài nguyên trực quan
		</p>
	</div>

	<!-- Dashboard Bố cục Grid -->
	<div class="dashboard-grid">
		<!-- CỘT TRÁI: Thông số hệ thống & Lịch sử RAM -->
		<div class="col-left">
			<!-- Card 1: Thông số Hệ thống & Runtime -->
			<div class="glass-card system-stats-card">
				<div class="card-header">
					<h2 class="card-title-with-icon">
						<svg class="title-icon" viewBox="0 0 24 24" width="22" height="22">
							<path
								fill="currentColor"
								d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-2 10h-4v4h-2v-4H7v-2h4V7h2v4h4v2z"
							/>
						</svg>
						Hệ thống & Runtime Go
					</h2>
					<button class="refresh-btn" onclick={fetchHealth} aria-label="Làm mới">
						<svg class="refresh-icon" viewBox="0 0 24 24" width="18" height="18">
							<path
								fill="currentColor"
								d="M17.65 6.35A7.958 7.958 0 0 0 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08A5.99 5.99 0 0 1 12 18c-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"
							/>
						</svg>
					</button>
				</div>

				<div class="stats-content">
					{#if error}
						<div class="offline-alert">
							<div class="offline-icon-container">
								<svg class="offline-alert-icon" viewBox="0 0 24 24" width="48" height="48">
									<path
										fill="currentColor"
										d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"
									/>
								</svg>
							</div>
							<p class="error-msg">{error}</p>
							<p class="error-hint">
								Hãy đảm bảo máy chủ Backend đang chạy ở <!-- eslint-disable-next-line svelte/no-navigation-without-resolve --><a
									href={PUBLIC_BACKEND_URL}
									target="_blank"
									rel="noopener noreferrer">{PUBLIC_BACKEND_URL}</a
								>
							</p>
						</div>
					{:else if healthData}
						<div class="mini-stats-grid">
							<!-- HĐH & Kiến trúc -->
							<div class="stat-item">
								<div class="stat-icon-wrapper cyan">
									<svg viewBox="0 0 24 24" width="20" height="20">
										<path
											fill="currentColor"
											d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"
										/>
									</svg>
								</div>
								<div class="stat-info">
									<span class="stat-label">Hệ điều hành / Arch</span>
									<span class="stat-value"
										>{healthData.server_info.os} ({healthData.server_info.architecture})</span
									>
								</div>
							</div>

							<!-- Go Version -->
							<div class="stat-item">
								<div class="stat-icon-wrapper purple">
									<svg viewBox="0 0 24 24" width="20" height="20">
										<path
											fill="currentColor"
											d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.53c-.26-.81-1-1.4-1.9-1.4h-1v-3c0-.55-.45-1-1-1h-6v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"
										/>
									</svg>
								</div>
								<div class="stat-info">
									<span class="stat-label">Phiên bản Go</span>
									<span class="stat-value font-mono">{healthData.server_info.go_version}</span>
								</div>
							</div>

							<!-- CPUs Cores -->
							<div class="stat-item">
								<div class="stat-icon-wrapper green">
									<svg viewBox="0 0 24 24" width="20" height="20">
										<path
											fill="currentColor"
											d="M4 6H2v2h2v2H2v2h2v2H2v2h2v2h2v2h2v-2h2v2h2v-2h2v2h2v-2h2v2h2v-2h-2v-2h2v-2h-2v-2h2v-2h-2V8h2V6h-2V4h-2v2h-2V4h-2v2h-2V4H8v2H6V4H4v2zm4 4h8v8H8v-8z"
										/>
									</svg>
								</div>
								<div class="stat-info">
									<span class="stat-label">Số luồng CPU</span>
									<span class="stat-value">{healthData.server_info.num_cpu} Cores</span>
								</div>
							</div>

							<!-- Goroutines -->
							<div class="stat-item">
								<div class="stat-icon-wrapper gold">
									<svg viewBox="0 0 24 24" width="20" height="20">
										<path
											fill="currentColor"
											d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17h-2v-2h2v2zm2.07-7.75l-.9.92C13.45 12.9 13 13.5 13 15h-2v-.5c0-1.1.45-2.1 1.17-2.83l1.24-1.26c.37-.36.59-.86.59-1.41 0-1.1-.9-2-2-2s-2 .9-2 2H7c0-2.76 2.24-5 5-5s5 2.24 5 5c0 1.04-.42 1.99-1.07 2.75z"
										/>
									</svg>
								</div>
								<div class="stat-info">
									<span class="stat-label">Goroutines đang chạy</span>
									<span class="stat-value font-mono">{healthData.server_info.num_goroutine}</span>
								</div>
							</div>

							<!-- Uptime -->
							<div class="stat-item">
								<div class="stat-icon-wrapper accent">
									<svg viewBox="0 0 24 24" width="20" height="20">
										<path
											fill="currentColor"
											d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm.5-13H11v6l5.25 3.15.75-1.23-4.5-2.67V7z"
										/>
									</svg>
								</div>
								<div class="stat-info">
									<span class="stat-label">Thời gian đã chạy (Uptime)</span>
									<span class="stat-value font-mono">{healthData.server_info.uptime}</span>
								</div>
							</div>

							<!-- Cơ sở dữ liệu -->
							<div class="stat-item">
								<div class="stat-icon-wrapper {healthData.database_status === 'CONNECTED' ? 'green' : 'red'}">
									<svg viewBox="0 0 24 24" width="20" height="20">
										<path
											fill="currentColor"
											d="M12 2C7.58 2 4 3.79 4 6v12c0 2.21 3.58 4 8 4s8-1.79 8-4V6c0-2.21-3.58-4-8-4zm0 2c3.87 0 6 1.34 6 2s-2.13 2-6 2s-6-1.34-6-2s2.13-2 6-2zm0 14c-3.87 0-6-1.34-6-2v-2.56c1.47.88 3.55 1.56 6 1.56s4.53-.68 6-1.56V16c0 .66-2.13 2-6 2zm0-5c-3.87 0-6-1.34-6-2V8.44c1.47.88 3.55 1.56 6 1.56s4.53-.68 6-1.56V11c0 .66-2.13 2-6 2z"
										/>
									</svg>
								</div>
								<div class="stat-info">
									<span class="stat-label">Cơ sở dữ liệu (Postgres)</span>
									<span class="stat-value {healthData.database_status === 'CONNECTED' ? 'green-text' : 'red-text'}">
										{healthData.database_status}
									</span>
								</div>
							</div>
						</div>
					{/if}
				</div>
			</div>

			<!-- Card 2: Biểu đồ Area Line Chart lịch sử bộ nhớ -->
			<div class="glass-card memory-chart-card">
				<h2 class="card-title-with-icon">
					<svg class="title-icon" viewBox="0 0 24 24" width="22" height="22">
						<path
							fill="currentColor"
							d="M3.5 18.49l6-6.01 4 4L22 6.92l-1.41-1.41-7.09 7.97-4-4L2 16.99l1.5 1.5z"
						/>
					</svg>
					Lịch sử Bộ nhớ RAM (Alloc vs Sys)
				</h2>

				<div class="chart-container">
					{#if memoryHistory.length === 0}
						<div class="chart-empty">
							<p>Đang tích lũy điểm dữ liệu...</p>
							<span class="pulse-dot"></span>
						</div>
					{:else}
						<div class="chart-legend">
							<span class="legend-item"
								><span class="legend-color alloc"></span>RAM Sử dụng (Alloc)</span
							>
							<span class="legend-item"
								><span class="legend-color sys"></span>RAM Hệ thống cấp (Sys)</span
							>
						</div>

						<svg class="svg-chart" viewBox="0 0 500 220" width="100%" height="220">
							<defs>
								<linearGradient id="allocGrad" x1="0" y1="0" x2="0" y2="1">
									<stop offset="0%" stop-color="var(--accent-cyan)" stop-opacity="0.35" />
									<stop offset="100%" stop-color="var(--accent-cyan)" stop-opacity="0.0" />
								</linearGradient>
								<linearGradient id="sysGrad" x1="0" y1="0" x2="0" y2="1">
									<stop offset="0%" stop-color="var(--accent-purple)" stop-opacity="0.15" />
									<stop offset="100%" stop-color="var(--accent-purple)" stop-opacity="0.0" />
								</linearGradient>
							</defs>

							<!-- Các đường lưới Grid ngang -->
							{#each yGridValues as val (val)}
								{@const y = 170 - (val / maxMem) * 140}
								<line
									x1="50"
									x2="480"
									y1={y}
									y2={y}
									stroke="rgba(255,255,255,0.06)"
									stroke-dasharray="3,3"
								/>
								<text x="42" y={y + 3} fill="var(--text-muted)" font-size="9" text-anchor="end"
									>{val.toFixed(1)} M</text
								>
							{/each}

							<!-- Trục X mờ -->
							<line x1="50" x2="480" y1="170" y2="170" stroke="rgba(255,255,255,0.15)" />

							<!-- Vẽ vùng bao phủ (Area) -->
							{#if sysAreaPath}
								<path d={sysAreaPath} fill="url(#sysGrad)" />
							{/if}
							{#if allocAreaPath}
								<path d={allocAreaPath} fill="url(#allocGrad)" />
							{/if}

							<!-- Vẽ đường (Line) -->
							{#if sysPath}
								<path
									d={sysPath}
									fill="none"
									stroke="var(--accent-purple)"
									stroke-width="2.5"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
							{/if}
							{#if allocPath}
								<path
									d={allocPath}
									fill="none"
									stroke="var(--accent-cyan)"
									stroke-width="2.5"
									stroke-linecap="round"
									stroke-linejoin="round"
								/>
							{/if}

							<!-- Các chấm tròn tại các điểm lấy mẫu -->
							{#each allocPoints as pt, idx (idx)}
								<circle
									cx={pt.x}
									cy={pt.y}
									r="3.5"
									fill="var(--bg-primary)"
									stroke="var(--accent-cyan)"
									stroke-width="2"
								/>
							{/each}

							<!-- Nhãn trục X thời gian (Hiển thị giãn cách để tránh chồng chéo) -->
							{#each allocPoints as pt, idx (idx)}
								{#if idx % 3 === 0 || idx === memoryHistory.length - 1}
									<text
										x={pt.x}
										y="192"
										fill="var(--text-muted)"
										font-size="9"
										text-anchor="middle"
										font-family="monospace">{pt.point.time}</text
									>
									<line x1={pt.x} x2={pt.x} y1="170" y2="174" stroke="rgba(255,255,255,0.2)" />
								{/if}
							{/each}
						</svg>
					{/if}
				</div>
			</div>
		</div>

		<!-- CỘT PHẢI: Biểu đồ Gauge & Logs máy chủ -->
		<div class="col-right">
			<!-- Card 3: RAM Donut Gauge Chart -->
			<div class="glass-card ram-gauge-card">
				<h2 class="card-title-with-icon">
					<svg class="title-icon" viewBox="0 0 24 24" width="22" height="22">
						<path
							fill="currentColor"
							d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 17.93c-3.95-.49-7-3.85-7-7.93h2c0 3.14 2.5 5.5 5 5v2.93zM12 4c3.95 0 7 3.05 7 7h-2c0-3.14-2.5-5.5-5-5V4.07z"
						/>
					</svg>
					Tải lượng RAM Hiện tại
				</h2>

				<div class="gauge-section">
					{#if error}
						<div class="gauge-placeholder offline">
							<span>Mất kết nối</span>
						</div>
					{:else if healthData}
						<div class="gauge-container">
							<svg class="gauge-svg" width="160" height="160" viewBox="0 0 160 160">
								<defs>
									<linearGradient id="gaugeGrad" x1="0" y1="0" x2="1" y2="1">
										<stop offset="0%" stop-color="var(--accent-cyan)" />
										<stop offset="100%" stop-color="var(--accent-purple)" />
									</linearGradient>
								</defs>
								<!-- Vòng tròn xám làm nền (100%) -->
								<circle class="gauge-bg" cx="80" cy="80" r="66" stroke-width="12" fill="none" />
								<!-- Vòng tiến trình % thực tế -->
								<!-- Chu vi = 2 * PI * r = 2 * 3.14159 * 66 = 414.7 -->
								<circle
									class="gauge-progress"
									cx="80"
									cy="80"
									r="66"
									stroke-width="12"
									fill="none"
									stroke="url(#gaugeGrad)"
									stroke-dasharray="414.7"
									stroke-dashoffset={414.7 - (ramPercentage / 100) * 414.7}
									stroke-linecap="round"
								/>
							</svg>
							<div class="gauge-text-overlay">
								<span class="gauge-value">{ramPercentage}%</span>
								<span class="gauge-label">RAM Occupied</span>
							</div>
						</div>

						<div class="gauge-stats">
							<div class="gauge-stat-row">
								<span class="label">Đang dùng (Alloc):</span>
								<span class="val cyan-text">{healthData.server_info.memory_usage.alloc}</span>
							</div>
							<div class="gauge-stat-row">
								<span class="label">Hệ thống cấp (Sys):</span>
								<span class="val purple-text">{healthData.server_info.memory_usage.sys}</span>
							</div>
							<div class="gauge-stat-row">
								<span class="label">Tổng tích lũy (Total):</span>
								<span class="val text-muted">{healthData.server_info.memory_usage.total_alloc}</span
								>
							</div>
							<div class="gauge-stat-row">
								<span class="label">Chu kỳ dọn rác (GC):</span>
								<span class="val gold-text">{healthData.server_info.memory_usage.num_gc} lần</span>
							</div>
						</div>
					{:else}
						<div class="gauge-placeholder">
							<span>Đang tải...</span>
						</div>
					{/if}
				</div>
			</div>

			<!-- Card 4: Terminal Live Logs Console -->
			<div class="glass-card console-logs-card">
				<h2 class="card-title-with-icon">
					<svg class="title-icon" viewBox="0 0 24 24" width="22" height="22">
						<path
							fill="currentColor"
							d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm-8 12H4v-2h8v2zm8-4H4V8h16v4z"
						/>
					</svg>
					Nhật ký Hệ thống (Server Live Logs)
				</h2>

				<div class="terminal-wrapper">
					<div class="terminal-header">
						<div class="window-dot red"></div>
						<div class="window-dot yellow"></div>
						<div class="window-dot green"></div>
						<span class="terminal-title">bash - meow-backend.log</span>
					</div>
					<div class="terminal-body" bind:this={logsContainer}>
						{#if error}
							<div class="log-line text-danger font-semibold">
								<span class="log-arrow">&gt;</span> ERROR: Không thể tải log. Kết nối đến server thất
								bại.
							</div>
						{:else if healthData && healthData.app_logs}
							{#if healthData.app_logs.length === 0}
								<div class="log-line text-muted">
									<span class="log-arrow">&gt;</span> Chờ ghi nhận log hoạt động mới...
								</div>
							{:else}
								{#each healthData.app_logs as log, idx (idx)}
									<div class="log-line">
										<span class="log-arrow">&gt;</span>
										<!-- eslint-disable-next-line svelte/no-at-html-tags -->
										<span class="log-content">{@html formatLogLine(log)}</span>
									</div>
								{/each}
							{/if}
						{:else}
							<div class="log-line text-muted">
								<span class="log-arrow">&gt;</span> Đang khởi tạo kết nối terminal...
							</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	/* Grid Layout */
	.dashboard {
		display: flex;
		flex-direction: column;
		gap: 30px;
		max-width: 1200px;
		margin: 0 auto;
		padding: 20px;
		transition: var(--transition-smooth);
	}

	.hero-section {
		display: flex;
		flex-direction: column;
		gap: 10px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		padding-bottom: 25px;
	}

	.title-container {
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: 15px;
	}

	.dashboard-title {
		font-family: var(--font-display);
		font-weight: 800;
		font-size: 2.8rem;
		letter-spacing: -0.03em;
		transition: var(--transition-smooth);
	}

	.dashboard-subtitle {
		color: var(--text-secondary);
		font-size: 1.1rem;
		transition: var(--transition-smooth);
	}

	.dashboard-grid {
		display: grid;
		grid-template-columns: 1.1fr 0.9fr;
		gap: 25px;
		transition: var(--transition-smooth);
	}

	.col-left,
	.col-right {
		display: flex;
		flex-direction: column;
		gap: 25px;
	}

	/* Responsive Tablet & Mobile */
	@media (max-width: 1024px) {
		.dashboard {
			padding: 16px;
			gap: 20px;
		}
		.dashboard-grid {
			gap: 20px;
		}
		.glass-card {
			padding: 20px;
		}
	}

	@media (max-width: 768px) {
		.dashboard-grid {
			grid-template-columns: 1fr;
		}
		.dashboard-title {
			font-size: 2rem;
		}
		.dashboard-subtitle {
			font-size: 0.95rem;
		}
		.col-left,
		.col-right {
			gap: 20px;
		}
	}

	@media (max-width: 480px) {
		.dashboard {
			padding: 12px;
		}
		.dashboard-title {
			font-size: 1.8rem;
		}
		.glass-card {
			padding: 16px;
		}
		.card-title-with-icon {
			font-size: 1.15rem;
		}
		.terminal-wrapper {
			height: 200px;
		}
	}

	/* Card common styling overrides */
	.glass-card {
		position: relative;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 200px;
	}

	.card-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 20px;
	}

	.card-title-with-icon {
		display: flex;
		align-items: center;
		gap: 10px;
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-primary);
	}

	.title-icon {
		color: var(--accent-cyan);
	}

	/* Refresh Button */
	.refresh-btn {
		background: rgba(255, 255, 255, 0.05);
		border: 1px solid var(--glass-border);
		color: var(--text-secondary);
		width: 36px;
		height: 36px;
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		transition: var(--transition-fast);
	}

	.refresh-btn:hover {
		background: rgba(79, 172, 254, 0.15);
		color: var(--accent-cyan);
		border-color: rgba(0, 242, 254, 0.3);
		transform: rotate(45deg);
	}

	.refresh-btn:active {
		transform: scale(0.9);
	}

	/* Offline State Alert */
	.offline-alert {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		padding: 30px 10px;
		gap: 12px;
		min-height: 220px;
	}

	.offline-icon-container {
		color: var(--accent-red);
		animation: pulseGlowError 2s infinite;
		background: rgba(255, 75, 92, 0.1);
		border: 1px solid rgba(255, 75, 92, 0.2);
		border-radius: 50%;
		width: 80px;
		height: 80px;
		display: flex;
		align-items: center;
		justify-content: center;
		margin-bottom: 10px;
	}

	.error-msg {
		color: var(--accent-red);
		font-weight: 600;
		font-size: 1.1rem;
	}

	.error-hint {
		color: var(--text-muted);
		font-size: 0.88rem;
		max-width: 320px;
	}

	/* Mini Stats Grid */
	.mini-stats-grid {
		display: grid;
		grid-template-columns: repeat(2, 1fr);
		gap: 15px;
	}

	.stat-item {
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid rgba(255, 255, 255, 0.04);
		border-radius: var(--radius-sm);
		padding: 15px;
		display: flex;
		align-items: center;
		gap: 15px;
		transition: var(--transition-fast);
	}

	.stat-item:hover {
		background: rgba(255, 255, 255, 0.05);
		border-color: rgba(255, 255, 255, 0.08);
		transform: translateY(-2px);
	}

	@media (max-width: 600px) {
		.mini-stats-grid {
			grid-template-columns: 1fr;
		}
	}

	.stat-icon-wrapper {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 42px;
		height: 42px;
		border-radius: 10px;
	}

	.stat-icon-wrapper.cyan {
		background: rgba(0, 242, 254, 0.1);
		color: var(--accent-cyan);
	}
	.stat-icon-wrapper.purple {
		background: rgba(79, 172, 254, 0.1);
		color: var(--accent-purple);
	}
	.stat-icon-wrapper.green {
		background: rgba(56, 239, 125, 0.1);
		color: var(--accent-green);
	}
	.stat-icon-wrapper.gold {
		background: rgba(248, 173, 157, 0.1);
		color: var(--accent-gold);
	}
	.stat-icon-wrapper.accent {
		background: rgba(79, 172, 254, 0.15);
		color: var(--text-primary);
	}

	.stat-info {
		display: flex;
		flex-direction: column;
		gap: 4px;
	}

	.stat-label {
		color: var(--text-muted);
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
	}

	.stat-value {
		color: var(--text-primary);
		font-weight: 600;
		font-size: 0.95rem;
	}

	/* SVG Area Line Chart */
	.memory-chart-card {
		min-height: 280px;
	}

	.chart-container {
		position: relative;
		margin-top: 10px;
		width: 100%;
		height: 100%;
		display: flex;
		flex-direction: column;
		justify-content: center;
	}

	.chart-empty {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		gap: 15px;
		min-height: 180px;
		color: var(--text-muted);
		font-size: 0.95rem;
	}

	.pulse-dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		background-color: var(--accent-cyan);
		box-shadow: 0 0 10px var(--accent-cyan);
		animation: pulseGlow 1.5s infinite;
	}

	.chart-legend {
		display: flex;
		gap: 20px;
		justify-content: flex-end;
		margin-bottom: 12px;
		font-size: 0.82rem;
	}

	.legend-item {
		display: flex;
		align-items: center;
		gap: 6px;
		color: var(--text-secondary);
	}

	.legend-color {
		width: 12px;
		height: 12px;
		border-radius: 3px;
	}

	.legend-color.alloc {
		background-color: var(--accent-cyan);
	}
	.legend-color.sys {
		background-color: var(--accent-purple);
	}

	.svg-chart {
		overflow: visible;
	}

	/* RAM Donut Gauge */
	.ram-gauge-card {
		align-items: center;
	}

	.gauge-section {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 20px;
		width: 100%;
		padding: 10px 0;
	}

	.gauge-container {
		position: relative;
		width: 160px;
		height: 160px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.gauge-svg {
		transform: rotate(-90deg);
	}

	.gauge-bg {
		stroke: rgba(255, 255, 255, 0.04);
	}

	.gauge-progress {
		transition: stroke-dashoffset 0.8s cubic-bezier(0.16, 1, 0.3, 1);
		filter: drop-shadow(0 0 4px rgba(0, 242, 254, 0.2));
	}

	.gauge-text-overlay {
		position: absolute;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
	}

	.gauge-value {
		font-family: var(--font-display);
		font-size: 2.2rem;
		font-weight: 800;
		color: var(--text-primary);
		line-height: 1;
	}

	.gauge-label {
		font-size: 0.72rem;
		color: var(--text-muted);
		text-transform: uppercase;
		letter-spacing: 0.05em;
		margin-top: 4px;
	}

	.gauge-placeholder {
		width: 160px;
		height: 160px;
		border-radius: 50%;
		border: 8px solid rgba(255, 255, 255, 0.04);
		display: flex;
		align-items: center;
		justify-content: center;
		color: var(--text-muted);
		font-size: 0.9rem;
	}

	.gauge-placeholder.offline {
		border-color: rgba(255, 75, 92, 0.1);
		color: var(--accent-red);
		font-weight: 600;
		background: rgba(255, 75, 92, 0.02);
	}

	.gauge-stats {
		display: flex;
		flex-direction: column;
		width: 100%;
		gap: 10px;
		border-top: 1px solid rgba(255, 255, 255, 0.05);
		padding-top: 15px;
	}

	.gauge-stat-row {
		display: flex;
		justify-content: space-between;
		font-size: 0.88rem;
	}

	.gauge-stat-row .label {
		color: var(--text-secondary);
	}

	.gauge-stat-row .val {
		font-weight: 600;
		font-family: monospace;
	}

	.cyan-text {
		color: var(--accent-cyan);
	}
	.purple-text {
		color: var(--accent-purple);
	}
	.gold-text {
		color: var(--accent-gold);
	}

	/* Terminal Console logs */
	.console-logs-card {
		flex-grow: 1;
	}

	.terminal-wrapper {
		background: #06070a;
		border: 1px solid rgba(255, 255, 255, 0.08);
		border-radius: 8px;
		display: flex;
		flex-direction: column;
		height: 250px;
		box-shadow: inset 0 0 10px rgba(0, 0, 0, 0.8);
		overflow: hidden;
	}

	.terminal-header {
		background: #11131a;
		padding: 8px 12px;
		display: flex;
		align-items: center;
		gap: 6px;
		border-bottom: 1px solid rgba(255, 255, 255, 0.05);
	}

	.window-dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
	}

	.window-dot.red {
		background-color: #ff5f56;
	}
	.window-dot.yellow {
		background-color: #ffbd2e;
	}
	.window-dot.green {
		background-color: #27c93f;
	}

	.terminal-title {
		color: var(--text-muted);
		font-size: 0.72rem;
		font-family: monospace;
		margin-left: 10px;
	}

	.terminal-body {
		padding: 12px;
		overflow-y: auto;
		font-family: 'Consolas', 'Courier New', Courier, monospace;
		font-size: 0.82rem;
		display: flex;
		flex-direction: column;
		gap: 6px;
		flex-grow: 1;
		scroll-behavior: smooth;
	}

	.log-line {
		display: flex;
		align-items: flex-start;
		gap: 8px;
		line-height: 1.4;
		color: #d1d5db;
		word-break: break-all;
	}

	.log-arrow {
		color: var(--accent-cyan);
		user-select: none;
	}

	.log-content {
		white-space: pre-wrap;
	}

	/* Log syntax highlight classes */
	:global(.log-method) {
		font-weight: 700;
		padding: 1px 4px;
		border-radius: 3px;
		font-size: 0.75rem;
	}
	:global(.log-method.get) {
		background: rgba(0, 242, 254, 0.15);
		color: var(--accent-cyan);
	}
	:global(.log-method.post) {
		background: rgba(79, 172, 254, 0.15);
		color: var(--accent-purple);
	}
	:global(.log-method.put) {
		background: rgba(248, 173, 157, 0.15);
		color: var(--accent-gold);
	}
	:global(.log-method.delete) {
		background: rgba(255, 75, 92, 0.15);
		color: var(--accent-red);
	}

	:global(.log-level) {
		font-weight: bold;
		text-transform: uppercase;
	}
	:global(.log-level.info) {
		color: #38bdf8;
	}
	:global(.log-level.warn) {
		color: #fbbf24;
	}
	:global(.log-level.error) {
		color: #f87171;
	}

	:global(.log-status) {
		font-weight: bold;
	}
	:global(.log-status.success) {
		color: var(--accent-green);
	}
	:global(.log-status.warning) {
		color: var(--accent-gold);
	}
	:global(.log-status.danger) {
		color: var(--accent-red);
	}

	:global(.log-url) {
		color: #c084fc;
		text-decoration: underline;
	}

	/* Animations & Shadows overrides */
	.animate-title {
		animation: float 5s ease-in-out infinite;
	}

	.shadow-success {
		box-shadow: 0 0 10px rgba(56, 239, 125, 0.2);
	}

	.shadow-error {
		box-shadow: 0 0 10px rgba(255, 75, 92, 0.2);
	}

	/* Pulse Animations & Classes for Status Badges */
	.pulse-success {
		animation: pulseGlowSuccess 2s infinite ease-in-out;
	}

	.pulse-error {
		animation: pulseGlowError 2s infinite ease-in-out;
	}

	@keyframes pulseGlowSuccess {
		0%,
		100% {
			box-shadow: 0 0 5px rgba(56, 239, 125, 0.2);
			opacity: 1;
		}
		50% {
			box-shadow: 0 0 15px rgba(56, 239, 125, 0.6);
			opacity: 0.85;
		}
	}

	@keyframes pulseGlowError {
		0%,
		100% {
			box-shadow: 0 0 5px rgba(255, 75, 92, 0.2);
			border-color: rgba(255, 75, 92, 0.2);
			opacity: 1;
		}
		50% {
			box-shadow: 0 0 20px rgba(255, 75, 92, 0.5);
			border-color: rgba(255, 75, 92, 0.5);
			opacity: 0.85;
		}
	}

	.stat-icon-wrapper.red {
		background: rgba(255, 75, 92, 0.1);
		color: var(--accent-red);
	}

	.green-text {
		color: var(--accent-green) !important;
	}

	.red-text {
		color: var(--accent-red) !important;
	}
</style>
