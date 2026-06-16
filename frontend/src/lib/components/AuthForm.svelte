<script lang="ts">
	import { auth } from '$lib/auth.svelte';

	// State cho Form Login / Register
	let authMode = $state<'login' | 'register'>('login');
	let nameInput = $state('');
	let emailInput = $state('');
	let passwordInput = $state('');
	let confirmPasswordInput = $state('');
	let showPassword = $state(false);
	let showConfirmPassword = $state(false);
	let formError = $state<string | null>(null);
	let formSuccess = $state<string | null>(null);

	const handleAuthSubmit = async (e: Event) => {
		e.preventDefault();
		formError = null;
		formSuccess = null;

		if (authMode === 'register' && !nameInput.trim()) {
			formError = 'Vui lòng nhập họ và tên của bạn.';
			return;
		}
		if (!emailInput.trim()) {
			formError = 'Vui lòng nhập địa chỉ email.';
			return;
		}
		if (!passwordInput) {
			formError = 'Vui lòng nhập mật khẩu.';
			return;
		}
		if (passwordInput.length < 4) {
			formError = 'Mật khẩu phải có độ dài tối thiểu 4 ký tự.';
			return;
		}
		if (authMode === 'register') {
			if (!confirmPasswordInput) {
				formError = 'Vui lòng xác nhận lại mật khẩu.';
				return;
			}
			if (passwordInput !== confirmPasswordInput) {
				formError = 'Mật khẩu xác nhận không khớp.';
				return;
			}
		}

		if (authMode === 'login') {
			const success = await auth.login(emailInput, passwordInput);
			if (success) {
				formSuccess = 'Đăng nhập thành công! Đang chuyển hướng...';
				nameInput = '';
				emailInput = '';
				passwordInput = '';
				confirmPasswordInput = '';
				showPassword = false;
				showConfirmPassword = false;
			} else {
				formError = auth.error;
			}
		} else {
			const success = await auth.register(nameInput, emailInput, passwordInput);
			if (success) {
				formSuccess = 'Đăng ký tài khoản thành công!';
				nameInput = '';
				emailInput = '';
				passwordInput = '';
				confirmPasswordInput = '';
				showPassword = false;
				showConfirmPassword = false;
			} else {
				formError = auth.error;
			}
		}
	};
</script>

<div class="auth-container">
	<div class="auth-card">
		<div class="auth-header">
			<div class="gradient-text auth-title">🐱 Meow App</div>
			<p class="auth-subtitle">
				{#if authMode === 'login'}
					Đăng nhập để giám sát hiệu năng hệ thống
				{:else}
					Tạo tài khoản quản trị hệ thống mới
				{/if}
			</p>
		</div>

		<!-- Tab control chuyển đổi nhanh giữa đăng nhập/đăng ký -->
		<div class="auth-tabs">
			<button
				class="auth-tab {authMode === 'login' ? 'active' : ''}"
				onclick={() => {
					authMode = 'login';
					formError = null;
					formSuccess = null;
					passwordInput = '';
					confirmPasswordInput = '';
					showPassword = false;
					showConfirmPassword = false;
				}}
			>
				Đăng nhập
			</button>
			<button
				class="auth-tab {authMode === 'register' ? 'active' : ''}"
				onclick={() => {
					authMode = 'register';
					formError = null;
					formSuccess = null;
					passwordInput = '';
					confirmPasswordInput = '';
					showPassword = false;
					showConfirmPassword = false;
				}}
			>
				Đăng ký
			</button>
		</div>

		<form class="auth-form" onsubmit={handleAuthSubmit}>
			{#if authMode === 'register'}
				<div class="form-group">
					<label for="name" class="form-label">Họ và tên</label>
					<div class="input-wrapper">
						<span class="input-icon">
							<svg viewBox="0 0 24 24" width="18" height="18">
								<path
									fill="currentColor"
									d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
								/>
							</svg>
						</span>
						<input
							type="text"
							id="name"
							placeholder="Nguyễn Văn A"
							class="input-field"
							bind:value={nameInput}
							disabled={auth.loading}
						/>
					</div>
				</div>
			{/if}

			<div class="form-group">
				<label for="email" class="form-label">Email</label>
				<div class="input-wrapper">
					<span class="input-icon">
						<svg viewBox="0 0 24 24" width="18" height="18">
							<path
								fill="currentColor"
								d="M20 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"
							/>
						</svg>
					</span>
					<input
						type="email"
						id="email"
						placeholder="admin@meow.io"
						class="input-field"
						bind:value={emailInput}
						disabled={auth.loading}
					/>
				</div>
			</div>

			<div class="form-group">
				<label for="password" class="form-label">Mật khẩu</label>
				<div class="input-wrapper">
					<span class="input-icon">
						<svg viewBox="0 0 24 24" width="18" height="18">
							<path
								fill="currentColor"
								d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"
							/>
						</svg>
					</span>
					<input
						type={showPassword ? 'text' : 'password'}
						id="password"
						placeholder="••••••••"
						class="input-field"
						bind:value={passwordInput}
						disabled={auth.loading}
					/>
					<button
						type="button"
						class="toggle-password"
						onclick={() => showPassword = !showPassword}
						aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiển thị mật khẩu'}
						disabled={auth.loading}
					>
						{#if showPassword}
							<svg viewBox="0 0 24 24" width="18" height="18">
								<path
									fill="currentColor"
									d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.82l2.92 2.92c1.51-1.2 2.7-2.78 3.44-4.74-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-4 .7l2.17 2.17C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46C3.08 8.3 1.78 10.02 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.34-4.3l1.6 1.6C12.94 7.04 12.48 7 12 7c-2.76 0-5 2.24-5 5 0 .48.04.94.1 1.47l1.6 1.6C8.83 14.47 9 13.27 9 12c0-1.66 1.34-3 3-3 1.27 0 2.47-.17 3.07-.3z"
								/>
							</svg>
						{:else}
							<svg viewBox="0 0 24 24" width="18" height="18">
								<path
									fill="currentColor"
									d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"
								/>
							</svg>
						{/if}
					</button>
				</div>
			</div>

			{#if authMode === 'register'}
				<div class="form-group">
					<label for="confirm-password" class="form-label">Xác nhận mật khẩu</label>
					<div class="input-wrapper">
						<span class="input-icon">
							<svg viewBox="0 0 24 24" width="18" height="18">
								<path
									fill="currentColor"
									d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"
								/>
							</svg>
						</span>
						<input
							type={showConfirmPassword ? 'text' : 'password'}
							id="confirm-password"
							placeholder="••••••••"
							class="input-field"
							bind:value={confirmPasswordInput}
							disabled={auth.loading}
						/>
						<button
							type="button"
							class="toggle-password"
							onclick={() => showConfirmPassword = !showConfirmPassword}
							aria-label={showConfirmPassword ? 'Ẩn mật khẩu' : 'Hiển thị mật khẩu'}
							disabled={auth.loading}
						>
							{#if showConfirmPassword}
								<svg viewBox="0 0 24 24" width="18" height="18">
									<path
										fill="currentColor"
										d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.82l2.92 2.92c1.51-1.2 2.7-2.78 3.44-4.74-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-4 .7l2.17 2.17C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46C3.08 8.3 1.78 10.02 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.34-4.3l1.6 1.6C12.94 7.04 12.48 7 12 7c-2.76 0-5 2.24-5 5 0 .48.04.94.1 1.47l1.6 1.6C8.83 14.47 9 13.27 9 12c0-1.66 1.34-3 3-3 1.27 0 2.47-.17 3.07-.3z"
									/>
								</svg>
							{:else}
								<svg viewBox="0 0 24 24" width="18" height="18">
									<path
										fill="currentColor"
										d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"
									/>
								</svg>
							{/if}
						</button>
					</div>
				</div>
			{/if}

			{#if formError}
				<div class="auth-error">
					<svg viewBox="0 0 24 24" width="18" height="18">
						<path
							fill="currentColor"
							d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"
						/>
					</svg>
					<span>{formError}</span>
				</div>
			{/if}

			{#if formSuccess}
				<div class="auth-success">
					<svg viewBox="0 0 24 24" width="18" height="18">
						<path
							fill="currentColor"
							d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"
						/>
					</svg>
					<span>{formSuccess}</span>
				</div>
			{/if}

			<button type="submit" class="gradient-btn btn-submit" disabled={auth.loading}>
				{#if auth.loading}
					<span class="spinner"></span>
					Đang xử lý...
				{:else if authMode === 'login'}
					Đăng nhập
				{:else}
					Đăng ký tài khoản
				{/if}
			</button>
		</form>
	</div>
</div>

<style>
	.auth-container {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: calc(100vh - 200px);
		padding: 20px;
	}

	.auth-card {
		width: 100%;
		max-width: 420px;
		background: var(--glass-bg);
		border: 1px solid var(--glass-border);
		backdrop-filter: blur(16px);
		-webkit-backdrop-filter: blur(16px);
		border-radius: var(--radius-md);
		padding: 40px 30px;
		box-shadow: var(--shadow-lg);
		transition: var(--transition-smooth);
	}

	.auth-card:hover {
		border-color: rgba(0, 242, 254, 0.2);
		box-shadow:
			var(--shadow-lg),
			0 0 30px rgba(0, 242, 254, 0.1);
	}

	.auth-header {
		text-align: center;
		margin-bottom: 30px;
	}

	.auth-title {
		font-family: var(--font-display);
		font-size: 2.2rem;
		font-weight: 800;
		letter-spacing: -0.02em;
		margin-bottom: 8px;
	}

	.auth-subtitle {
		color: var(--text-secondary);
		font-size: 0.9rem;
	}

	/* Tabs */
	.auth-tabs {
		display: flex;
		background: rgba(255, 255, 255, 0.03);
		border: 1px solid var(--glass-border);
		border-radius: 30px;
		padding: 4px;
		margin-bottom: 30px;
		position: relative;
	}

	.auth-tab {
		flex: 1;
		background: transparent;
		border: none;
		color: var(--text-secondary);
		font-family: var(--font-sans);
		font-weight: 600;
		font-size: 0.92rem;
		padding: 10px 0;
		border-radius: 26px;
		cursor: pointer;
		transition: var(--transition-smooth);
		z-index: 1;
	}

	.auth-tab.active {
		color: var(--bg-primary);
		background: linear-gradient(135deg, var(--accent-cyan) 0%, var(--accent-purple) 100%);
		box-shadow: 0 4px 10px rgba(0, 242, 254, 0.2);
	}

	/* Form Group */
	.auth-form {
		display: flex;
		flex-direction: column;
		gap: 20px;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.form-label {
		color: var(--text-primary);
		font-size: 0.88rem;
		font-weight: 600;
		letter-spacing: 0.02em;
	}

	.input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.input-icon {
		position: absolute;
		left: 14px;
		color: var(--text-muted);
		display: flex;
		align-items: center;
		justify-content: center;
		pointer-events: none;
	}

	.input-field {
		width: 100%;
		background: rgba(255, 255, 255, 0.02);
		border: 1px solid var(--glass-border);
		border-radius: var(--radius-sm);
		padding: 12px 42px 12px 42px;
		color: var(--text-primary);
		font-family: var(--font-sans);
		font-size: 0.95rem;
		transition: var(--transition-fast);
	}

	.input-field:focus {
		outline: none;
		background: rgba(255, 255, 255, 0.05);
		border-color: var(--accent-cyan);
		box-shadow: 0 0 10px rgba(0, 242, 254, 0.15);
	}

	/* Alert messages */
	.auth-error {
		background: rgba(255, 75, 92, 0.1);
		border: 1px solid rgba(255, 75, 92, 0.25);
		color: var(--accent-red);
		padding: 12px 16px;
		border-radius: var(--radius-sm);
		font-size: 0.88rem;
		display: flex;
		align-items: center;
		gap: 10px;
		animation: shake 0.4s ease-in-out;
	}

	.auth-success {
		background: rgba(56, 239, 125, 0.1);
		border: 1px solid rgba(56, 239, 125, 0.25);
		color: var(--accent-green);
		padding: 12px 16px;
		border-radius: var(--radius-sm);
		font-size: 0.88rem;
		display: flex;
		align-items: center;
		gap: 10px;
	}

	/* Submit Button */
	.btn-submit {
		margin-top: 10px;
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 10px;
	}

	.btn-submit:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
		box-shadow: none;
	}

	/* Loading Spinner */
	.spinner {
		width: 18px;
		height: 18px;
		border: 2px solid rgba(255, 255, 255, 0.3);
		border-top-color: #ffffff;
		border-radius: 50%;
		animation: spin-slow 0.8s linear infinite;
	}

	/* Shake animation for errors */
	@keyframes shake {
		0%,
		100% {
			transform: translateX(0);
		}
		20%,
		60% {
			transform: translateX(-6px);
		}
		40%,
		80% {
			transform: translateX(6px);
		}
	}

	.toggle-password {
		position: absolute;
		right: 14px;
		background: transparent;
		border: none;
		color: var(--text-muted);
		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;
		padding: 0;
		transition: var(--transition-fast);
	}

	.toggle-password:hover {
		color: var(--accent-cyan);
		filter: drop-shadow(0 0 4px rgba(0, 242, 254, 0.4));
	}

	.toggle-password:focus {
		outline: none;
		color: var(--accent-cyan);
	}
</style>
