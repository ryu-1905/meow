<script lang="ts">
	import { PUBLIC_BACKEND_URL } from '$env/static/public';
	import type { HealthResponse, HistoryPoint } from '$lib';

	// Import các components con
	import SystemStats from '$lib/components/SystemStats.svelte';
	import MemoryChart from '$lib/components/MemoryChart.svelte';
	import RamGauge from '$lib/components/RamGauge.svelte';
	import ConsoleLogs from '$lib/components/ConsoleLogs.svelte';

	// Svelte 5 Runes
	let healthData = $state<HealthResponse | null>(null);
	let error = $state<string | null>(null);
	let lastRefreshed = $state<string>('');
	let memoryHistory = $state<HistoryPoint[]>([]);

	// Parse RAM string (ví dụ "2.50 MB") -> số thực
	const parseMemory = (memStr: string): number => {
		if (!memStr) return 0;
		return parseFloat(memStr.replace(/[^\d.]/g, '')) || 0;
	};

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
			<SystemStats {healthData} {error} onRefresh={fetchHealth} />
			<MemoryChart {memoryHistory} />
		</div>

		<!-- CỘT PHẢI: Biểu đồ Gauge & Logs máy chủ -->
		<div class="col-right">
			<RamGauge {healthData} {error} />
			<ConsoleLogs {healthData} {error} />
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
</style>
