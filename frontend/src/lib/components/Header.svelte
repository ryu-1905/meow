<script lang="ts">
	import { resolve } from '$app/paths';
	import { auth } from '$lib/auth.svelte';
</script>

<header class="app-header">
	<div class="header-content">
		<div class="logo">
			<span class="logo-emoji">🐱</span>
			<span class="logo-text gradient-text">Meow App</span>
		</div>
		<nav class="nav-links">
			<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
			<a href={resolve('/')} class="nav-link active">Dashboard</a>
			<div class="user-profile">
				<div class="user-avatar">
					<svg viewBox="0 0 24 24" width="16" height="16">
						<path
							fill="currentColor"
							d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 4c1.93 0 3.5 1.57 3.5 3.5S13.93 13 12 13s-3.5-1.57-3.5-3.5S10.07 6 12 6zm0 14c-2.03 0-4.43-1-5.48-2.53C7.62 16.21 9.69 15 12 15s4.38 1.21 5.48 2.47C16.43 19 14.03 20 12 20z"
						/>
					</svg>
				</div>
				<span class="user-email">{auth.email}</span>
				<button onclick={() => auth.logout()} class="logout-btn">
					<span>Đăng xuất</span>
					<svg class="logout-icon" viewBox="0 0 24 24" width="16" height="16">
						<path
							fill="currentColor"
							d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.58L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"
						/>
					</svg>
				</button>
			</div>
		</nav>
	</div>
</header>

<style>
	.app-header {
		background: rgba(6, 7, 11, 0.7);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
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
		transition: var(--transition-smooth);
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
		align-items: center;
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

	.user-profile {
		display: flex;
		align-items: center;
		gap: 12px;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid var(--glass-border);
		padding: 6px 14px;
		border-radius: 20px;
		margin-left: 10px;
	}

	.user-avatar {
		color: var(--accent-cyan);
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.user-email {
		color: var(--text-primary);
		font-size: 0.88rem;
		font-weight: 500;
	}

	.logout-btn {
		background: transparent;
		border: none;
		color: var(--text-secondary);
		cursor: pointer;
		display: flex;
		align-items: center;
		gap: 6px;
		font-size: 0.85rem;
		font-weight: 600;
		padding: 4px 8px;
		border-radius: 4px;
		transition: var(--transition-fast);
		border-left: 1px solid rgba(255, 255, 255, 0.1);
		padding-left: 12px;
		margin-left: 4px;
	}

	.logout-btn:hover {
		color: var(--accent-red);
		text-shadow: 0 0 8px rgba(255, 75, 92, 0.3);
	}

	.logout-icon {
		transition: transform 0.2s ease;
	}

	.logout-btn:hover .logout-icon {
		transform: translateX(2px);
	}

	@media (max-width: 640px) {
		.header-content {
			flex-direction: column;
			gap: 16px;
			padding: 16px;
		}

		.nav-links {
			gap: 16px;
		}
	}
</style>
