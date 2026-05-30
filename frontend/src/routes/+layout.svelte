<script lang="ts">
	import { resolve } from '$app/paths';
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>🐱 Meow Dashboard</title>
</svelte:head>

<div class="app-container">
	<header class="app-header">
		<div class="header-content">
			<div class="logo">
				<span class="logo-emoji">🐱</span>
				<span class="logo-text gradient-text">Meow App</span>
			</div>
			<nav class="nav-links">
				<a href={resolve('/')} class="nav-link active">Dashboard</a>
			</nav>
		</div>
	</header>

	<main class="main-content">
		{@render children()}
	</main>

	<footer class="app-footer">
		<p class="footer-text">Built with 💖 using Svelte 5, Go, & Antigravity IDE</p>
	</footer>
</div>

<style>
	.app-container {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	.app-header {
		background: rgba(10, 11, 16, 0.6);
		backdrop-filter: blur(12px);
		-webkit-backdrop-filter: blur(12px);
		border-bottom: 1px solid var(--glass-border);
		position: sticky;
		top: 0;
		z-index: 100;
	}

	.header-content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 16px 24px;
		display: flex;
		justify-content: space-between;
		align-items: center;
	}

	.logo {
		display: flex;
		align-items: center;
		gap: 10px;
	}

	.logo-emoji {
		font-size: 1.8rem;
		animation: float 4s ease-in-out infinite;
	}

	.logo-text {
		font-family: var(--font-display);
		font-weight: 800;
		font-size: 1.4rem;
		letter-spacing: -0.02em;
	}

	.nav-links {
		display: flex;
		gap: 24px;
	}

	.nav-link {
		color: var(--text-secondary);
		font-weight: 500;
		font-size: 0.95rem;
		position: relative;
		padding: 4px 0;
	}

	.nav-link.active,
	.nav-link:hover {
		color: var(--text-primary);
	}

	.nav-link.active::after {
		content: '';
		position: absolute;
		bottom: 0;
		left: 0;
		width: 100%;
		height: 2px;
		background: linear-gradient(90deg, var(--accent-cyan), var(--accent-purple));
		border-radius: 2px;
	}

	.main-content {
		flex: 1;
		width: 100%;
		max-width: 1200px;
		margin: 0 auto;
		padding: 40px 24px;
	}

	.app-footer {
		text-align: center;
		padding: 24px;
		border-top: 1px solid var(--glass-border);
		background: rgba(10, 11, 16, 0.4);
	}

	.footer-text {
		color: var(--text-muted);
		font-size: 0.85rem;
	}
</style>
