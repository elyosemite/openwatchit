import { defineCollection, z } from 'astro:content';
import { docsLoader } from '@astrojs/starlight/loaders';
import { docsSchema } from '@astrojs/starlight/schema';
import { glob } from 'astro/loaders';

export const collections = {
	docs: defineCollection({ loader: docsLoader(), schema: docsSchema() }),

	// Reads versioned markdown files from the repo-root changelog/ directory.
	// Each file (e.g. v0.1.0.md) becomes an entry in the changelog collection.
	changelog: defineCollection({
		loader: glob({ pattern: 'v*.md', base: '../changelog' }),
		schema: z.object({
			title: z.string().optional(),
		}),
	}),
};
