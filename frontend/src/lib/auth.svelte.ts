import { PUBLIC_BACKEND_URL } from '$env/static/public';

export class AuthStore {
	token = $state<string | null>(null);
	email = $state<string | null>(null);
	name = $state<string | null>(null);
	error = $state<string | null>(null);
	loading = $state<boolean>(false);

	constructor() {
		if (typeof window !== 'undefined') {
			this.token = localStorage.getItem('auth_token');
			this.email = localStorage.getItem('auth_email');
			this.name = localStorage.getItem('auth_name');
		}
	}

	get isAuthenticated() {
		return !!this.token;
	}

	async login(emailInput: string, passwordInput: string): Promise<boolean> {
		this.loading = true;
		this.error = null;
		try {
			const res = await fetch(`${PUBLIC_BACKEND_URL}/api/auth/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email: emailInput, password: passwordInput })
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Đăng nhập thất bại. Vui lòng kiểm tra lại thông tin.');
			}

			this.token = data.token;
			this.email = emailInput;
			this.error = null;

			if (typeof window !== 'undefined') {
				localStorage.setItem('auth_token', data.token);
				localStorage.setItem('auth_email', emailInput);
			}
			return true;
		} catch (err) {
			this.error = (err as Error).message || 'Lỗi kết nối máy chủ';
			return false;
		} finally {
			this.loading = false;
		}
	}

	async register(nameInput: string, emailInput: string, passwordInput: string): Promise<boolean> {
		this.loading = true;
		this.error = null;
		try {
			const res = await fetch(`${PUBLIC_BACKEND_URL}/api/auth/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ name: nameInput, email: emailInput, password: passwordInput })
			});

			const data = await res.json();
			if (!res.ok) {
				throw new Error(data.error || 'Đăng ký thất bại. Email có thể đã tồn tại.');
			}

			this.token = data.token;
			this.name = nameInput;
			this.email = emailInput;
			this.error = null;

			if (typeof window !== 'undefined') {
				localStorage.setItem('auth_token', data.token);
				localStorage.setItem('auth_name', nameInput);
				localStorage.setItem('auth_email', emailInput);
			}
			return true;
		} catch (err) {
			this.error = (err as Error).message || 'Lỗi kết nối máy chủ';
			return false;
		} finally {
			this.loading = false;
		}
	}

	logout() {
		this.token = null;
		this.email = null;
		this.name = null;
		this.error = null;

		if (typeof window !== 'undefined') {
			localStorage.removeItem('auth_token');
			localStorage.removeItem('auth_email');
			localStorage.removeItem('auth_name');
		}
	}
}

export const auth = new AuthStore();
