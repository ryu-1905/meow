<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import { auth } from '$lib/auth.svelte';
	import AuthForm from '$lib/components/AuthForm.svelte';
	import Header from '$lib/components/Header.svelte';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>🐱 Meow Dashboard</title>
</svelte:head>

<div class="app-container">
	{#if auth.isAuthenticated}
		<Header />

		<main class="main-content">
			{@render children()}
		</main>

		<footer class="app-footer">
			<p class="footer-text">Built with 💖 using Svelte 5, Go, & Antigravity IDE</p>
		</footer>
	{:else}
		<AuthForm />
	{/if}
</div>

<style>
	.app-container {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	.main-content {
		flex: 1;
		width: 100%;
		max-width: 1200px;
		margin: 0 auto;
		padding: 40px 24px;
		transition: var(--transition-smooth);
	}

	.app-footer {
		text-align: center;
		padding: 24px;
		border-top: 1px solid var(--glass-border);
		background: rgba(6, 7, 11, 0.5);
	}

	.footer-text {
		color: var(--text-muted);
		font-size: 0.85rem;
	}

	@media (max-width: 640px) {
		.main-content {
			padding: 20px 16px;
		}
	}
</style>
