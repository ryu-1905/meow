<script lang="ts">
	import type { HealthResponse } from '$lib';

	interface Props {
		healthData: HealthResponse | null;
		error: string | null;
	}

	let { healthData, error }: Props = $props();

	let logsContainer = $state<HTMLDivElement | null>(null);

	// Định dạng màu sắc log trong console
	const formatLogLine = (log: string): string => {
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
	};

	// Tự động cuộn logs container xuống cuối
	$effect(() => {
		if (healthData && healthData.app_logs && logsContainer) {
			const container = logsContainer;
			setTimeout(() => {
				container.scrollTop = container.scrollHeight;
			}, 50);
		}
	});
</script>

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
					<span class="log-arrow">&gt;</span> ERROR: Không thể tải log. Kết nối đến server thất bại.
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

<style>
	.console-logs-card {
		position: relative;
		overflow: hidden;
		display: flex;
		flex-direction: column;
		min-height: 200px;
		flex-grow: 1;
	}

	.card-title-with-icon {
		display: flex;
		align-items: center;
		gap: 10px;
		font-family: var(--font-display);
		font-size: 1.25rem;
		font-weight: 700;
		color: var(--text-primary);
		margin-bottom: 20px;
	}

	.title-icon {
		color: var(--accent-cyan);
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

	.text-danger {
		color: var(--accent-red) !important;
	}

	.text-muted {
		color: var(--text-muted) !important;
	}

	.font-semibold {
		font-weight: 600;
	}
</style>
