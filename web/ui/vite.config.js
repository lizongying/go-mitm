import {fileURLToPath, URL} from 'node:url'

import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'
import Components from 'unplugin-vue-components/vite';
import {AntDesignVueResolver} from 'unplugin-vue-components/resolvers';
import {prismjsPlugin} from "vite-plugin-prismjs"

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        Components({
            resolvers: [
                AntDesignVueResolver({
                    importStyle: false, // css in js
                }),
            ],
        }),
        prismjsPlugin({
            languages: ['json', 'html', 'shell', 'javascript', 'python', 'go', 'java', 'php'],
            plugins: ['line-numbers', 'copy-to-clipboard'],
            theme: 'okaidia',
            css: true,
        }),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    }
})
