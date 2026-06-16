<script lang="ts">
	import { PUBLIC_BACKEND_URL } from '$env/static/public';
	import type { HealthResponse } from '$lib';

	// Định nghĩa props sử dụng Svelte 5 $props
	interface Props {
		healthData: HealthResponse | null;
		error: string | null;
		onRefresh: () => void;
	}

	let { healthData, error, onRefresh }: Props = $props();
</script>

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
		<button class="refresh-btn" onclick={onRefresh} aria-label="Làm mới">
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
					<div
						class="stat-icon-wrapper {healthData.database_status === 'CONNECTED' ? 'green' : 'red'}"
					>
						<svg viewBox="0 0 24 24" width="20" height="20">
							<path
								fill="currentColor"
								d="M12 2C7.58 2 4 3.79 4 6v12c0 2.21 3.58 4 8 4s8-1.79 8-4V6c0-2.21-3.58-4-8-4zm0 2c3.87 0 6 1.34 6 2s-2.13 2-6 2s-6-1.34-6-2s2.13-2 6-2zm0 14c-3.87 0-6-1.34-6-2v-2.56c1.47.88 3.55 1.56 6 1.56s4.53-.68 6-1.56V16c0 .66-2.13 2-6 2zm0-5c-3.87 0-6-1.34-6-2V8.44c1.47.88 3.55 1.56 6 1.56s4.53-.68 6-1.56V11c0 .66-2.13 2-6 2z"
							/>
						</svg>
					</div>
					<div class="stat-info">
						<span class="stat-label">Cơ sở dữ liệu (Postgres)</span>
						<span
							class="stat-value {healthData.database_status === 'CONNECTED'
								? 'green-text'
								: 'red-text'}"
						>
							{healthData.database_status}
						</span>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	.system-stats-card {
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
	.stat-icon-wrapper.red {
		background: rgba(255, 75, 92, 0.1);
		color: var(--accent-red);
	}

	.stat-info {
		display: flex;
		flex-direction: column;
		gap: 4px;
		min-width: 0;
	}

	.stat-label {
		color: var(--text-muted);
		font-size: 0.72rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.stat-value {
		color: var(--text-primary);
		font-weight: 600;
		font-size: 0.95rem;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.green-text {
		color: var(--accent-green) !important;
	}

	.red-text {
		color: var(--accent-red) !important;
	}
</style>
