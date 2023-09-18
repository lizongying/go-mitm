<template>
    <pre class="language-md"><code class="language-md"><span v-for="{k, v} in headers" key="k">{{ k }}: {{
        v
        }}
</span></code></pre>
</template>
<script setup>
import {computed} from 'vue'

const props = defineProps({
    header: {},
})
const headers = computed(() => {
    let arr = []
    if (props.header === undefined) {
        return arr
    }
    Object.entries(props.header).forEach((v) => {
        arr.push({k: v[0], v: v[1]})
    })
    return arr
})
</script>