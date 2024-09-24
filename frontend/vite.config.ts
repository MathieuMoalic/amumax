import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/api': {
				target: 'http://localhost:35367',
				changeOrigin: true
			}
		}
	},
	define: {
		'process.env': process.env // Ensures process.env variables are passed
	}
});
