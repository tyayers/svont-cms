import { defineConfig } from 'astro/config';

// https://astro.build/config
import lit from "@astrojs/lit";

// https://astro.build/config
import node from "@astrojs/node";

// https://astro.build/config
export default defineConfig({
  output: 'server',
  integrations: [lit()],
  adapter: node({
    mode: 'standalone'
  })
});