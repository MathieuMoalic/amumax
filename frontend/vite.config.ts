import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:35367',
				changeOrigin: true
			},
			'/ws': {
				target: 'ws://localhost:35367',
				ws: true
			},
		}
	}
});
