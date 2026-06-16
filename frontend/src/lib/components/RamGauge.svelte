<script lang="ts">
	import type { HealthResponse } from '$lib';

	interface Props {
		healthData: HealthResponse | null;
		error: string | null;
	}

	let { healthData, error }: Props = $props();

	// Parse RAM string (ví dụ "2.50 MB") -> số thực
	const parseMemory = (memStr: string): number => {
		if (!memStr) return 0;
		return parseFloat(memStr.replace(/[^\d.]/g, '')) || 0;
	};

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
</script>

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
					<span class="label">Tổng tích luỹ (Total):</span>
					<span class="val text-muted">{healthData.server_info.memory_usage.total_alloc}</span>
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

<style>
	.ram-gauge-card {
		position: relative;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 200px;
		align-items: center;
	}

	.card-title-with-icon {
		display: flex;
		align-items: center;
		gap: 10px;
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-primary);
		align-self: flex-start;
	}

	.title-icon {
		color: var(--accent-cyan);
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
</style>
