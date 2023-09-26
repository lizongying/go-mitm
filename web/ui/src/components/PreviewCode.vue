<template>
    <pre :class="'hx-scroll ' + lineNumbers"><code :class="'language-'+ type"
                                                   v-html="html"></code></pre>
</template>
<script setup>
import {computed, onMounted, onUpdated} from 'vue'
import Prism from 'prismjs'

const props = defineProps({
  code: {
    type: String,
    default: '',
  },
  type: {
    type: String,
    default: 'html',
  },
  lineNumbers: {
    type: Boolean,
    default: false,
  },
})
const lineNumbers = computed(() => {
  return props.lineNumbers ? 'line-numbers' : ''
});
const html = computed(() => {
  return props.code === undefined || props.code == null || props.code === '' ? '' : Prism.highlight(props.code, Prism.languages[props.type], props.type)
});
onMounted(() => {
  Prism.highlightAll()
});
onUpdated(() => {
      Prism.highlightAll()
    }
)
</script>