<template>
  <pre class="language-java"><code class="language-java">{{ headers }}</code></pre>
</template>
<script setup>
import {computed, onMounted, onUpdated} from 'vue'
import Prism from "prismjs";

const props = defineProps({
  header: {},
})

const headers = computed(() => {
  let arr = []
  if (props.header === undefined || props.header === null) {
    return ''
  }
  Object.entries(props.header).forEach((v) => {
    arr.push(v[0] + ": " + v[1])
  })
  return arr.join("\n")
})
onMounted(() => {
  Prism.highlightAll()
});
onUpdated(() => {
      Prism.highlightAll()
    }
)
</script>