// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
	integrations: [
		starlight({
			title: 'OpenWatchIt',
			description: 'One query language for all your observability backends.',
			customCss: ['./src/styles/custom.css'],
			head: [
				{
					tag: 'link',
					attrs: {
						rel: 'preconnect',
						href: 'https://fonts.googleapis.com',
					},
				},
				{
					tag: 'link',
					attrs: {
						rel: 'preconnect',
						href: 'https://fonts.gstatic.com',
						crossorigin: true,
					},
				},
				{
					tag: 'link',
					attrs: {
						rel: 'stylesheet',
						href: 'https://fonts.googleapis.com/css2?family=Lora:wght@400;500;600;700&display=swap',
					},
				},
			],
			social: [
				{ icon: 'github', label: 'GitHub', href: 'https://github.com/elyosemite/openwatchit' },
			],
			sidebar: [
				{
					label: 'Getting Started',
					items: [
						{ label: 'Introduction', slug: 'getting-started/introduction' },
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Configuration', slug: 'getting-started/configuration' },
					],
				},
				{
					label: 'CLI Reference',
					items: [
						{ label: 'owit logs', slug: 'cli/logs' },
						{ label: 'owit metrics', slug: 'cli/metrics' },
						{ label: 'owit traces', slug: 'cli/traces' },
						{ label: 'owit profiles', slug: 'cli/profiles' },
						{ label: 'owit tail', slug: 'cli/tail' },
					],
				},
				{
					label: 'OWL Language',
					items: [
						{ label: 'Filter Operators', slug: 'owl/operators' },
						{ label: 'Flags Reference', slug: 'owl/flags' },
					],
				},
				{
					label: 'Plugin Development',
					items: [
						{ label: 'Writing a Plugin', slug: 'plugins/writing-a-plugin' },
						{ label: 'gRPC Contract', slug: 'plugins/grpc-contract' },
					],
				},
				{
					label: 'Changelog',
					link: '/changelog/',
				},
			],
		}),
	],
});
