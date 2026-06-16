<script lang="ts">
	import type { HistoryPoint } from '$lib';

	interface Props {
		memoryHistory: HistoryPoint[];
	}

	let { memoryHistory }: Props = $props();

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
</script>

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
				<span class="legend-item"><span class="legend-color alloc"></span>RAM Sử dụng (Alloc)</span>
				<span class="legend-item"><span class="legend-color sys"></span>RAM Hệ thống cấp (Sys)</span
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

<style>
	.memory-chart-card {
		position: relative;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 280px;
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
</style>
