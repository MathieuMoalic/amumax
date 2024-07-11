import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://amumax-backend-dev:35367',
				changeOrigin: true,
			},
			'/ws': {
				target: 'ws://amumax-backend-dev:35367',
				changeOrigin: true,
				ws: true,
			},
		}
	}
});
